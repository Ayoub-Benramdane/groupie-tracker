package Tracker

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		temp, err := template.ParseFiles("src/Error.html")
		if err != nil {
			http.Error(w, "Status Internal Server Error 500", http.StatusInternalServerError)
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

func Filter(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/filter" {
		temp, err := template.ParseFiles("src/Error.html")
		if err != nil {
			http.Error(w, "Status Internal Server Error 500", http.StatusInternalServerError)
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
		http.Error(w, "Status Internal Server Error 500b", http.StatusInternalServerError)
		return
	}
	loc := strings.ToLower(r.FormValue("search"))
	table := []string{}
	table = append(table, r.FormValue("creMin"))
	table = append(table, r.FormValue("creMax"))
	table = append(table, r.FormValue("fAlbMin"))
	table = append(table, r.FormValue("fAlbMax"))
	membe := r.Form["member[]"]
	filters := List(loc, table, membe)
	if len(filters) != 0 {
		temp.Execute(w, filters)
	} else {
		temp.Execute(w, nil)
	}
}

func Geo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Status Method Not Allowed 405", http.StatusMethodNotAllowed)
		return
	}
	temp, err := template.ParseFiles("src/geo.html")
	if err != nil {
		http.Error(w, "Status Internal Server Error 500", http.StatusInternalServerError)
		return
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		temp, err := template.ParseFiles("src/Error.html")
		if err != nil {
			http.Error(w, "Status Internal Server Error 500", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(404)
		temp.Execute(w, nil)
		return
	}
	fmt.Println(id)
	temp.Execute(w, nil)
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
			http.Error(w, "Status Internal Server Error 500", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(404)
		temp.Execute(w, nil)
		return
	}
	temp.Execute(w, Artist(id))
}

func Stylise(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Status Method Not Allowed 405", http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Path[len("/css/"):]
	if path == "style.css" {
		http.ServeFile(w, r, "css/style.css")
	} else {
		temp, err := template.ParseFiles("src/Error.html")
		if err != nil {
			http.Error(w, "Status Forbidden 403", http.StatusForbidden)
			return
		}
		w.WriteHeader(403)
		temp.Execute(w, nil)
		return
	}
}
