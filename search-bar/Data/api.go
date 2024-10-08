package Tracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type rela struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type loca struct {
	ID   int      `json:"id"`
	Loca []string `json:"locations"`
}

type artists struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Date         dates
	Location     loca
	Relation     rela
}

func fetchAll(id int, wg *sync.WaitGroup, result *artists, fetch string) {
	defer wg.Done()
	url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/%s/%d", fetch, id)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response code")
		return
	}
	switch fetch {
	case "artists":
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			fmt.Println("Error decoding JSON:", err)
		}
	case "locations":
		var location loca
		if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
			fmt.Println("Error decoding JSON:", err)
		} else {
			result.Location = location
		}
	case "relation":
		var relation rela
		if err := json.NewDecoder(resp.Body).Decode(&relation); err != nil {
			fmt.Println("Error decoding JSON:", err)
		} else {
			result.Relation = relation
		}
	case "dates":
		var date dates
		if err := json.NewDecoder(resp.Body).Decode(&date); err != nil {
			fmt.Println("Error decoding JSON:", err)
		} else {
			result.Date = date
		}
	}
}

func Artist(id int) artists {
	var wg sync.WaitGroup
	result := artists{}
	dataTypes := []string{"artists", "locations", "dates", "relation"}
	for _, fetch := range dataTypes {
		wg.Add(1)
		go fetchAll(id, &wg, &result, fetch)
	}
	wg.Wait()
	return result
}

func List(src string) []artists {
	result := make([]artists, 0)
	date, _ := strconv.Atoi(src)
	for _, c := range Artists() {
		if c.CreationDate == date || strings.Contains(strings.ToLower(c.Name), src) || strings.Contains(c.FirstAlbum, src) {
			result = append(result, c)
		} else {
			isValid := true
			for _, v := range c.Members {
				if strings.Contains(strings.ToLower(v), src) {
					result = append(result, c)
					isValid = false
					break
				}
			}
			if isValid {
				for _, v := range c.Location.Loca {
					if strings.Contains(strings.ToLower(v), src) {
						result = append(result, c)
						break
					}
				}
			}
		}
	}
	return result
}

func Artists() []artists {
	var wg sync.WaitGroup
	dataTypes := []string{"artists", "locations"}
	result := make([]artists, 52)
	for i := 0; i < len(result); i++ {
		for _, c := range dataTypes {
			wg.Add(1)
			go fetchAll(i+1, &wg, &result[i], c)
		}
	}
	wg.Wait()
	return result
}
