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

	// Add client server
	clientMux := http.NewServeMux()

	// Serve static files
	static := http.FileServer(http.Dir("./static"))
	clientMux.Handle("/static/", http.StripPrefix("/static/", static))

	// Serve image files
	images := http.FileServer(http.Dir("./images"))
	clientMux.Handle("/images/", http.StripPrefix("/images/", images))

	// Handlers
	clientMux.HandleFunc("/", handlers.Index(env))
	clientMux.HandleFunc("/basics/", handlers.EnglishAndMaths(env))
	clientMux.HandleFunc("/headlines/", handlers.Headlines(env))
	clientMux.HandleFunc("/progress8/", handlers.Progress8(env))
	clientMux.HandleFunc("/effort/", handlers.Effort(env))
	clientMux.HandleFunc("/subjects/", handlers.SubjectOverview(env))
	clientMux.HandleFunc("/progressgrid/", handlers.ProgressGrid(env))
	clientMux.HandleFunc("/classlist/", handlers.ClassList(env))
	clientMux.HandleFunc("/students/", handlers.Student(env))
	clientMux.HandleFunc("/studentsearch/", handlers.SearchRedirect(env))
	clientMux.HandleFunc("/search/", handlers.Search(env))

	// Start client server
	go func() {
		http.ListenAndServe(":8080", clientMux)
	}()
}
