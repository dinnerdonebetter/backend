package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	_ types.HouseholdInvitationDataManager = (*SQLQuerier)(nil)

	// householdInvitationsTableColumns are the columns for the household invitations table.
	householdInvitationsTableColumns = []string{
		"household_invitations.id",
		"household_invitations.destination_household",
		"household_invitations.to_email",
		"household_invitations.to_user",
		"household_invitations.from_user",
		"household_invitations.status",
		"household_invitations.note",
		"household_invitations.status_note",
		"household_invitations.token",
		"household_invitations.created_on",
		"household_invitations.last_updated_on",
		"household_invitations.archived_on",
	}
)

// scanHouseholdInvitation is a consistent way to turn a *sql.Row into an invitation struct.
func (q *SQLQuerier) scanHouseholdInvitation(ctx context.Context, scan database.Scanner, includeCounts bool) (householdInvitation *types.HouseholdInvitation, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)
	householdInvitation = &types.HouseholdInvitation{}

	targetVars := []interface{}{
		&householdInvitation.ID,
		&householdInvitation.DestinationHousehold,
		&householdInvitation.ToEmail,
		&householdInvitation.ToUser,
		&householdInvitation.FromUser,
		&householdInvitation.Status,
		&householdInvitation.Note,
		&householdInvitation.StatusNote,
		&householdInvitation.Token,
		&householdInvitation.CreatedOn,
		&householdInvitation.LastUpdatedOn,
		&householdInvitation.ArchivedOn,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "scanning householdInvitation")
	}

	return householdInvitation, filteredCount, totalCount, nil
}

// scanHouseholdInvitations provides a consistent way to turn sql rows into a slice of household_invitations.
func (q *SQLQuerier) scanHouseholdInvitations(ctx context.Context, rows database.ResultIterator, includeCounts bool) (householdInvitations []*types.HouseholdInvitation, filteredCount, totalCount uint64, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		householdInvitation, fc, tc, scanErr := q.scanHouseholdInvitation(ctx, rows, includeCounts)
		if scanErr != nil {
			return nil, 0, 0, scanErr
		}

		if includeCounts {
			if filteredCount == 0 {
				filteredCount = fc
			}

			if totalCount == 0 {
				totalCount = tc
			}
		}

		householdInvitations = append(householdInvitations, householdInvitation)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "fetching household invitation from database")
	}

	if err = rows.Close(); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "fetching household invitation from database")
	}

	return householdInvitations, filteredCount, totalCount, nil
}

const householdInvitationExistenceQuery = "SELECT EXISTS ( SELECT household_invitations.id FROM household_invitations WHERE household_invitations.archived_on IS NULL AND household_invitations.id = $1 )"

// HouseholdInvitationExists fetches whether a household invitation exists from the database.
func (q *SQLQuerier) HouseholdInvitationExists(ctx context.Context, householdInvitationID string) (bool, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdInvitationID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)
	tracing.AttachHouseholdInvitationIDToSpan(span, householdInvitationID)

	args := []interface{}{
		householdInvitationID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, householdInvitationExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing household invitation existence check")
	}

	return result, nil
}

const getHouseholdInvitationByHouseholdAndIDQuery = `
SELECT
	household_invitations.id,
	household_invitations.destination_household,
	household_invitations.to_email,
	household_invitations.to_user,
	household_invitations.from_user,
	household_invitations.status,
	household_invitations.note,
	household_invitations.status_note,
	household_invitations.token,
	household_invitations.created_on,
	household_invitations.last_updated_on,
	household_invitations.archived_on
FROM household_invitations 
WHERE household_invitations.archived_on IS NULL
AND household_invitations.destination_household = $1
AND household_invitations.id = $2
`

// GetHouseholdInvitationByHouseholdAndID fetches an invitation from the database.
func (q *SQLQuerier) GetHouseholdInvitationByHouseholdAndID(ctx context.Context, householdID, householdInvitationID string) (*types.HouseholdInvitation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	if householdInvitationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)
	tracing.AttachHouseholdInvitationIDToSpan(span, householdInvitationID)

	args := []interface{}{
		householdID,
		householdInvitationID,
	}

	row := q.getOneRow(ctx, q.db, "household invitation", getHouseholdInvitationByHouseholdAndIDQuery, args)

	householdInvitation, _, _, err := q.scanHouseholdInvitation(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning household invitation")
	}

	return householdInvitation, nil
}

/* #nosec */
const getHouseholdInvitationByEmailAndTokenQuery = `
SELECT
	household_invitations.id,
	household_invitations.destination_household,
	household_invitations.to_email,
	household_invitations.to_user,
	household_invitations.from_user,
	household_invitations.status,
	household_invitations.note,
	household_invitations.status_note,
	household_invitations.token,
	household_invitations.created_on,
	household_invitations.last_updated_on,
	household_invitations.archived_on
FROM household_invitations 
WHERE household_invitations.archived_on IS NULL 
AND household_invitations.to_email = $1
AND household_invitations.token = $2
`

