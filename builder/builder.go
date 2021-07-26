package builder

import (
	"fmt"
	"github.com/ulovecode/gdcache/schemas"
	"strings"
)

func GetEntryByIdSQL(entry schemas.IEntry, entryParams []schemas.EntryKey) string {
	var idWhereString []string
	for _, entryParam := range entryParams {
		idWhereString = append(idWhereString, fmt.Sprintf(" %s = %s ", entryParam.Name, entryParam.Param))
	}
	return fmt.Sprintf(`SELECT * FROM %s  WHERE %s;`, entry.GetTableName(), strings.Join(idWhereString, "AND"))
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
	return fmt.Sprintf(`SELECT * FROM %s  WHERE %s;`, entry.GetTableName(), strings.Join(idWhereString, " OR "))
}
