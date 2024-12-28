# Safe Squirrel - fluent SQL generator for Go [![Go Reference](https://pkg.go.dev/badge/github.com/bored-engineer/safe-squirrel.svg)](https://pkg.go.dev/github.com/bored-engineer/safe-squirrel)
"safe" squirrel is a _fork_ of [github.com/Masterminds/squirrel](https://github.com/Masterminds/squirrel) that enforces more secure/safe usage via the Golang type system. 

It can be adopted as a drop-in replacement ([caveat](#caveats)) for squirrel just by changing the replacing the import with `github.com/bored-engineer/safe-squirrel`.

# Why?
The [squirrel](https://github.com/Masterminds/squirrel) package already encourages the use of [parameterized queries (aka placeholders)](https://cheatsheetseries.owasp.org/cheatsheets/Query_Parameterization_Cheat_Sheet.html) to avoid the risk of [SQL injection](https://owasp.org/www-community/attacks/SQL_Injection) in, ex:
```go
username := "bored-engineer" // untrusted input

// "SELECT * FROM users WHERE github = ?"
sq.Select("*").From("users").Where(sq.Eq{"github": username}).ToSql()
```
However, not all methods/parameters is safe/protected against SQL injection, ex:
```go
provider := "is_superadmin=true OR github" // untrusted input
username := "uh oh" // untrusted input

// "SELECT * FROM users WHERE is_superadmin=true OR github = ?"
sq.Select("*").From("users").Where(sq.Eq{provider: username}).ToSql()
```
While this is a contrived example, [SQL injection](https://owasp.org/www-community/attacks/SQL_Injection) vulnerabilities have been found in real-world applications/services that use [squirrel](https://github.com/Masterminds/squirrel) due to incorrect usage of these APIs.

# How?
By taking advantage of the Golang type system/compiler, it is possible to create a function that will _only_ accept a `const` string at compile-time, ex:
```go
type stringConst string

func refuseDynamicStrings(foo stringConst) {
    println(foo)
}
```

When this function is invoked Golang will automatically cast `const` strings to the private (otherwise inaccessible) `stringConst` type:
```go
pkg.refuseDynamicStrings("this is an implicit const string")
const foo = "this is an explicit const string"
pkg.refuseDynamicStrings(foo)
```

However, if we try to pass a dynamic string (such as one generated from user input), it will fail to build/compile:
```go
var bar = fmt.Sprintf("this is a %s string", "dynamic")
pkg.refuseDynamicStrings(bar) // cannot use bar (variable of type string) as stringConst value in argument to refuseDynamicStrings
```

This package/fork takes advantage of this "feature" to enforce that all parameters passed to [squirrel](https://github.com/Masterminds/squirrel) are `const` strings at compile-time. APIs that are already secure due to their use of [parameterized queries](https://cheatsheetseries.owasp.org/cheatsheets/Query_Parameterization_Cheat_Sheet.html) like `sq.Expr("foo = ?", untrustedStringVar)` or `sq.Eq{"column": untrustedStringVar}` continue to work as-is, accepting dynamic values on the relevant types.

# Caveats
Most of the [squirrel](https://github.com/Masterminds/squirrel) APIs were directly converted from `string` to `stringConst` requiring no refactoring for an application to adopt this fork. However, the [Where](https://pkg.go.dev/github.com/bored-engineer/safe-squirrel#SelectBuilder.Where), [Having](https://pkg.go.dev/github.com/bored-engineer/safe-squirrel#SelectBuilder.Having) and [Case](https://pkg.go.dev/github.com/bored-engineer/safe-squirrel#Case) APIs now only accept a [Sqlizer](https://pkg.go.dev/github.com/bored-engineer/safe-squirrel#Sqlizer) type. If you were previously using the less common `(sql string, args ...interface{})` invocation, some simple refactoring to insert [sq.Expr(...)](https://pkg.go.dev/github.com/bored-engineer/safe-squirrel#Expr) will be required, ex:
```go
// before
builder.Where("foo = ?", untrustedStringVar)

// after
builder.Where(sq.Expr("foo = ?", untrustedStringVar))
```
Notably this is _not_ required if the application is already using the comparison types from [squirrel](https://github.com/Masterminds/squirrel), ex:
```go
builder.Where(sq.Eq{"foo": untrustedStringVar})
```

If an unsafe/insecure [Sqlizer](https://pkg.go.dev/github.com/bored-engineer/safe-squirrel#Sqlizer) that was defined outside of the `safe-squirrel` package is used, this _could_ still result in a [SQL injection](https://owasp.org/www-community/attacks/SQL_Injection) vulnerability as the value returned by `ToSql` is used as-is.
