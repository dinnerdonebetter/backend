package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/Masterminds/squirrel"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
)

const (
	recipeIterationsTableName = "recipe_iterations"
)

var (
	recipeIterationsTableColumns = []string{
		"id",
		"recipe_id",
		"end_difficulty_rating",
		"end_complexity_rating",
		"end_taste_rating",
		"end_overall_rating",
		"created_on",
		"updated_on",
		"archived_on",
		"belongs_to",
	}
)

// scanRecipeIteration takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe Iteration struct
func scanRecipeIteration(scan database.Scanner) (*models.RecipeIteration, error) {
	x := &models.RecipeIteration{}

	if err := scan.Scan(
		&x.ID,
		&x.RecipeID,
		&x.EndDifficultyRating,
		&x.EndComplexityRating,
		&x.EndTasteRating,
		&x.EndOverallRating,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsTo,
	); err != nil {
		return nil, err
	}

	return x, nil
}

// scanRecipeIterations takes a logger and some database rows and turns them into a slice of recipe iterations
func scanRecipeIterations(logger logging.Logger, rows *sql.Rows) ([]models.RecipeIteration, error) {
	var list []models.RecipeIteration

	for rows.Next() {
		x, err := scanRecipeIteration(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		logger.Error(closeErr, "closing database rows")
	}

	return list, nil
}

// buildGetRecipeIterationQuery constructs a SQL query for fetching a recipe iteration with a given ID belong to a user with a given ID.
func (s *Sqlite) buildGetRecipeIterationQuery(recipeIterationID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = s.sqlBuilder.
		Select(recipeIterationsTableColumns...).
		From(recipeIterationsTableName).
		Where(squirrel.Eq{
			"id":         recipeIterationID,
			"belongs_to": userID,
		}).ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// GetRecipeIteration fetches a recipe iteration from the sqlite database
func (s *Sqlite) GetRecipeIteration(ctx context.Context, recipeIterationID, userID uint64) (*models.RecipeIteration, error) {
	query, args := s.buildGetRecipeIterationQuery(recipeIterationID, userID)
	row := s.db.QueryRowContext(ctx, query, args...)
	return scanRecipeIteration(row)
}

// buildGetRecipeIterationCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of recipe iterations belonging to a given user that meet a given query
func (s *Sqlite) buildGetRecipeIterationCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := s.sqlBuilder.
		Select(CountQuery).
		From(recipeIterationsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
			"belongs_to":  userID,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	s.logQueryBuildingError(err)

	return query, args
}

// GetRecipeIterationCount will fetch the count of recipe iterations from the database that meet a particular filter and belong to a particular user.
func (s *Sqlite) GetRecipeIterationCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := s.buildGetRecipeIterationCountQuery(filter, userID)
	err = s.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allRecipeIterationsCountQueryBuilder sync.Once
	allRecipeIterationsCountQuery        string
)

// buildGetAllRecipeIterationsCountQuery returns a query that fetches the total number of recipe iterations in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (s *Sqlite) buildGetAllRecipeIterationsCountQuery() string {
	allRecipeIterationsCountQueryBuilder.Do(func() {
		var err error
		allRecipeIterationsCountQuery, _, err = s.sqlBuilder.
			Select(CountQuery).
			From(recipeIterationsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		s.logQueryBuildingError(err)
	})

	return allRecipeIterationsCountQuery
}

// GetAllRecipeIterationsCount will fetch the count of recipe iterations from the database
func (s *Sqlite) GetAllRecipeIterationsCount(ctx context.Context) (count uint64, err error) {
	err = s.db.QueryRowContext(ctx, s.buildGetAllRecipeIterationsCountQuery()).Scan(&count)
	return count, err
}

// buildGetRecipeIterationsQuery builds a SQL query selecting recipe iterations that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (s *Sqlite) buildGetRecipeIterationsQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := s.sqlBuilder.
		Select(recipeIterationsTableColumns...).
		From(recipeIterationsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
			"belongs_to":  userID,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	s.logQueryBuildingError(err)

	return query, args
}

// GetRecipeIterations fetches a list of recipe iterations from the database that meet a particular filter
func (s *Sqlite) GetRecipeIterations(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeIterationList, error) {
	query, args := s.buildGetRecipeIterationsQuery(filter, userID)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe iterations")
	}

	list, err := scanRecipeIterations(s.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := s.GetRecipeIterationCount(ctx, filter, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching recipe iteration count: %w", err)
	}

	x := &models.RecipeIterationList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		RecipeIterations: list,
	}

	return x, nil
}

