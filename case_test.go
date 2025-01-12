package squirrel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaseWithVal(t *testing.T) {
	caseStmt := Case(Expr("number")).
		When(Expr("1"), Expr("one")).
		When(Expr("2"), Expr("two")).
		Else(Expr("?", "big number"))

	qb := Select().
		Column(caseStmt).
		From("table")
	sql, args, err := qb.ToSql()

	assert.NoError(t, err)

	expectedSql := "SELECT CASE number " +
		"WHEN 1 THEN one " +
		"WHEN 2 THEN two " +
		"ELSE ? " +
		"END " +
		"FROM table"
	assert.Equal(t, expectedSql, sql)

	expectedArgs := []interface{}{"big number"}
	assert.Equal(t, expectedArgs, args)
}

func TestCaseWithComplexVal(t *testing.T) {
	caseStmt := Case(Expr("? > ?", 10, 5)).
		When(Expr("true"), Expr("'T'"))

	qb := Select().
		Column(Alias(caseStmt, "complexCase")).
		From("table")
	sql, args, err := qb.ToSql()

	assert.NoError(t, err)

	expectedSql := "SELECT (CASE ? > ? " +
		"WHEN true THEN 'T' " +
		"END) AS complexCase " +
		"FROM table"
	assert.Equal(t, expectedSql, sql)

	expectedArgs := []interface{}{10, 5}
	assert.Equal(t, expectedArgs, args)
}

func TestCaseWithNoVal(t *testing.T) {
	caseStmt := Case().
		When(Eq{"x": 0}, Expr("x is zero")).
		When(Expr("x > ?", 1), Expr("CONCAT('x is greater than ', ?)", 2))

	qb := Select().Column(caseStmt).From("table")
	sql, args, err := qb.ToSql()

	assert.NoError(t, err)

	expectedSql := "SELECT CASE " +
		"WHEN x = ? THEN x is zero " +
		"WHEN x > ? THEN CONCAT('x is greater than ', ?) " +
		"END " +
		"FROM table"

	assert.Equal(t, expectedSql, sql)

	expectedArgs := []interface{}{0, 1, 2}
	assert.Equal(t, expectedArgs, args)
}

func TestCaseWithExpr(t *testing.T) {
	caseStmt := Case(Expr("x = ?", true)).
		When(Expr("true"), Expr("?", "it's true!")).
		Else(Expr("42"))

	qb := Select().Column(caseStmt).From("table")
	sql, args, err := qb.ToSql()

	assert.NoError(t, err)

	expectedSql := "SELECT CASE x = ? " +
		"WHEN true THEN ? " +
		"ELSE 42 " +
		"END " +
		"FROM table"

	assert.Equal(t, expectedSql, sql)

	expectedArgs := []interface{}{true, "it's true!"}
	assert.Equal(t, expectedArgs, args)
}

func TestMultipleCase(t *testing.T) {
	caseStmtNoval := Case(Expr("x = ?", true)).
		When(Expr("true"), Expr("?", "it's true!")).
		Else(Expr("42"))
	caseStmtExpr := Case().
		When(Eq{"x": 0}, Expr("'x is zero'")).
		When(Expr("x > ?", 1), Expr("CONCAT('x is greater than ', ?)", 2))

	qb := Select().
		Column(Alias(caseStmtNoval, "case_noval")).
		Column(Alias(caseStmtExpr, "case_expr")).
		From("table")

	sql, args, err := qb.ToSql()

	assert.NoError(t, err)

	expectedSql := "SELECT " +
		"(CASE x = ? WHEN true THEN ? ELSE 42 END) AS case_noval, " +
		"(CASE WHEN x = ? THEN 'x is zero' WHEN x > ? THEN CONCAT('x is greater than ', ?) END) AS case_expr " +
		"FROM table"

	assert.Equal(t, expectedSql, sql)

	expectedArgs := []interface{}{
		true, "it's true!",
		0, 1, 2,
	}
	assert.Equal(t, expectedArgs, args)
}

func TestCaseWithNoWhenClause(t *testing.T) {
	caseStmt := Case(Expr("something")).
		Else(Expr("42"))

	qb := Select().Column(caseStmt).From("table")

	_, _, err := qb.ToSql()

	assert.Error(t, err)

	assert.Equal(t, "case expression must contain at lease one WHEN clause", err.Error())
}

func TestCaseBuilderMustSql(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("TestCaseBuilderMustSql should have panicked!")
		}
	}()
	Case(Expr("")).MustSql()
}
