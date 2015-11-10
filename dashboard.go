package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/andrewcharlton/school-dashboard/env"
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
	env, err := env.Connect(filename)
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
	http.HandleFunc("/index/", handlers.Index(env))
<<<<<<< HEAD
	http.HandleFunc("/basics/", handlers.EnglishAndMaths(env))
=======
	http.HandleFunc("/headlines/", handlers.Headlines(env))
>>>>>>> befe32a338d926d73344879d41d2d204a44b5013
	http.HandleFunc("/effort/", handlers.Effort(env))
	http.HandleFunc("/classlist/", handlers.ClassList(env))
	http.HandleFunc("/students/", handlers.Student(env))
	http.HandleFunc("/studentsearch/", handlers.SearchRedirect(env))
	http.HandleFunc("/search/", handlers.Search(env))

	// Start server
	http.ListenAndServe(":8080", nil)
}
