package main

import (
	"net/http"

	Tracker "Tracker/Data"
)

func main() {
	http.HandleFunc("/", Tracker.Index)
	http.HandleFunc("/{{id}}", Tracker.Groupie)
}
