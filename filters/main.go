package main

import (
	"fmt"
	"net/http"

	Tracker "Tracker/Data"
)

func main() {
	http.HandleFunc("/css/", Tracker.Stylise)
	http.HandleFunc("/", Tracker.Index)
	http.HandleFunc("/filter", Tracker.Filter)
	http.HandleFunc("/artists/{id}", Tracker.Groupie)
	fmt.Println("Server is Running... http://localhost:8404")
	err := http.ListenAndServe(":8404", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