// GetHouseholdInvitationByEmailAndToken fetches an invitation from the database.
func (q *SQLQuerier) GetHouseholdInvitationByEmailAndToken(ctx context.Context, emailAddress, token string) (*types.HouseholdInvitation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if emailAddress == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserEmailAddressKey, emailAddress)
	tracing.AttachEmailAddressToSpan(span, emailAddress)

	if token == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationTokenKey, token)
	tracing.AttachHouseholdInvitationTokenToSpan(span, token)

	args := []interface{}{
		emailAddress,
		token,
	}

	row := q.getOneRow(ctx, q.db, "household invitation", getHouseholdInvitationByEmailAndTokenQuery, args)

	invitation, _, _, err := q.scanHouseholdInvitation(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning invitation")
	}

	return invitation, nil
}

const getAllHouseholdInvitationsCountQuery = `
	SELECT COUNT(household_invitations.id) FROM household_invitations WHERE household_invitations.archived_on IS NULL
`

// GetAllHouseholdInvitationsCount fetches the count of household invitations from the database that meet a particular filter.
func (q *SQLQuerier) GetAllHouseholdInvitationsCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	count, err := q.performCountQuery(ctx, q.db, getAllHouseholdInvitationsCountQuery, "fetching count of household invitations")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of household invitations")
	}

	return count, nil
}

const createHouseholdInvitationQuery = `
	INSERT INTO household_invitations (id,from_user,to_user,note,to_email,token,destination_household) VALUES ($1,$2,$3,$4,$5,$6,$7)
`

// CreateHouseholdInvitation creates an invitation in a database.
func (q *SQLQuerier) CreateHouseholdInvitation(ctx context.Context, input *types.HouseholdInvitationDatabaseCreationInput) (*types.HouseholdInvitation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.HouseholdInvitationIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.FromUser,
		input.ToUser,
		input.Note,
		input.ToEmail,
		input.Token,
		input.DestinationHousehold,
	}

	if err := q.performWriteQuery(ctx, q.db, "household invitation creation", createHouseholdInvitationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "performing household invitation creation query")
	}

	x := &types.HouseholdInvitation{
		ID:                   input.ID,
		FromUser:             input.FromUser,
		ToUser:               input.ToUser,
		Note:                 input.Note,
		ToEmail:              input.ToEmail,
		Token:                input.Token,
		StatusNote:           "",
		Status:               types.PendingHouseholdInvitationStatus,
		DestinationHousehold: input.DestinationHousehold,
		CreatedOn:            q.currentTime(),
	}

	tracing.AttachHouseholdInvitationIDToSpan(span, x.ID)
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, x.ID)

	logger.Info("household invitation created")

	return x, nil
}

// BuildGetPendingHouseholdInvitationsFromUserQuery builds a query for fetching pending household invitations sent by a given user.
func (q *SQLQuerier) BuildGetPendingHouseholdInvitationsFromUserQuery(ctx context.Context, userID string, filter *types.QueryFilter) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	queryBuilder := q.sqlBuilder.Select(householdInvitationsTableColumns...).
		From("household_invitations").
		Where(squirrel.Eq{
			"household_invitations.from_user":   userID,
			"household_invitations.archived_on": nil,
			"household_invitations.status":      types.PendingHouseholdInvitationStatus,
		})

	queryBuilder = applyFilterToQueryBuilder(filter, "household_invitations", queryBuilder)

	query, args, err := queryBuilder.ToSql()
	q.logQueryBuildingError(span, err)

	return query, args
}

// GetPendingHouseholdInvitationsFromUser fetches pending household invitations sent from a given user.
func (q *SQLQuerier) GetPendingHouseholdInvitationsFromUser(ctx context.Context, userID string, filter *types.QueryFilter) (*types.HouseholdInvitationList, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue(keys.UserIDKey, userID)
	filter.AttachToLogger(logger)

	query, args := q.BuildGetPendingHouseholdInvitationsFromUserQuery(ctx, userID, filter)

	rows, err := q.performReadQuery(ctx, q.db, "household invitations from user", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "reading household invitations from user")
	}

	householdInvitations, fc, tc, err := q.scanHouseholdInvitations(ctx, rows, true)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "reading household invitations from user")
	}

	returnList := &types.HouseholdInvitationList{
		Pagination: types.Pagination{
			Page:          filter.Page,
			Limit:         filter.Limit,
			FilteredCount: fc,
			TotalCount:    tc,
		},
		HouseholdInvitations: householdInvitations,
	}

	return returnList, nil
}

