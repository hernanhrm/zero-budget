package core

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"backend/core/auth/port"
	organizationPort "backend/core/organization/port"
	rolePort "backend/core/role/port"
	userPort "backend/core/user/port"
	workspacePort "backend/core/workspace/port"
	workspaceMemberPort "backend/core/workspace_member/port"
	basedomain "backend/port"
	apperrors "backend/port/errors"
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
	userSvc            userPort.Service
	organizationSvc    organizationPort.Service
	workspaceSvc       workspacePort.Service
	roleSvc            rolePort.Service
	workspaceMemberSvc workspaceMemberPort.Service
	logger             basedomain.Logger
	jwtSecret          []byte
}

func New(
	userSvc userPort.Service,
	organizationSvc organizationPort.Service,
	workspaceSvc workspacePort.Service,
	roleSvc rolePort.Service,
	workspaceMemberSvc workspaceMemberPort.Service,
	logger basedomain.Logger,
	jwtSecret string,
) port.Service {
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

func (s service) WithTx(tx basedomain.Transaction) port.Service {
	return service{
		userSvc:            s.userSvc.WithTx(tx),
		organizationSvc:    s.organizationSvc.WithTx(tx),
		workspaceSvc:       s.workspaceSvc.WithTx(tx),
		roleSvc:            s.roleSvc.WithTx(tx),
		workspaceMemberSvc: s.workspaceMemberSvc.WithTx(tx),
		logger:             s.logger,
		jwtSecret:          s.jwtSecret,
	}
}

func (s service) SignupWithEmail(ctx context.Context, input port.SignupWithEmail) (port.SignupResponse, error) {
	if err := input.Validate(ctx); err != nil {
		return port.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	// Check if user with email already exists
	existingUserCriteria := dafi.Where("email", dafi.Equal, input.Email)
	_, err := s.userSvc.FindOne(ctx, existingUserCriteria)
	if err == nil {
		return port.SignupResponse{}, oops.WithContext(ctx).
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
		return port.SignupResponse{}, oops.WithContext(ctx).
			In(apperrors.LayerService).
			Code(apperrors.CodeAlreadyExists).
			With("slug", orgSlug).
			Public("An organization with this name already exists").
			Errorf("organization slug already exists: %s", orgSlug)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return port.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	// Generate user ID in service
	userID := uuid.New().String()

	// Generate organization ID in service
	orgID := uuid.New()

	// Create user with generated ID
	userInput := userPort.CreateUser{
		ID:        userID,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  string(hashedPassword),
	}

	if err := s.userSvc.Create(ctx, userInput); err != nil {
		return port.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	// Create organization
	orgInput := organizationPort.CreateOrganization{
		ID:      orgID,
		Name:    input.OrganizationName,
		Slug:    orgSlug,
		OwnerID: uuid.MustParse(userID),
	}

	if err := s.organizationSvc.Create(ctx, orgInput); err != nil {
		return port.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	// Query back the organization to get the generated ID
	orgCriteria := dafi.Where("slug", dafi.Equal, orgSlug)
	org, err := s.organizationSvc.FindOne(ctx, orgCriteria)
	if err != nil {
		return port.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	// Check if workspace already exists for this organization
	existingWorkspaceCriteria := dafi.Where("organizationId", dafi.Equal, org.ID).
		And("slug", dafi.Equal, workspaceSlug)
	_, err = s.workspaceSvc.FindOne(ctx, existingWorkspaceCriteria)
	if err == nil {
		return port.SignupResponse{}, oops.WithContext(ctx).
			In(apperrors.LayerService).
			Code(apperrors.CodeAlreadyExists).
			With("organizationId", org.ID).
			With("slug", workspaceSlug).
			Public("A workspace with this name already exists in your organization").
			Errorf("workspace slug already exists for org %s: %s", org.ID, workspaceSlug)
	}

	// Create workspace
	workspaceInput := workspacePort.CreateWorkspace{
		ID:             uuid.New(),
		OrganizationID: org.ID,
		Name:           input.OrganizationName,
		Slug:           workspaceSlug,
	}

	if err := s.workspaceSvc.Create(ctx, workspaceInput); err != nil {
		return port.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	// Find Admin role (system role with workspace_id = NULL)
	roleCriteria := dafi.Where("name", dafi.Equal, "Admin").
		And("workspaceId", dafi.IsNull, nil)
	adminRole, err := s.roleSvc.FindOne(ctx, roleCriteria)
	if err != nil {
		return port.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	userIDParsed := uuid.MustParse(userID)

	memberInput := workspaceMemberPort.CreateWorkspaceMember{
		WorkspaceID: workspaceInput.ID,
		UserID:      userIDParsed,
		RoleID:      adminRole.ID,
	}

	if err := s.workspaceMemberSvc.Create(ctx, memberInput); err != nil {
		return port.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	tokenPair, err := s.generateTokenPair(userID, workspaceInput.ID.String())
	if err != nil {
		return port.SignupResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("user signed up",
		"email", input.Email,
		"user_id", userID,
		"workspace_id", workspaceInput.ID,
	)

	return port.SignupResponse{
		TokenPair: tokenPair,
		User: port.UserInfo{
			ID:        userID,
			Email:     input.Email,
			FirstName: input.FirstName,
			LastName:  input.LastName,
		},
		Workspace: port.WorkspaceInfo{
			ID:   workspaceInput.ID.String(),
			Name: workspaceInput.Name,
			Slug: workspaceInput.Slug,
		},
	}, nil
}

func (s service) generateTokenPair(userID, workspaceID string) (port.TokenPair, error) {
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
		return port.TokenPair{}, err
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
		return port.TokenPair{}, err
	}

	return port.TokenPair{
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

func (s service) LoginWithEmail(ctx context.Context, input port.LoginWithEmail) (port.LoginResponse, error) {
	if err := input.Validate(ctx); err != nil {
		return port.LoginResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	userCriteria := dafi.Where("email", dafi.Equal, input.Email)
	user, err := s.userSvc.FindOne(ctx, userCriteria)
	if err != nil {
		return port.LoginResponse{}, oops.WithContext(ctx).
			In(apperrors.LayerService).
			Code(apperrors.CodeUnauthorized).
			With("email", input.Email).
			Public("Invalid email or password").
			Errorf("user not found with email: %s", input.Email)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(input.Password)); err != nil {
		return port.LoginResponse{}, oops.WithContext(ctx).
			In(apperrors.LayerService).
			Code(apperrors.CodeUnauthorized).
			With("email", input.Email).
			Public("Invalid email or password").
			Errorf("invalid password for user: %s", input.Email)
	}

	memberCriteria := dafi.Where("userId", dafi.Equal, user.ID)
	member, err := s.workspaceMemberSvc.FindOne(ctx, memberCriteria)
	if err != nil {
		return port.LoginResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	workspaceCriteria := dafi.Where("id", dafi.Equal, member.WorkspaceID)
	workspace, err := s.workspaceSvc.FindOne(ctx, workspaceCriteria)
	if err != nil {
		return port.LoginResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	tokenPair, err := s.generateTokenPair(user.ID, member.WorkspaceID.String())
	if err != nil {
		return port.LoginResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("user logged in",
		"email", input.Email,
		"user_id", user.ID,
		"workspace_id", member.WorkspaceID,
	)

	return port.LoginResponse{
		TokenPair: tokenPair,
		User: port.UserInfo{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
		Workspace: port.WorkspaceInfo{
			ID:   member.WorkspaceID.String(),
			Name: workspace.Name,
			Slug: workspace.Slug,
		},
	}, nil
}

func (s service) RefreshToken(ctx context.Context, accessToken, refreshToken string) (port.RefreshResponse, error) {
	accessClaims, err := s.parseToken(accessToken, true)
	if err != nil {
		return port.RefreshResponse{}, oops.WithContext(ctx).
			In(apperrors.LayerService).
			Code(apperrors.CodeUnauthorized).
			Public("Invalid or expired refresh token").
			Errorf("failed to parse access token: %w", err)
	}

	userID, ok := accessClaims["user_id"].(string)
	if !ok || userID == "" {
		return port.RefreshResponse{}, oops.WithContext(ctx).
			In(apperrors.LayerService).
			Code(apperrors.CodeUnauthorized).
			Public("Invalid or expired refresh token").
			Errorf("user_id not found in access token claims")
	}

	memberCriteria := dafi.Where("userId", dafi.Equal, userID)
	member, err := s.workspaceMemberSvc.FindOne(ctx, memberCriteria)
	if err != nil {
		return port.RefreshResponse{}, oops.WithContext(ctx).
			In(apperrors.LayerService).
			Code(apperrors.CodeUnauthorized).
			Public("Invalid or expired refresh token").
			With("user_id", userID).
			Errorf("user does not belong to any workspace: %s", userID)
	}

	refreshClaims, err := s.parseToken(refreshToken, false)
	if err != nil {
		return port.RefreshResponse{}, oops.WithContext(ctx).
			In(apperrors.LayerService).
			Code(apperrors.CodeUnauthorized).
			Public("Invalid or expired refresh token").
			Errorf("failed to parse refresh token: %w", err)
	}

	tokenType, ok := refreshClaims["type"].(string)
	if !ok || tokenType != "refresh" {
		return port.RefreshResponse{}, oops.WithContext(ctx).
			In(apperrors.LayerService).
			Code(apperrors.CodeUnauthorized).
			Public("Invalid or expired refresh token").
			Errorf("invalid token type: expected 'refresh', got '%s'", tokenType)
	}

	refreshUserID, ok := refreshClaims["user_id"].(string)
	if !ok || refreshUserID == "" {
		return port.RefreshResponse{}, oops.WithContext(ctx).
			In(apperrors.LayerService).
			Code(apperrors.CodeUnauthorized).
			Public("Invalid or expired refresh token").
			Errorf("user_id not found in refresh token claims")
	}

	if refreshUserID != userID {
		return port.RefreshResponse{}, oops.WithContext(ctx).
			In(apperrors.LayerService).
			Code(apperrors.CodeUnauthorized).
			Public("Invalid or expired refresh token").
			Errorf("user_id mismatch between access and refresh tokens")
	}

	workspaceID, ok := accessClaims["workspace_id"].(string)
	if !ok || workspaceID == "" {
		workspaceID = member.WorkspaceID.String()
	}

	tokenPair, err := s.generateTokenPair(userID, workspaceID)
	if err != nil {
		return port.RefreshResponse{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("token refreshed",
		"user_id", userID,
	)

	return port.RefreshResponse{
		TokenPair: tokenPair,
	}, nil
}

func (s service) parseToken(tokenString string, allowExpired bool) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, oops.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		if allowExpired && errors.Is(err, jwt.ErrTokenExpired) {
			token, _ = jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, oops.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return s.jwtSecret, nil
			})
		} else {
			return nil, err
		}
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, oops.Errorf("invalid token claims")
	}

	return *claims, nil
}
