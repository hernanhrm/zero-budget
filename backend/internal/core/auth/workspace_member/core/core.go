package core

import (
	"context"
	"slices"

	roleport "backend/core/auth/role/port"
	userport "backend/core/auth/user/port"
	"backend/core/auth/workspace_member/port"
	"backend/infra/dafi"
	basedomain "backend/port"
	apperrors "backend/port/errors"
	"github.com/google/uuid"
	"github.com/samber/oops"
)

type service struct {
	repo    port.Repository
	userSvc basedomain.UseCaseFindAll[basedomain.List[userport.User]]
	roleSvc basedomain.UseCaseFindAll[basedomain.List[roleport.Role]]
	logger  basedomain.Logger
}

func New(
	repo port.Repository,
	userSvc basedomain.UseCaseFindAll[basedomain.List[userport.User]],
	roleSvc basedomain.UseCaseFindAll[basedomain.List[roleport.Role]],
	logger basedomain.Logger,
) port.Service {
	return service{
		repo:    repo,
		userSvc: userSvc,
		roleSvc: roleSvc,
		logger:  logger.With("component", "workspace_member.service"),
	}
}

func (s service) WithTx(tx basedomain.Transaction) port.Service {
	return service{
		repo:    s.repo.WithTx(tx),
		userSvc: s.userSvc,
		roleSvc: s.roleSvc,
		logger:  s.logger,
	}
}

func (s service) FindOne(ctx context.Context, criteria dafi.Criteria) (port.WorkspaceMember, error) {
	member, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return port.WorkspaceMember{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return member, nil
}

func (s service) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[port.WorkspaceMember], error) {
	members, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return members, nil
}

func (s service) FindOneRelation(ctx context.Context, criteria dafi.Criteria) (port.WorkspaceMemberRelation, error) {
	member, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return port.WorkspaceMemberRelation{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	result := port.WorkspaceMemberRelation{WorkspaceMember: member}

	if slices.Contains(criteria.Relations, "user") {
		userCriteria := dafi.Where("id", dafi.Equal, member.UserID.String())
		users, err := s.userSvc.FindAll(ctx, userCriteria)
		if err != nil {
			return port.WorkspaceMemberRelation{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
		}
		if len(users) > 0 {
			result.User = &users[0]
		}
	}

	if slices.Contains(criteria.Relations, "role") {
		roleCriteria := dafi.Where("id", dafi.Equal, member.RoleID.String())
		roles, err := s.roleSvc.FindAll(ctx, roleCriteria)
		if err != nil {
			return port.WorkspaceMemberRelation{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
		}
		if len(roles) > 0 {
			result.Role = &roles[0]
		}
	}

	return result, nil
}

func (s service) FindAllRelation(ctx context.Context, criteria dafi.Criteria) ([]port.WorkspaceMemberRelation, error) {
	members, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	result := make([]port.WorkspaceMemberRelation, len(members))
	for i, m := range members {
		result[i] = port.WorkspaceMemberRelation{WorkspaceMember: m}
	}

	if len(members) == 0 {
		return result, nil
	}

	if slices.Contains(criteria.Relations, "user") {
		userIDs := uniqueIDs(members, func(m port.WorkspaceMember) uuid.UUID { return m.UserID })
		userCriteria := dafi.Where("id", dafi.In, userIDs)
		users, err := s.userSvc.FindAll(ctx, userCriteria)
		if err != nil {
			return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
		}

		userMap := make(map[string]*userport.User, len(users))
		for i := range users {
			userMap[users[i].ID] = &users[i]
		}

		for i := range result {
			if u, ok := userMap[result[i].UserID.String()]; ok {
				result[i].User = u
			}
		}
	}

	if slices.Contains(criteria.Relations, "role") {
		roleIDs := uniqueIDs(members, func(m port.WorkspaceMember) uuid.UUID { return m.RoleID })
		roleCriteria := dafi.Where("id", dafi.In, roleIDs)
		roles, err := s.roleSvc.FindAll(ctx, roleCriteria)
		if err != nil {
			return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
		}

		roleMap := make(map[uuid.UUID]*roleport.Role, len(roles))
		for i := range roles {
			roleMap[roles[i].ID] = &roles[i]
		}

		for i := range result {
			if r, ok := roleMap[result[i].RoleID]; ok {
				result[i].Role = r
			}
		}
	}

	return result, nil
}

func (s service) Create(ctx context.Context, input port.CreateWorkspaceMember) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	if err := s.repo.Create(ctx, input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("workspace member created")

	return nil
}

func (s service) CreateBulk(ctx context.Context, inputs basedomain.List[port.CreateWorkspaceMember]) error {
	for _, input := range inputs {
		if err := input.Validate(ctx); err != nil {
			return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
		}
	}

	if err := s.repo.CreateBulk(ctx, inputs); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("workspace members created", "count", len(inputs))

	return nil
}

func (s service) Update(ctx context.Context, input port.UpdateWorkspaceMember, filters ...dafi.Filter) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	if err := s.repo.Update(ctx, input, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("workspace member updated")

	return nil
}

func (s service) Delete(ctx context.Context, filters ...dafi.Filter) error {
	if err := s.repo.Delete(ctx, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("workspace member deleted")

	return nil
}

func uniqueIDs(members basedomain.List[port.WorkspaceMember], getID func(port.WorkspaceMember) uuid.UUID) []string {
	seen := make(map[uuid.UUID]struct{}, len(members))
	ids := make([]string, 0, len(members))

	for _, m := range members {
		id := getID(m)
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
			ids = append(ids, id.String())
		}
	}

	return ids
}