// BuildGetPendingHouseholdInvitationsForUserQuery builds a query for fetching pending household invitations sent to a given user.
func (q *SQLQuerier) BuildGetPendingHouseholdInvitationsForUserQuery(ctx context.Context, userID string, filter *types.QueryFilter) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	queryBuilder := q.sqlBuilder.Select(householdInvitationsTableColumns...).
		From("household_invitations").
		Where(squirrel.Eq{
			"household_invitations.to_user":     userID,
			"household_invitations.archived_on": nil,
			"household_invitations.status":      types.PendingHouseholdInvitationStatus,
		})

	queryBuilder = applyFilterToQueryBuilder(filter, "household_invitations", queryBuilder)

	query, args, err := queryBuilder.ToSql()
	q.logQueryBuildingError(span, err)

	return query, args
}

// GetPendingHouseholdInvitationsForUser fetches pending household invitations sent to a given user.
func (q *SQLQuerier) GetPendingHouseholdInvitationsForUser(ctx context.Context, userID string, filter *types.QueryFilter) (*types.HouseholdInvitationList, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue(keys.UserIDKey, userID)
	filter.AttachToLogger(logger)

	query, args := q.BuildGetPendingHouseholdInvitationsForUserQuery(ctx, userID, filter)

	rows, err := q.performReadQuery(ctx, q.db, "household invitations from user", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "reading household invitations from user")
	}

	householdInvitations, fc, tc, err := q.scanHouseholdInvitations(ctx, rows, true)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "reading household invitations from user")
	}

	returnList := &types.HouseholdInvitationList{
		Pagination: types.Pagination{
			Page:          filter.Page,
			Limit:         filter.Limit,
			FilteredCount: fc,
			TotalCount:    tc,
		},
		HouseholdInvitations: householdInvitations,
	}

	return returnList, nil
}

const setInvitationStatusQuery = `
UPDATE household_invitations SET
	status = $1,
	status_note = $2,
	last_updated_on = extract(epoch FROM NOW()), 
	archived_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL 
AND destination_household = $3
AND id = $4
`

func (q *SQLQuerier) setInvitationStatus(ctx context.Context, querier database.SQLQueryExecutor, householdID, householdInvitationID, note string, status types.HouseholdInvitationStatus) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("new_status", status)

	if householdID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	if householdInvitationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)
	tracing.AttachHouseholdInvitationIDToSpan(span, householdInvitationID)

	args := []interface{}{
		status,
		note,
		householdID,
		householdInvitationID,
	}

	if err := q.performWriteQuery(ctx, querier, "household invitation status change", setInvitationStatusQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "changing household invitation status")
	}

	logger.Debug("household invitation updated")

	return nil
}

// CancelHouseholdInvitation cancels a household invitation by its ID with a note.
func (q *SQLQuerier) CancelHouseholdInvitation(ctx context.Context, householdID, householdInvitationID, note string) error {
	return q.setInvitationStatus(ctx, q.db, householdID, householdInvitationID, note, types.CancelledHouseholdInvitationStatus)
}

// AcceptHouseholdInvitation accepts a household invitation by its ID with a note.
func (q *SQLQuerier) AcceptHouseholdInvitation(ctx context.Context, householdID, householdInvitationID, note string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	if householdInvitationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)
	tracing.AttachHouseholdInvitationIDToSpan(span, householdInvitationID)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	if err = q.setInvitationStatus(ctx, tx, householdID, householdInvitationID, note, types.AcceptedHouseholdInvitationStatus); err != nil {
		return observability.PrepareError(err, logger, span, "accepting household invitation")
	}

	invitation, err := q.GetHouseholdInvitationByHouseholdAndID(ctx, householdID, householdInvitationID)
	if err != nil {
		return observability.PrepareError(err, logger, span, "fetching household invitation")
	}

	if err = q.addUserToHousehold(ctx, tx, &types.HouseholdUserMembershipDatabaseCreationInput{
		ID:             ksuid.New().String(),
		Reason:         fmt.Sprintf("accepted household invitation %q", householdInvitationID),
		UserID:         *invitation.ToUser,
		HouseholdID:    householdID,
		HouseholdRoles: []string{"household_member"},
	}); err != nil {
		return observability.PrepareError(err, logger, span, "adding user to household")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	return nil
}

// RejectHouseholdInvitation rejects a household invitation by its ID with a note.
func (q *SQLQuerier) RejectHouseholdInvitation(ctx context.Context, householdID, householdInvitationID, note string) error {
	return q.setInvitationStatus(ctx, q.db, householdID, householdInvitationID, note, types.RejectedHouseholdInvitationStatus)
}
