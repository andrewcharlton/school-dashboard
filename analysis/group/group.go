// Package group provides the tools for analysing a group
// of students.
package group

import (
	"github.com/andrewcharlton/school-dashboard/analysis/student"
)

// A Group collects a group of students together
type Group struct {
	Students []student.Student
}
