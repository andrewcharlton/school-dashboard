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
