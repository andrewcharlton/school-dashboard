package database

// A Config holds the details of the schools
type Config struct {

	// Name of the school
	School string

	// URN/LEA numbers (as found on performance tables)
	URN string
	LEA string

	// Default filter options
	Options map[string]string
}

// DefaultFilter produces a Filter object with the
// default values specified in the database.
func (cfg Config) DefaultFilter() Filter {

	return Filter{Date: cfg.Options["Date"],
		Resultset: cfg.Options["Resultset"],
		NatYear:   cfg.Options["NatYear"],
		Year:      cfg.Options["Year"],
	}
}

// config retrieves the config values from the database
func (db *Database) loadConfig() error {

	rows, err := db.conn.Query("SELECT key, value FROM config")
	if err != nil {
		return err
	}
	defer rows.Close()

	cfg := Config{Options: map[string]string{}}
	for rows.Next() {
		var key, value string
		err := rows.Scan(&key, &value)
		if err != nil {
			return err
		}
		switch key {
		case "School":
			cfg.School = value
		case "URN":
			cfg.URN = value
		case "LEA":
			cfg.LEA = value
		default:
			cfg.Options[key] = value
		}
	}

	db.Config = cfg
	return nil
}
