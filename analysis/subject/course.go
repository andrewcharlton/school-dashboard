package subject

// A Course brings together a subject and the grade achieved
// in that subject.
type Course struct {
	Subject
	Grade
	TM      TransitionMatrix
	Effort  int
	Class   string
	Teacher string
}
