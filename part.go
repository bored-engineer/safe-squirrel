package squirrel

import (
	"io"
)

func nestedToSql(s Sqlizer) (string, []interface{}, error) {
	if raw, ok := s.(rawSqlizer); ok {
		return raw.toSqlRaw()
	} else {
		return s.ToSql()
	}
}

func appendToSql(parts []Sqlizer, w io.Writer, sep string, args []interface{}) ([]interface{}, error) {
	for i, p := range parts {
		partSql, partArgs, err := nestedToSql(p)
		if err != nil {
			return nil, err
		} else if len(partSql) == 0 {
			continue
		}

		if i > 0 {
			_, err := io.WriteString(w, sep)
			if err != nil {
				return nil, err
			}
		}

		_, err = io.WriteString(w, partSql)
		if err != nil {
			return nil, err
		}
		args = append(args, partArgs...)
	}
	return args, nil
}
