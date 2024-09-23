package main

import (
	"fmt"
	"net/http"

	Tracker "Tracker/Data"
)

func main() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
	http.HandleFunc("/", Tracker.Index)
	http.HandleFunc("/artists/{id}", Tracker.Groupie)
	fmt.Println("Server is Running... http://localhost:8404")
	err := http.ListenAndServe(":8404", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
