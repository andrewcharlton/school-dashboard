package subject

// A Result brings together a subject and the grade achieved
// in that subject.
type Result struct {
	*Subject
	Grade
	Effort  int
	Class   string
	Teacher string
}
