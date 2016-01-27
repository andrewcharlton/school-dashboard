package database

// A Config holds the details of the schools
type Config struct {

	// Name of the school
	School string

	// URN/LEA numbers (as found on performance tables)
	URN string
	LEA string

	// Default filter options
	date      string
	resultset string
	natYear   string
	year      string
}

// DefaultFilter produces a Filter object with the
// default values specified in the database.
func (cfg Config) DefaultFilter() Filter {

	return Filter{Date: cfg.date,
		Resultset: cfg.resultset,
		NatYear:   cfg.natYear,
		Year:      cfg.year,
	}
}

// config retrieves the config values from the database
func (db Database) config() (Config, error) {

	rows, err := db.conn.Query("SELECT key, value FROM config")
	if err != nil {
		return Config{}, err
	}
	defer rows.Close()

	cfg := Config{}
	for rows.Next() {
		var key, value string
		err := rows.Scan(&key, &value)
		if err != nil {
			return Config{}, err
		}
		switch key {
		case "School":
			cfg.School = value
		case "URN":
			cfg.URN = value
		case "LEA":
			cfg.LEA = value
		case "Date":
			cfg.Date = value
		case "Resultset":
			cfg.Resultset = value
		case "NatYear":
			cfg.NatYear = value
		case "Year":
			cfg.Year = value
		}
	}
	return cfg, nil
}
