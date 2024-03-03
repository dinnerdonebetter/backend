package postgres

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.ValidIngredientGroupDataManager = (*Querier)(nil)
)

// ValidIngredientGroupExists fetches whether a valid ingredient group exists from the database.
func (q *Querier) ValidIngredientGroupExists(ctx context.Context, validIngredientGroupID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientGroupID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientGroupIDKey, validIngredientGroupID)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, validIngredientGroupID)

	result, err := q.generatedQuerier.CheckValidIngredientGroupExistence(ctx, q.db, validIngredientGroupID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient group existence check")
	}

	return result, nil
}

// GetValidIngredientGroup fetches a valid ingredient group from the database.
func (q *Querier) GetValidIngredientGroup(ctx context.Context, validIngredientGroupID string) (*types.ValidIngredientGroup, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientGroupID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientGroupIDKey, validIngredientGroupID)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, validIngredientGroupID)

	result, err := q.generatedQuerier.GetValidIngredientGroup(ctx, q.db, validIngredientGroupID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredients group from database")
	}

	validIngredientGroup := &types.ValidIngredientGroup{
		CreatedAt:     result.CreatedAt,
		LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
		ID:            result.ID,
		Name:          result.Name,
		Slug:          result.Slug,
		Description:   result.Description,
		Members:       nil,
	}

	membersResults, err := q.generatedQuerier.GetValidIngredientGroupMembers(ctx, q.db, validIngredientGroupID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredients group members from database")
	}

	for _, memberResult := range membersResults {
		validIngredientGroup.Members = append(validIngredientGroup.Members, &types.ValidIngredientGroupMember{
			CreatedAt:      memberResult.CreatedAt,
			ArchivedAt:     database.TimePointerFromNullTime(memberResult.ArchivedAt),
			ID:             memberResult.ID,
			BelongsToGroup: memberResult.BelongsToGroup,
			ValidIngredient: types.ValidIngredient{
				CreatedAt:                               memberResult.ValidIngredientCreatedAt,
				LastUpdatedAt:                           database.TimePointerFromNullTime(memberResult.ValidIngredientLastUpdatedAt),
				ArchivedAt:                              database.TimePointerFromNullTime(memberResult.ValidIngredientArchivedAt),
				MaximumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(memberResult.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
				MinimumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(memberResult.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
				IconPath:                                memberResult.ValidIngredientIconPath,
				Warning:                                 memberResult.ValidIngredientWarning,
				PluralName:                              memberResult.ValidIngredientPluralName,
				StorageInstructions:                     memberResult.ValidIngredientStorageInstructions,
				Name:                                    memberResult.ValidIngredientName,
				ID:                                      memberResult.ValidIngredientID,
				Description:                             memberResult.ValidIngredientDescription,
				Slug:                                    memberResult.ValidIngredientSlug,
				ShoppingSuggestions:                     memberResult.ValidIngredientShoppingSuggestions,
				ContainsShellfish:                       memberResult.ValidIngredientContainsShellfish,
				IsMeasuredVolumetrically:                memberResult.ValidIngredientVolumetric,
				IsLiquid:                                database.BoolFromNullBool(memberResult.ValidIngredientIsLiquid),
				ContainsPeanut:                          memberResult.ValidIngredientContainsPeanut,
				ContainsTreeNut:                         memberResult.ValidIngredientContainsTreeNut,
				ContainsEgg:                             memberResult.ValidIngredientContainsEgg,
				ContainsWheat:                           memberResult.ValidIngredientContainsWheat,
				ContainsSoy:                             memberResult.ValidIngredientContainsSoy,
				AnimalDerived:                           memberResult.ValidIngredientAnimalDerived,
				RestrictToPreparations:                  memberResult.ValidIngredientRestrictToPreparations,
				ContainsSesame:                          memberResult.ValidIngredientContainsSesame,
				ContainsFish:                            memberResult.ValidIngredientContainsFish,
				ContainsGluten:                          memberResult.ValidIngredientContainsGluten,
				ContainsDairy:                           memberResult.ValidIngredientContainsDairy,
				ContainsAlcohol:                         memberResult.ValidIngredientContainsAlcohol,
				AnimalFlesh:                             memberResult.ValidIngredientAnimalFlesh,
				IsStarch:                                memberResult.ValidIngredientIsStarch,
				IsProtein:                               memberResult.ValidIngredientIsProtein,
				IsGrain:                                 memberResult.ValidIngredientIsGrain,
				IsFruit:                                 memberResult.ValidIngredientIsFruit,
				IsSalt:                                  memberResult.ValidIngredientIsSalt,
				IsFat:                                   memberResult.ValidIngredientIsFat,
				IsAcid:                                  memberResult.ValidIngredientIsAcid,
				IsHeat:                                  memberResult.ValidIngredientIsHeat,
			},
		})
	}

	return validIngredientGroup, nil
}

// SearchForValidIngredientGroups fetches a valid ingredient group from the database.
func (q *Querier) SearchForValidIngredientGroups(ctx context.Context, query string, filter *types.QueryFilter) ([]*types.ValidIngredientGroup, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, query)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)
	filter.AttachToLogger(logger)

	results, err := q.generatedQuerier.SearchForValidIngredientGroups(ctx, q.db, &generated.SearchForValidIngredientGroupsParams{
		Name:          query,
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching webhook from database")
	}

	validIngredientGroups := []*types.ValidIngredientGroup{}
	for _, result := range results {
		validIngredientGroup := &types.ValidIngredientGroup{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			ID:            result.ID,
			Name:          result.Name,
			Slug:          result.Slug,
			Description:   result.Description,
			Members:       []*types.ValidIngredientGroupMember{},
		}

		var membersResults []*generated.GetValidIngredientGroupMembersRow
		membersResults, err = q.generatedQuerier.GetValidIngredientGroupMembers(ctx, q.db, result.ID)
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredients group members from database")
		}

		for _, memberResult := range membersResults {
			validIngredientGroup.Members = append(validIngredientGroup.Members, &types.ValidIngredientGroupMember{
				CreatedAt:      memberResult.CreatedAt,
				ArchivedAt:     database.TimePointerFromNullTime(memberResult.ArchivedAt),
				ID:             memberResult.ID,
				BelongsToGroup: memberResult.BelongsToGroup,
				ValidIngredient: types.ValidIngredient{
					CreatedAt:                               memberResult.ValidIngredientCreatedAt,
					LastUpdatedAt:                           database.TimePointerFromNullTime(memberResult.ValidIngredientLastUpdatedAt),
					ArchivedAt:                              database.TimePointerFromNullTime(memberResult.ValidIngredientArchivedAt),
					MaximumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(memberResult.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
					MinimumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(memberResult.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
					IconPath:                                memberResult.ValidIngredientIconPath,
					Warning:                                 memberResult.ValidIngredientWarning,
					PluralName:                              memberResult.ValidIngredientPluralName,
					StorageInstructions:                     memberResult.ValidIngredientStorageInstructions,
					Name:                                    memberResult.ValidIngredientName,
					ID:                                      memberResult.ValidIngredientID,
					Description:                             memberResult.ValidIngredientDescription,
					Slug:                                    memberResult.ValidIngredientSlug,
					ShoppingSuggestions:                     memberResult.ValidIngredientShoppingSuggestions,
					ContainsShellfish:                       memberResult.ValidIngredientContainsShellfish,
					IsMeasuredVolumetrically:                memberResult.ValidIngredientVolumetric,
					IsLiquid:                                database.BoolFromNullBool(memberResult.ValidIngredientIsLiquid),
					ContainsPeanut:                          memberResult.ValidIngredientContainsPeanut,
					ContainsTreeNut:                         memberResult.ValidIngredientContainsTreeNut,
					ContainsEgg:                             memberResult.ValidIngredientContainsEgg,
					ContainsWheat:                           memberResult.ValidIngredientContainsWheat,
					ContainsSoy:                             memberResult.ValidIngredientContainsSoy,
					AnimalDerived:                           memberResult.ValidIngredientAnimalDerived,
					RestrictToPreparations:                  memberResult.ValidIngredientRestrictToPreparations,
					ContainsSesame:                          memberResult.ValidIngredientContainsSesame,
					ContainsFish:                            memberResult.ValidIngredientContainsFish,
					ContainsGluten:                          memberResult.ValidIngredientContainsGluten,
					ContainsDairy:                           memberResult.ValidIngredientContainsDairy,
					ContainsAlcohol:                         memberResult.ValidIngredientContainsAlcohol,
					AnimalFlesh:                             memberResult.ValidIngredientAnimalFlesh,
					IsStarch:                                memberResult.ValidIngredientIsStarch,
					IsProtein:                               memberResult.ValidIngredientIsProtein,
					IsGrain:                                 memberResult.ValidIngredientIsGrain,
					IsFruit:                                 memberResult.ValidIngredientIsFruit,
					IsSalt:                                  memberResult.ValidIngredientIsSalt,
					IsFat:                                   memberResult.ValidIngredientIsFat,
					IsAcid:                                  memberResult.ValidIngredientIsAcid,
					IsHeat:                                  memberResult.ValidIngredientIsHeat,
				},
			})
		}

		validIngredientGroups = append(validIngredientGroups, validIngredientGroup)
	}

	return validIngredientGroups, nil
}

// GetValidIngredientGroups fetches a list of valid ingredients group from the database that meet a particular filter.
func (q *Querier) GetValidIngredientGroups(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidIngredientGroup], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ValidIngredientGroup]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetValidIngredientGroups(ctx, q.db, &generated.GetValidIngredientGroupsParams{
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching webhook from database")
	}

	for _, result := range results {
		validIngredientGroup := &types.ValidIngredientGroup{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			ID:            result.ID,
			Name:          result.Name,
			Slug:          result.Slug,
			Description:   result.Description,
			Members:       []*types.ValidIngredientGroupMember{},
		}

		var membersResults []*generated.GetValidIngredientGroupMembersRow
		membersResults, err = q.generatedQuerier.GetValidIngredientGroupMembers(ctx, q.db, result.ID)
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredients group members from database")
		}

		for _, memberResult := range membersResults {
			validIngredientGroup.Members = append(validIngredientGroup.Members, &types.ValidIngredientGroupMember{
				CreatedAt:      memberResult.CreatedAt,
				ArchivedAt:     database.TimePointerFromNullTime(memberResult.ArchivedAt),
				ID:             memberResult.ID,
				BelongsToGroup: memberResult.BelongsToGroup,
				ValidIngredient: types.ValidIngredient{
					CreatedAt:                               memberResult.ValidIngredientCreatedAt,
					LastUpdatedAt:                           database.TimePointerFromNullTime(memberResult.ValidIngredientLastUpdatedAt),
					ArchivedAt:                              database.TimePointerFromNullTime(memberResult.ValidIngredientArchivedAt),
					MaximumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(memberResult.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
					MinimumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(memberResult.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
					IconPath:                                memberResult.ValidIngredientIconPath,
					Warning:                                 memberResult.ValidIngredientWarning,
					PluralName:                              memberResult.ValidIngredientPluralName,
					StorageInstructions:                     memberResult.ValidIngredientStorageInstructions,
					Name:                                    memberResult.ValidIngredientName,
					ID:                                      memberResult.ValidIngredientID,
					Description:                             memberResult.ValidIngredientDescription,
					Slug:                                    memberResult.ValidIngredientSlug,
					ShoppingSuggestions:                     memberResult.ValidIngredientShoppingSuggestions,
					ContainsShellfish:                       memberResult.ValidIngredientContainsShellfish,
					IsMeasuredVolumetrically:                memberResult.ValidIngredientVolumetric,
					IsLiquid:                                database.BoolFromNullBool(memberResult.ValidIngredientIsLiquid),
					ContainsPeanut:                          memberResult.ValidIngredientContainsPeanut,
					ContainsTreeNut:                         memberResult.ValidIngredientContainsTreeNut,
					ContainsEgg:                             memberResult.ValidIngredientContainsEgg,
					ContainsWheat:                           memberResult.ValidIngredientContainsWheat,
					ContainsSoy:                             memberResult.ValidIngredientContainsSoy,
					AnimalDerived:                           memberResult.ValidIngredientAnimalDerived,
					RestrictToPreparations:                  memberResult.ValidIngredientRestrictToPreparations,
					ContainsSesame:                          memberResult.ValidIngredientContainsSesame,
					ContainsFish:                            memberResult.ValidIngredientContainsFish,
					ContainsGluten:                          memberResult.ValidIngredientContainsGluten,
					ContainsDairy:                           memberResult.ValidIngredientContainsDairy,
					ContainsAlcohol:                         memberResult.ValidIngredientContainsAlcohol,
					AnimalFlesh:                             memberResult.ValidIngredientAnimalFlesh,
					IsStarch:                                memberResult.ValidIngredientIsStarch,
					IsProtein:                               memberResult.ValidIngredientIsProtein,
					IsGrain:                                 memberResult.ValidIngredientIsGrain,
					IsFruit:                                 memberResult.ValidIngredientIsFruit,
					IsSalt:                                  memberResult.ValidIngredientIsSalt,
					IsFat:                                   memberResult.ValidIngredientIsFat,
					IsAcid:                                  memberResult.ValidIngredientIsAcid,
					IsHeat:                                  memberResult.ValidIngredientIsHeat,
				},
			})
		}

		x.Data = append(x.Data, validIngredientGroup)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// CreateValidIngredientGroup creates a valid ingredient group in the database.
func (q *Querier) CreateValidIngredientGroup(ctx context.Context, input *types.ValidIngredientGroupDatabaseCreationInput) (*types.ValidIngredientGroup, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, input.ID)
	logger := q.logger.WithValue(keys.ValidIngredientGroupIDKey, input.ID)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "starting transaction")
	}

	// create the valid ingredient group.
	if err = q.generatedQuerier.CreateValidIngredientGroup(ctx, tx, &generated.CreateValidIngredientGroupParams{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		Slug:        input.Slug,
	}); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient group creation query")
	}

	x := &types.ValidIngredientGroup{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		Slug:        input.Slug,
		CreatedAt:   q.currentTime(),
	}

	for i := range input.Members {
		m := input.Members[i]
		var member *types.ValidIngredientGroupMember
		member, err = q.CreateValidIngredientGroupMember(ctx, tx, x.ID, m)
		if err != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient group member")
		}

		x.Members = append(x.Members, member)
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.WithValue("member_count", len(input.Members)).Info("valid ingredient group created")

	return x, nil
}

