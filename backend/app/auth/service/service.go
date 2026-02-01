package service

import (
	"context"
	"fmt"
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

	// Check if user with email already exists
	existingUserCriteria := dafi.Where("email", dafi.Equal, input.Email)
	_, err := s.userSvc.FindOne(ctx, existingUserCriteria)
	if err == nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).
			In(apperrors.LayerService).
			Code(apperrors.CodeAlreadyExists).
			With("email", input.Email).
			Public("A user with this email already exists").
			Errorf("email already exists: %s", input.Email)
	}

	// Generate unique slugs by appending short UUID suffix
	baseSlug := generateSlug(input.OrganizationName)
	uniqueSuffix := generateShortSuffix()
	orgSlug := fmt.Sprintf("%s-%s", baseSlug, uniqueSuffix)
	workspaceSlug := baseSlug

	// Check if organization slug already exists
	existingOrgCriteria := dafi.Where("slug", dafi.Equal, orgSlug)
	_, err = s.organizationSvc.FindOne(ctx, existingOrgCriteria)
	if err == nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).
			In(apperrors.LayerService).
			Code(apperrors.CodeAlreadyExists).
			With("slug", orgSlug).
			Public("An organization with this name already exists").
			Errorf("organization slug already exists: %s", orgSlug)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	// Generate user ID in service
	userID := uuid.New().String()

	// Generate organization ID in service
	orgID := uuid.New()

	// Create user with generated ID
	userInput := userDomain.CreateUser{
		ID:        userID,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  string(hashedPassword),
	}

	if err := s.userSvc.Create(ctx, userInput); err != nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	// Create organization
	orgInput := organizationDomain.CreateOrganization{
		ID:      orgID,
		Name:    input.OrganizationName,
		Slug:    orgSlug,
		OwnerID: uuid.MustParse(userID),
	}

	if err := s.organizationSvc.Create(ctx, orgInput); err != nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	// Query back the organization to get the generated ID
	orgCriteria := dafi.Where("slug", dafi.Equal, orgSlug)
	org, err := s.organizationSvc.FindOne(ctx, orgCriteria)
	if err != nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	// Check if workspace already exists for this organization
	existingWorkspaceCriteria := dafi.Where("organizationId", dafi.Equal, org.ID).
		And("slug", dafi.Equal, workspaceSlug)
	_, err = s.workspaceSvc.FindOne(ctx, existingWorkspaceCriteria)
	if err == nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).
			In(apperrors.LayerService).
			Code(apperrors.CodeAlreadyExists).
			With("organizationId", org.ID).
			With("slug", workspaceSlug).
			Public("A workspace with this name already exists in your organization").
			Errorf("workspace slug already exists for org %s: %s", org.ID, workspaceSlug)
	}

	// Create workspace
	workspaceInput := workspaceDomain.CreateWorkspace{
		ID:             uuid.New(),
		OrganizationID: org.ID,
		Name:           input.OrganizationName,
		Slug:           workspaceSlug,
	}

	if err := s.workspaceSvc.Create(ctx, workspaceInput); err != nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	// Find Admin role (system role with workspace_id = NULL)
	roleCriteria := dafi.Where("name", dafi.Equal, "Admin").
		And("workspaceId", dafi.IsNull, nil)
	adminRole, err := s.roleSvc.FindOne(ctx, roleCriteria)
	if err != nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	userIDParsed := uuid.MustParse(userID)

	memberInput := workspaceMemberDomain.CreateWorkspaceMember{
		WorkspaceID: workspaceInput.ID,
		UserID:      userIDParsed,
		RoleID:      adminRole.ID,
	}

	if err := s.workspaceMemberSvc.Create(ctx, memberInput); err != nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	tokenPair, err := s.generateTokenPair(userID, workspaceInput.ID.String())
	if err != nil {
		return domain.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("user signed up",
		"email", input.Email,
		"user_id", userID,
		"workspace_id", workspaceInput.ID,
	)

	return domain.SignupResponse{
		TokenPair: tokenPair,
		User: domain.UserInfo{
			ID:        userID,
			Email:     input.Email,
			FirstName: input.FirstName,
			LastName:  input.LastName,
		},
		Workspace: domain.WorkspaceInfo{
			ID:   workspaceInput.ID.String(),
			Name: workspaceInput.Name,
			Slug: workspaceInput.Slug,
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

func generateShortSuffix() string {
	return uuid.New().String()[:8]
}
