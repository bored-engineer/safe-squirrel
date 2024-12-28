package squirrel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateBuilderToSql(t *testing.T) {
	b := Update("").
		Prefix("WITH prefix AS ?", 0).
		Table("a").
		Set("b", Expr("? + 1", 1)).
		SetMap(Eq{"c": 2}).
		Set("c1", Case(Expr("status")).When(Expr("1"), Expr("2")).When(Expr("2"), Expr("1"))).
		Set("c2", Case().When(Expr("a = 2"), Expr("?", "foo")).When(Expr("a = 3"), Expr("?", "bar"))).
		Set("c3", Select("a").From("b")).
		Where(Expr("d = ?", 3)).
		OrderBy("e").
		Limit(4).
		Offset(5).
		Suffix("RETURNING ?", 6)

	sql, args, err := b.ToSql()
	assert.NoError(t, err)

	expectedSql :=
		"WITH prefix AS ? " +
			"UPDATE a SET b = ? + 1, c = ?, " +
			"c1 = CASE status WHEN 1 THEN 2 WHEN 2 THEN 1 END, " +
			"c2 = CASE WHEN a = 2 THEN ? WHEN a = 3 THEN ? END, " +
			"c3 = (SELECT a FROM b) " +
			"WHERE d = ? " +
			"ORDER BY e LIMIT 4 OFFSET 5 " +
			"RETURNING ?"
	assert.Equal(t, expectedSql, sql)

	expectedArgs := []interface{}{0, 1, 2, "foo", "bar", 3, 6}
	assert.Equal(t, expectedArgs, args)
}

func TestUpdateBuilderToSqlErr(t *testing.T) {
	_, _, err := Update("").Set("x", 1).ToSql()
	assert.Error(t, err)

	_, _, err = Update("x").ToSql()
	assert.Error(t, err)
}

func TestUpdateBuilderMustSql(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("TestUpdateBuilderMustSql should have panicked!")
		}
	}()
	Update("").MustSql()
}

func TestUpdateBuilderPlaceholders(t *testing.T) {
	b := Update("test").SetMap(Eq{"x": 1, "y": 2})

	sql, _, _ := b.PlaceholderFormat(Question).ToSql()
	assert.Equal(t, "UPDATE test SET x = ?, y = ?", sql)

	sql, _, _ = b.PlaceholderFormat(Dollar).ToSql()
	assert.Equal(t, "UPDATE test SET x = $1, y = $2", sql)
}

func TestUpdateBuilderRunners(t *testing.T) {
	db := &DBStub{}
	b := Update("test").Set("x", 1).RunWith(db)

	expectedSql := "UPDATE test SET x = ?"

	b.Exec()
	assert.Equal(t, expectedSql, db.LastExecSql)
}

func TestUpdateBuilderNoRunner(t *testing.T) {
	b := Update("test").Set("x", 1)

	_, err := b.Exec()
	assert.Equal(t, RunnerNotSet, err)
}

func TestUpdateBuilderFrom(t *testing.T) {
	sql, _, err := Update("employees").Set("sales_count", 100).From("accounts").Where(Expr("accounts.name = ?", "ACME")).ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE employees SET sales_count = ? FROM accounts WHERE accounts.name = ?", sql)
}

func TestUpdateBuilderFromSelect(t *testing.T) {
	sql, _, err := Update("employees").
		Set("sales_count", 100).
		FromSelect(Select("id").
			From("accounts").
			Where(Expr("accounts.name = ?", "ACME")), "subquery").
		Where(Expr("employees.account_id = subquery.id")).ToSql()
	assert.NoError(t, err)

	expectedSql :=
		"UPDATE employees " +
			"SET sales_count = ? " +
			"FROM (SELECT id FROM accounts WHERE accounts.name = ?) AS subquery " +
			"WHERE employees.account_id = subquery.id"
	assert.Equal(t, expectedSql, sql)
}
