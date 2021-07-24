package schemas

type IEntry interface {
	GetTableName() string
}

type IEntries []IEntry

func GetPKsByEntries(entries []IEntry) (PK, error) {
	pks := make(PK, 0)
	for i := 0; i < len(entries); i++ {
		_, pk, err := GetEntryKey(entries[i])
		if err != nil {
			return pks, err
		}
		pks = append(pks, pk)
	}
	return pks, nil
}
