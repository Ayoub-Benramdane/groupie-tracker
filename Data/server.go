package Tracker

import (
	"html/template"
	"net/http"
	"strconv"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		temp, err := template.ParseFiles("src/Error.html")
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(404)
		temp.Execute(w, nil)
		return
	} else if r.Method != http.MethodGet {
		http.Error(w, "Status Method Not Allowed 405", http.StatusMethodNotAllowed)
		return
	}
	temp, err := template.ParseFiles("src/home.html")
	if err != nil {
		http.Error(w, "Status Internal Server Error 500", http.StatusInternalServerError)
		return
	}
	temp.Execute(w, Artists())
}

func Groupie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Status Method Not Allowed 405", http.StatusMethodNotAllowed)
		return
	}
	temp, err := template.ParseFiles("src/artist.html")
	if err != nil {
		http.Error(w, "Status Internal Server Error 500", http.StatusInternalServerError)
		return
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || (Artist(id).Name) == "" {
		temp, err := template.ParseFiles("src/Error.html")
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(404)
		temp.Execute(w, nil)
		return
	}
	temp.Execute(w, Artist(id))
}
