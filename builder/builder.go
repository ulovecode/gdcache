package builder

import (
	"fmt"
	"github.com/ulovecode/gdcache/schemas"
	"reflect"
	"strings"
)

func GetEntryByIdSQL(entry schemas.IEntry, entryParams []schemas.EntryKey) string {
	var idWhereString []string
	for _, entryParam := range entryParams {
		idWhereString = append(idWhereString, fmt.Sprintf(" %s = %s ", entryParam.Name, entryParam.Param))
	}
	return fmt.Sprintf(`SELECT * FROM %s  WHERE %s;`, entry.TableName(), strings.Join(idWhereString, "AND"))
}

func GetEntriesByIdSQL(entry schemas.IEntry, entryKeys []schemas.EntryKeys) string {
	var idWhereString []string
	for _, entryParams := range entryKeys {
		var idString []string
		for _, entryParam := range entryParams {
			idString = append(idString, fmt.Sprintf(" %s = %s ", entryParam.Name, entryParam.Param))
		}
		idWhereString = append(idWhereString, fmt.Sprintf(`( %s )`, strings.Join(idString, "AND")))
	}
	return fmt.Sprintf(`SELECT * FROM %s  WHERE %s;`, entry.TableName(), strings.Join(idWhereString, " OR "))
}

func GenerateSql(sql string, args ...interface{}) string {
	if len(args) == 0 {
		return sql
	}
	params := make([]interface{}, 0)
	for _, arg := range args {
		definitelyNotEqualString := fmt.Sprintf(" ( 1 != 1 ) ")
		if arg == nil {
			params = append(params, definitelyNotEqualString)
		}

		argValue := reflect.ValueOf(arg)
		if argValue.Kind() == reflect.Ptr {
			if argValue.IsNil() {
				params = append(params, definitelyNotEqualString)
				continue
			}
			arg = argValue.Elem()
		}

		if argValue.Kind() == reflect.Slice {
			if argValue.Len() == 0 {
				params = append(params, definitelyNotEqualString)
				continue
			}
			argSQL := fmt.Sprint(arg)
			arg = "(" + strings.Replace(argSQL[1:len(argSQL)-1], " ", ",", -1) + ")"
		}
		params = append(params, fmt.Sprint(arg))
	}
	sql = strings.Replace(sql, "?", "%s", -1)
	return fmt.Sprintf(sql, params...)
}

func GenerateCountSql(sql string, args ...interface{}) string {
	sql = GenerateSql(sql, args...)
	var indexOf int
	findSql := strings.ToUpper(sql)
	indexOf = strings.LastIndex(findSql, "LIMIT")
	if indexOf != -1 {
		s := sql[indexOf:]
		if !strings.ContainsAny(s, ")") {
			sql = sql[:indexOf]
		}
	}
	return fmt.Sprintf(`SELECT COUNT(*) FROM (%s) t`, sql)
}
