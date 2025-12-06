package database

import (
	"fmt"
	"strings"
	"time"
)

func DebugSQL(query string, args ...interface{}) string {
	for i, arg := range args {
		placeholder := fmt.Sprintf("$%d", i+1)
		query = strings.Replace(query, placeholder, fmt.Sprintf("'%v'", arg), 1)
	}
	return query
}

func DebugSQL2(query string, args ...interface{}) string {
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			query = strings.Replace(query, "?", fmt.Sprintf("'%v'", v), 1)
		case int, int64, float64:
			query = strings.Replace(query, "?", fmt.Sprintf("%v", v), 1)
		case time.Time:
			query = strings.Replace(query, "?", fmt.Sprintf("'%v'", v.Format("2006-01-02 15:04:05")), 1)
		default:
			query = strings.Replace(query, "?", fmt.Sprintf("'%v'", v), 1)
		}
	}
	return query
}
