package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/andrewcharlton/school-dashboard/handlers"
)

func main() {

	// Get database filename
	args := flag.Args()
	filename := "school.db"
	if len(args) >= 1 {
		filename = args[1]
	}

	// Connect to database
	env, err := database.Connect(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve image files
	is := http.FileServer(http.Dir("./images"))
	http.Handle("/images/", http.StripPrefix("/images/", is))

	// Handlers
	http.HandleFunc("/", handlers.Index(env))
	http.HandleFunc("/basics/", handlers.EnglishAndMaths(env))
	http.HandleFunc("/headlines/", handlers.Headlines(env))
	http.HandleFunc("/effort/", handlers.Effort(env))
	http.HandleFunc("/ks4subject/", handlers.KS4Subject(env))
	http.HandleFunc("/classlist/", handlers.ClassList(env))
	http.HandleFunc("/students/", handlers.Student(env))
	http.HandleFunc("/studentsearch/", handlers.SearchRedirect(env))
	http.HandleFunc("/search/", handlers.Search(env))

	// Start server
	http.ListenAndServe(":8080", nil)
}