// GetAllRecipeIterationsForUser fetches every recipe iteration belonging to a user
func (s *Sqlite) GetAllRecipeIterationsForUser(ctx context.Context, userID uint64) ([]models.RecipeIteration, error) {
	query, args := s.buildGetRecipeIterationsQuery(nil, userID)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching recipe iterations for user")
	}

	list, err := scanRecipeIterations(s.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreateRecipeIterationQuery takes a recipe iteration and returns a creation query for that recipe iteration and the relevant arguments.
func (s *Sqlite) buildCreateRecipeIterationQuery(input *models.RecipeIteration) (query string, args []interface{}) {
	var err error
	query, args, err = s.sqlBuilder.
		Insert(recipeIterationsTableName).
		Columns(
			"recipe_id",
			"end_difficulty_rating",
			"end_complexity_rating",
			"end_taste_rating",
			"end_overall_rating",
			"belongs_to",
		).
		Values(
			input.RecipeID,
			input.EndDifficultyRating,
			input.EndComplexityRating,
			input.EndTasteRating,
			input.EndOverallRating,
			input.BelongsTo,
		).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// buildRecipeIterationCreationTimeQuery takes a recipe iteration and returns a creation query for that recipe iteration and the relevant arguments
func (s *Sqlite) buildRecipeIterationCreationTimeQuery(recipeIterationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select("created_on").
		From(recipeIterationsTableName).
		Where(squirrel.Eq{"id": recipeIterationID}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// CreateRecipeIteration creates a recipe iteration in the database
func (s *Sqlite) CreateRecipeIteration(ctx context.Context, input *models.RecipeIterationCreationInput) (*models.RecipeIteration, error) {
	x := &models.RecipeIteration{
		RecipeID:            input.RecipeID,
		EndDifficultyRating: input.EndDifficultyRating,
		EndComplexityRating: input.EndComplexityRating,
		EndTasteRating:      input.EndTasteRating,
		EndOverallRating:    input.EndOverallRating,
		BelongsTo:           input.BelongsTo,
	}

	query, args := s.buildCreateRecipeIterationQuery(x)

	// create the recipe iteration
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing recipe iteration creation query: %w", err)
	}

	// fetch the last inserted ID
	id, idErr := res.LastInsertId()
	if idErr == nil {
		x.ID = uint64(id)

		query, args := s.buildRecipeIterationCreationTimeQuery(x.ID)
		s.logCreationTimeRetrievalError(s.db.QueryRowContext(ctx, query, args...).Scan(&x.CreatedOn))
	}

	return x, nil
}

// buildUpdateRecipeIterationQuery takes a recipe iteration and returns an update SQL query, with the relevant query parameters
func (s *Sqlite) buildUpdateRecipeIterationQuery(input *models.RecipeIteration) (query string, args []interface{}) {
	var err error
	query, args, err = s.sqlBuilder.
		Update(recipeIterationsTableName).
		Set("recipe_id", input.RecipeID).
		Set("end_difficulty_rating", input.EndDifficultyRating).
		Set("end_complexity_rating", input.EndComplexityRating).
		Set("end_taste_rating", input.EndTasteRating).
		Set("end_overall_rating", input.EndOverallRating).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":         input.ID,
			"belongs_to": input.BelongsTo,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipeIteration updates a particular recipe iteration. Note that UpdateRecipeIteration expects the provided input to have a valid ID.
func (s *Sqlite) UpdateRecipeIteration(ctx context.Context, input *models.RecipeIteration) error {
	query, args := s.buildUpdateRecipeIterationQuery(input)
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveRecipeIterationQuery returns a SQL query which marks a given recipe iteration belonging to a given user as archived.
func (s *Sqlite) buildArchiveRecipeIterationQuery(recipeIterationID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = s.sqlBuilder.
		Update(recipeIterationsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          recipeIterationID,
			"archived_on": nil,
			"belongs_to":  userID,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// ArchiveRecipeIteration marks a recipe iteration as archived in the database
func (s *Sqlite) ArchiveRecipeIteration(ctx context.Context, recipeIterationID, userID uint64) error {
	query, args := s.buildArchiveRecipeIterationQuery(recipeIterationID, userID)
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}
