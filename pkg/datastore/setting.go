package datastore

type Setting struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func FetchSettings() (settings []Setting, err error) {
	rows, err := db.Query(`SELECT * FROM settings ORDER BY name ASC`)
	if err != nil {
		return
	}

	// Fetch rows
	for rows.Next() {
		setting := new(Setting)

		// get RawBytes from data
		err = rows.Scan(&setting.Name, &setting.Value)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		settings = append(settings, *setting)

	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return
}

func GetSettingByName(name string) (setting Setting, err error) {
	stmtOut, err := db.Prepare(`SELECT name, value FROM settings WHERE name = ?`)
	if err != nil {
		return
	}
	defer stmtOut.Close()

	err = stmtOut.QueryRow(name).Scan(&setting.Name, &setting.Value)
	if err != nil {
		return
	}

	return
}

func SetSettingByName(name, value string) (err error) {
	_, err = GetSettingByName(name)
	if err != nil {
		//	Insert
		stmtIns, err := db.Prepare(`INSERT INTO settings(name, value) VALUES( ?, ? )`)
		if err != nil {
			return err
		}
		defer stmtIns.Close()

		_, err = stmtIns.Exec(name, value)
		if err != nil {
			return err
		}

		return err
	}

	// Update
	stmtUpd, err := db.Prepare(`UPDATE settings SET value = ? WHERE name = ?`)
	if err != nil {
		return
	}
	defer stmtUpd.Close()

	_, err = stmtUpd.Exec(value, name)
	if err != nil {
		return
	}

	return
}
