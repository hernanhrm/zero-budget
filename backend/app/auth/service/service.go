package service

import (
	"context"
	"strings"
	"time"

	"backend/app/auth/domain"
	organizationDomain "backend/app/organization/domain"
	roleDomain "backend/app/role/domain"
	userDomain "backend/app/user/domain"
	workspaceDomain "backend/app/workspace/domain"
	workspaceMemberDomain "backend/app/workspace_member/domain"
	basedomain "backend/domain"
	apperrors "backend/domain/errors"
	"backend/infra/dafi"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/samber/oops"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenExpiry  = 5 * time.Minute
	refreshTokenExpiry = 7 * 24 * time.Hour
)

type service struct {
	userSvc            userDomain.Service
	organizationSvc    organizationDomain.Service
	workspaceSvc       workspaceDomain.Service
	roleSvc            roleDomain.Service
	workspaceMemberSvc workspaceMemberDomain.Service
	logger             basedomain.Logger
	jwtSecret          []byte
}

func New(
	userSvc userDomain.Service,
	organizationSvc organizationDomain.Service,
	workspaceSvc workspaceDomain.Service,
	roleSvc roleDomain.Service,
	workspaceMemberSvc workspaceMemberDomain.Service,
	logger basedomain.Logger,
	jwtSecret string,
) domain.Service {
	return service{
		userSvc:            userSvc,
		organizationSvc:    organizationSvc,
		workspaceSvc:       workspaceSvc,
		roleSvc:            roleSvc,
		workspaceMemberSvc: workspaceMemberSvc,
		logger:             logger.With("component", "auth.service"),
		jwtSecret:          []byte(jwtSecret),
	}
}

func (s service) SignupWithEmail(ctx context.Context, input domain.SignupWithEmail) (domain.SignupResponse, error) {
	if err := input.Validate(ctx); err != nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	userInput := userDomain.CreateUser{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  string(hashedPassword),
	}

	if err := s.userSvc.Create(ctx, userInput); err != nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	userCriteria := dafi.Where("email", dafi.Equal, input.Email)
	user, err := s.userSvc.FindOne(ctx, userCriteria)
	if err != nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	orgSlug := generateSlug(input.OrganizationName)

	orgInput := organizationDomain.CreateOrganization{
		Name:    input.OrganizationName,
		Slug:    orgSlug,
		OwnerID: user.ID,
	}

	if err := s.organizationSvc.Create(ctx, orgInput); err != nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	orgCriteria := dafi.Where("slug", dafi.Equal, orgSlug)
	org, err := s.organizationSvc.FindOne(ctx, orgCriteria)
	if err != nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	workspaceSlug := generateSlug(input.OrganizationName)

	workspaceInput := workspaceDomain.CreateWorkspace{
		OrganizationID: org.ID,
		Name:           input.OrganizationName,
		Slug:           workspaceSlug,
	}

	if err := s.workspaceSvc.Create(ctx, workspaceInput); err != nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	workspaceCriteria := dafi.Where("organization_id", dafi.Equal, org.ID).
		And("slug", dafi.Equal, workspaceSlug)
	workspace, err := s.workspaceSvc.FindOne(ctx, workspaceCriteria)
	if err != nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	roleCriteria := dafi.Where("name", dafi.Equal, "Admin").
		And("workspace_id", dafi.IsNull, nil)
	adminRole, err := s.roleSvc.FindOne(ctx, roleCriteria)
	if err != nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	workspaceIDParsed, _ := uuid.Parse(workspace.ID)
	userIDParsed, _ := uuid.Parse(user.ID)
	roleIDParsed, _ := uuid.Parse(adminRole.ID)

	memberInput := workspaceMemberDomain.CreateWorkspaceMember{
		WorkspaceID: workspaceIDParsed,
		UserID:      userIDParsed,
		RoleID:      roleIDParsed,
	}

	if err := s.workspaceMemberSvc.Create(ctx, memberInput); err != nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	tokenPair, err := s.generateTokenPair(user.ID, workspace.ID)
	if err != nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("user signed up",
		"email", input.Email,
		"user_id", user.ID,
		"workspace_id", workspace.ID,
	)

	return domain.SignupResponse{
		TokenPair: tokenPair,
		User: domain.UserInfo{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
		Workspace: domain.WorkspaceInfo{
			ID:   workspace.ID,
			Name: workspace.Name,
			Slug: workspace.Slug,
		},
	}, nil
}

func (s service) generateTokenPair(userID, workspaceID string) (domain.TokenPair, error) {
	now := time.Now()

	accessClaims := jwt.MapClaims{
		"user_id":      userID,
		"workspace_id": workspaceID,
		"exp":          now.Add(accessTokenExpiry).Unix(),
		"iat":          now.Unix(),
		"type":         "access",
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return domain.TokenPair{}, err
	}

	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     now.Add(refreshTokenExpiry).Unix(),
		"iat":     now.Unix(),
		"type":    "refresh",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		return domain.TokenPair{}, err
	}

	return domain.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    now.Add(accessTokenExpiry),
	}, nil
}

func generateSlug(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}