// CreateValidIngredientGroupMember creates a valid ingredient group member in the database.
func (q *Querier) CreateValidIngredientGroupMember(ctx context.Context, db database.SQLQueryExecutorAndTransactionManager, groupID string, input *types.ValidIngredientGroupMemberDatabaseCreationInput) (*types.ValidIngredientGroupMember, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidIngredientGroupIDKey, input.ID).WithValue(keys.ValidIngredientIDKey, input.ValidIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, input.ID)

	// create the valid ingredient group.
	if err := q.generatedQuerier.CreateValidIngredientGroupMember(ctx, db, &generated.CreateValidIngredientGroupMemberParams{
		ID:              input.ID,
		BelongsToGroup:  groupID,
		ValidIngredient: input.ValidIngredientID,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient group member creation query")
	}

	x := &types.ValidIngredientGroupMember{
		ID:              input.ID,
		BelongsToGroup:  groupID,
		ValidIngredient: types.ValidIngredient{ID: input.ValidIngredientID},
		CreatedAt:       q.currentTime(),
	}

	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, x.ID)
	logger.Info("valid ingredient group member created")

	return x, nil
}

// UpdateValidIngredientGroup updates a particular valid ingredient group.
func (q *Querier) UpdateValidIngredientGroup(ctx context.Context, updated *types.ValidIngredientGroup) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.ValidIngredientGroupIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateValidIngredientGroup(ctx, q.db, &generated.UpdateValidIngredientGroupParams{
		Name:        updated.Name,
		Description: updated.Description,
		Slug:        updated.Slug,
		ID:          updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient group")
	}

	logger.Info("valid ingredient group updated")

	return nil
}

// ArchiveValidIngredientGroup archives a valid ingredient group from the database by its ID.
func (q *Querier) ArchiveValidIngredientGroup(ctx context.Context, validIngredientGroupID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientGroupID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientGroupIDKey, validIngredientGroupID)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, validIngredientGroupID)

	if _, err := q.generatedQuerier.ArchiveValidIngredientGroup(ctx, q.db, validIngredientGroupID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient group")
	}

	logger.Info("valid ingredient group archived")

	return nil
}
