package main

import (
	"flag"
	"fmt"
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

	// Static and Image servers
	static := http.FileServer(http.Dir("./static"))
	images := http.FileServer(http.Dir("./images"))

	// Client Server
	clientMux := http.NewServeMux()
	clientMux.Handle("/static/", http.StripPrefix("/static/", static))
	clientMux.Handle("/images/", http.StripPrefix("/images/", images))
	clientMux.HandleFunc("/", handlers.Index(env))
	clientMux.HandleFunc("/attendance/", handlers.AttendanceExplorer(env))
	clientMux.HandleFunc("/attendancegroups/", handlers.AttendanceGroups(env))
	//	clientMux.HandleFunc("/basics/", handlers.EnglishAndMaths(env))
	//clientMux.HandleFunc("/headlines/", handlers.Headlines(env))
	clientMux.HandleFunc("/progress8/", handlers.Progress8(env))
	clientMux.HandleFunc("/progress8groups/", handlers.Progress8Groups(env))
	//clientMux.HandleFunc("/effort/", handlers.Effort(env))
	//clientMux.HandleFunc("/export/headlines/", handlers.ExportHeadlines(env))
	//clientMux.HandleFunc("/subjects/", handlers.SubjectOverview(env))
	clientMux.HandleFunc("/progressgrid/", handlers.ProgressGrid(env))
	clientMux.HandleFunc("/subjectgroups/", handlers.SubjectGroups(env))
	//clientMux.HandleFunc("/classlist/", handlers.ClassList(env))
	//clientMux.HandleFunc("/students/", handlers.Student(env))
	//clientMux.HandleFunc("/studentsearch/", handlers.SearchRedirect(env))
	//clientMux.HandleFunc("/search/", handlers.Search(env))
	go func() {
		http.ListenAndServe(":8080", clientMux)
	}()

	adminMux := http.NewServeMux()
	adminMux.Handle("/static/", http.StripPrefix("/static/", static))
	adminMux.HandleFunc("/admin/", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "Hello") })
	http.ListenAndServe(":8081", adminMux)
}
