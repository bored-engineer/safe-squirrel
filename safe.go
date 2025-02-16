package squirrel

import "strings"

// safeString ensures that a given string is a 'const' value at compile-time.
type safeString string

func (s safeString) ToSql() (sqlStr string, args []interface{}, err error) {
	return string(s), nil, nil
}

// SafeString allows callers to explicitly declare a 'safeString'
func SafeString(val safeString) safeString {
	return val
}

// SafeString allows callers to explicitly declare multiple 'safeString's
func SafeStrings(vals ...safeString) []safeString {
	return vals
}

// JoinSafeStrings joins multiple 'safeString's into a single 'safeString'
func JoinSafeStrings(sep safeString, vals ...safeString) safeString {
	var sb strings.Builder
	for idx, val := range vals {
		if idx > 0 {
			sb.WriteString(string(sep))
		}
		sb.WriteString(string(val))
	}
	return safeString(sb.String())
}

// DangerouslyCastDynamicStringToSafeString converts a dynamic string to a safeString for use in the methods/types of this package.
// This should be used with _extreme_ caution, as it will lead to SQL injection if the string has not been properly sanitized.
//
// Deprecated: This function is dangerous and should not be used unless you are _very_ sure you know what you're doing.
func DangerouslyCastDynamicStringToSafeString(val string) safeString {
	return safeString(val)
}

// SetMap can be passed to the SetMap function in various builders
type SetMap map[safeString]interface{}
