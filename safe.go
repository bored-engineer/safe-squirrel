package squirrel

// stringConst ensures that a given string is a 'const' value at compile-time.
type stringConst string

func (s stringConst) ToSql() (sqlStr string, args []interface{}, err error) {
	return string(s), nil, nil
}
