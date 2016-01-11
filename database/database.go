// Package database implements a database access layer
// for the school-dashboard application.
//
// It wraps an SQL database connection, prepares
// statements to query the database and then
// provides methods to return the results of any
// queries.
package database

import (
	"github.com/andrewcharlton/school-dashboard/analysis"
	"github.com/andrewcharlton/school-dashboard/level"
	"github.com/andrewcharlton/school-dashboard/national"
)

// A Database provides a wrapper to the database
// connection and provides methods to query it.
type Database interface {

	// Close the database connection
	Close() error

	// Config returns a populated config object from the
	// database.
	Config() (Config, error)

	// Dates returns a sorted list of all effective dates
	// in the database that are marked to be listed.
	Dates() ([]Lookup, error)

	// Ethnicities returns a list of all the distinct
	// ethnicities present in the database, and the frequency
	// that each appears with.  Only students who are present
	// in listed dates are counted.  List should be sorted
	// with largest groups first.
	Ethnicities() ([]Ethnicity, error)

	// Resultsets returns a sorted list of all resultsets
	// in the database.
	Resultsets() ([]Lookup, error)

	// Student returns a student object with details relevant
	// to the given filter.
	Student(f StudentFilter) (analysis.Student, error)

	// GroupByFilter returns a list of students who satisfy the
	// criteria specified in the filter
	GroupByFilter(f Filter) (analysis.Group, error)

	// GroupByClass returns a group of students who are present
	// in the subject/class at the date specified in the filter.
	GroupByClass(subj_id, class string, f Filter) (analysis.Group, error)

	// GroupByFilteredClass returns a group of students who meet
	// the filter criteria and are also present in the subject/
	// class combination.  If class="", the group will include
	// all students who study that subject.
	GroupByFilteredClass(subj_id, class string, f Filter) (analysis.Group, error)

	// Search returns a list of students from a search query.
	Search(name, date string) ([]StudentLookup, error)

	// Subjects returns a list of all subjects available.
	Subjects() map[int]*analysis.Subject

	// Levels returns a list of all levels available.
	Level(name string) *level.Level

	// Classes returns a list of classes that exist for a subject,
	// at a particular date (date_id should be provided).
	Classes(subj_id, date string) ([]string, error)

	// NationalYears returns a sorted list of all years
	// for which national data is available
	NationalYears() ([]Lookup, error)

	// National returns a set of national data for the given year.
	National(yearID string) (national.National, error)

	// HistoricalResults returns a map of student results for
	// a particular subject, keyed by resultset id
	HistoricalResults(upn, subjID string) (map[string]string, error)
}

// A Lookup holds an ID/Name pair for an item in the database.
type Lookup struct {
	ID   string
	Name string
}

// An Ethnicity holds the name of an Ethnic group, as well as the
// number of students in that group.
type Ethnicity struct {
	Name  string
	Count int
}

// A StudentLookup holds the basic details of a student,
// returned from a search
type StudentLookup struct {
	UPN  string
	Name string
	Year string
	Form string
}

// A Filter holds all the fields available to select a group
// of students on.
type Filter struct {
	// Effective date to use for membership of the group.
	// Holds the id number, but held as string for easier
	// parsing from query strings.
	Date string
	// Which set of assessment results to use.
	// Holds id number as above.
	Resultset string
	// Which year to use for national data comparisons.
	NatYear string
	// Which yeargroup to include in the group
	Year string
	// Pupil premium: "", "1", or "0" for Any/True/False
	PP string
	// EAL students: "", "1", or "0" for Any/True/False
	EAL string
	// Gender: "", "1", "0" for Any/Male/Female
	Gender string
	// SEN types to include - Empty for any
	SEN []string
	// Which Ethnic groups to include, empty for any
	Ethnicities []string
	// Which KS2 bands to include, empty for any
	KS2Bands []string
}

// A StudentFilter holds all the fields necessary to uniquely
// identify a student and the grades etc.
type StudentFilter struct {
	UPN       string
	Date      string
	Resultset string
}

// A Config holds all of the
type Config struct {

	// Name of the school
	School string

	// URN/LEA numbers (as found on performance tables)
	URN string
	LEA string

	// Default filter options
	Date      string
	Resultset string
	NatYear   string
	Year      string
}

func (cfg Config) DefaultFilter() Filter {

	return Filter{Date: cfg.Date,
		Resultset: cfg.Resultset,
		NatYear:   cfg.NatYear,
		Year:      cfg.Year,
	}
}
