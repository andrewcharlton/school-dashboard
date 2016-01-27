package database

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
