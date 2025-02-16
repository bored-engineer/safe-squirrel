package squirrel

import "github.com/lann/builder"

// StatementBuilderType is the type of StatementBuilder.
type StatementBuilderType builder.Builder

// Select returns a SelectBuilder for this StatementBuilderType.
func (b StatementBuilderType) Select(columns ...safeString) SelectBuilder {
	return SelectBuilder(b).Columns(columns...)
}

// Insert returns a InsertBuilder for this StatementBuilderType.
func (b StatementBuilderType) Insert(into safeString) InsertBuilder {
	return InsertBuilder(b).Into(into)
}

// Replace returns a InsertBuilder for this StatementBuilderType with the
// statement keyword set to "REPLACE".
func (b StatementBuilderType) Replace(into safeString) InsertBuilder {
	return InsertBuilder(b).statementKeyword("REPLACE").Into(into)
}

// Update returns a UpdateBuilder for this StatementBuilderType.
func (b StatementBuilderType) Update(table safeString) UpdateBuilder {
	return UpdateBuilder(b).Table(table)
}

// Delete returns a DeleteBuilder for this StatementBuilderType.
func (b StatementBuilderType) Delete(from safeString) DeleteBuilder {
	return DeleteBuilder(b).From(from)
}

// PlaceholderFormat sets the PlaceholderFormat field for any child builders.
func (b StatementBuilderType) PlaceholderFormat(f PlaceholderFormat) StatementBuilderType {
	return builder.Set(b, "PlaceholderFormat", f).(StatementBuilderType)
}

// RunWith sets the RunWith field for any child builders.
func (b StatementBuilderType) RunWith(runner BaseRunner) StatementBuilderType {
	return setRunWith(b, runner).(StatementBuilderType)
}

// Where adds WHERE expressions to the query.
//
// See SelectBuilder.Where for more information.
func (b StatementBuilderType) Where(expr Sqlizer) StatementBuilderType {
	return builder.Append(b, "WhereParts", expr).(StatementBuilderType)
}

// StatementBuilder is a parent builder for other builders, e.g. SelectBuilder.
var StatementBuilder = StatementBuilderType(builder.EmptyBuilder).PlaceholderFormat(Question)

// Select returns a new SelectBuilder, optionally setting some result columns.
//
// See SelectBuilder.Columns.
func Select(columns ...safeString) SelectBuilder {
	return StatementBuilder.Select(columns...)
}

// Insert returns a new InsertBuilder with the given table name.
//
// See InsertBuilder.Into.
func Insert(into safeString) InsertBuilder {
	return StatementBuilder.Insert(into)
}

// Replace returns a new InsertBuilder with the statement keyword set to
// "REPLACE" and with the given table name.
//
// See InsertBuilder.Into.
func Replace(into safeString) InsertBuilder {
	return StatementBuilder.Replace(into)
}

// Update returns a new UpdateBuilder with the given table name.
//
// See UpdateBuilder.Table.
func Update(table safeString) UpdateBuilder {
	return StatementBuilder.Update(table)
}

// Delete returns a new DeleteBuilder with the given table name.
//
// See DeleteBuilder.Table.
func Delete(from safeString) DeleteBuilder {
	return StatementBuilder.Delete(from)
}

// Case returns a new CaseBuilder
// "what" represents optional case value
func Case(what ...Sqlizer) CaseBuilder {
	b := CaseBuilder(builder.EmptyBuilder)

	switch len(what) {
	case 0:
		return b
	default:
		return b.what(what[0])
	}
}
