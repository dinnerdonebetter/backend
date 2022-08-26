package postgres

import (
	"fmt"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/pkg/types"
)

func TestQueryFilter_ApplyFilterToQueryBuilder(T *testing.T) {
	T.Parallel()

	exampleTableName := "stuff"
	baseQueryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("things").
		From(exampleTableName).
		Where(squirrel.Eq{fmt.Sprintf("%s.condition", exampleTableName): true})

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		qf := &types.QueryFilter{
			Page:          func(x uint64) *uint64 { return &x }(100),
			Limit:         func(x uint8) *uint8 { return &x }(50),
			CreatedAfter:  func(x uint64) *uint64 { return &x }(123456789),
			CreatedBefore: func(x uint64) *uint64 { return &x }(123456789),
			UpdatedAfter:  func(x uint64) *uint64 { return &x }(123456789),
			UpdatedBefore: func(x uint64) *uint64 { return &x }(123456789),
			SortBy:        types.SortDescending,
		}

		sb := squirrel.StatementBuilder.Select("*").From("testing")
		sb = applyFilterToQueryBuilder(qf, exampleTableName, sb)
		expected := "SELECT * FROM testing WHERE stuff.created_on > ? AND stuff.created_on < ? AND stuff.last_updated_on > ? AND stuff.last_updated_on < ? LIMIT 50 OFFSET 4950"
		actual, _, err := sb.ToSql()

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		sb := squirrel.StatementBuilder.Select("*").From("testing")
		sb = applyFilterToQueryBuilder(nil, exampleTableName, sb)
		expected := "SELECT * FROM testing"
		actual, _, err := sb.ToSql()

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	T.Run("basic usage", func(t *testing.T) {
		t.Parallel()

		qf := &types.QueryFilter{
			Limit: func(x uint8) *uint8 { return &x }(15),
			Page:  func(x uint64) *uint64 { return &x }(2),
		}

		expected := "SELECT things FROM stuff WHERE stuff.condition = $1 LIMIT 15 OFFSET 15"
		x := applyFilterToQueryBuilder(qf, exampleTableName, baseQueryBuilder)
		actual, args, err := x.ToSql()

		assert.Equal(t, expected, actual, "expected and actual queries don't match")
		assert.Nil(t, err)
		assert.NotEmpty(t, args)
	})

	T.Run("whole kit and kaboodle", func(t *testing.T) {
		t.Parallel()

		qf := &types.QueryFilter{
			Limit:         func(x uint8) *uint8 { return &x }(20),
			Page:          func(x uint64) *uint64 { return &x }(6),
			CreatedAfter:  func(x uint64) *uint64 { return &x }(uint64(time.Now().Unix())),
			CreatedBefore: func(x uint64) *uint64 { return &x }(uint64(time.Now().Unix())),
			UpdatedAfter:  func(x uint64) *uint64 { return &x }(uint64(time.Now().Unix())),
			UpdatedBefore: func(x uint64) *uint64 { return &x }(uint64(time.Now().Unix())),
		}

		expected := "SELECT things FROM stuff WHERE stuff.condition = $1 AND stuff.created_on > $2 AND stuff.created_on < $3 AND stuff.last_updated_on > $4 AND stuff.last_updated_on < $5 LIMIT 20 OFFSET 100"
		x := applyFilterToQueryBuilder(qf, exampleTableName, baseQueryBuilder)
		actual, args, err := x.ToSql()

		assert.Equal(t, expected, actual, "expected and actual queries don't match")
		assert.Nil(t, err)
		assert.NotEmpty(t, args)
	})

	T.Run("with zero limit", func(t *testing.T) {
		t.Parallel()

		qf := &types.QueryFilter{
			Limit: func(x uint8) *uint8 { return &x }(0),
			Page:  func(x uint64) *uint64 { return &x }(1),
		}
		expected := "SELECT things FROM stuff WHERE stuff.condition = $1 LIMIT 250"
		x := applyFilterToQueryBuilder(qf, exampleTableName, baseQueryBuilder)
		actual, args, err := x.ToSql()

		assert.Equal(t, expected, actual, "expected and actual queries don't match")
		assert.Nil(t, err)
		assert.NotEmpty(t, args)
	})
}

func TestQueryFilter_ApplyFilterToSubCountQueryBuilder(T *testing.T) {
	T.Parallel()

	exampleTableName := "stuff"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		qf := &types.QueryFilter{
			Page:          func(x uint64) *uint64 { return &x }(100),
			Limit:         func(x uint8) *uint8 { return &x }(50),
			CreatedAfter:  func(x uint64) *uint64 { return &x }(123456789),
			CreatedBefore: func(x uint64) *uint64 { return &x }(123456789),
			UpdatedAfter:  func(x uint64) *uint64 { return &x }(123456789),
			UpdatedBefore: func(x uint64) *uint64 { return &x }(123456789),
			SortBy:        types.SortDescending,
		}

		sb := squirrel.StatementBuilder.Select("*").From("testing")
		sb = applyFilterToSubCountQueryBuilder(qf, exampleTableName, sb)
		expected := "SELECT * FROM testing WHERE stuff.created_on > ? AND stuff.created_on < ? AND stuff.last_updated_on > ? AND stuff.last_updated_on < ?"
		actual, _, err := sb.ToSql()

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		sb := squirrel.StatementBuilder.Select("*").From("testing")
		sb = applyFilterToSubCountQueryBuilder(nil, exampleTableName, sb)
		expected := "SELECT * FROM testing"
		actual, _, err := sb.ToSql()

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}