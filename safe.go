package squirrel

// stringConst ensures that a given string is a 'const' value at compile-time.
type stringConst string

func (s stringConst) ToSql() (sqlStr string, args []interface{}, err error) {
	return string(s), nil, nil
}

// DangerouslyCastDynamicStringToSQL converts a dynamic string to a stringConst for use in the methods/types of this package.
// This should be used with _extreme_ caution, as it will lead to SQL injection if the string has not been properly sanitized.
//
// Deprecated: This function is dangerous and should not be used unless you are _very_ sure you know what you're doing.
func DangerouslyCastDynamicStringToSQL(val string) stringConst {
	return stringConst(val)
}

// SetMap can be passed to the SetMap function in various builders
type SetMap map[stringConst]interface{}
