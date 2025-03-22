package Tracker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type Coord struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
	Type      string `json:"addresstype"`
	FullDisp  string `json:"display_name"`
	Name      string
}

type Country struct {
	Name struct {
		Long string `json:"common"`
	} `json:"name"`
}

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
	Coordinates  [][]Coord
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

func GetCountry(short string) (string, error) {
	short = strings.ToUpper(short) // Ensure uppercase
	url := fmt.Sprintf("https://restcountries.com/v3.1/alpha/%s", short)

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var countries []Country
	err = json.Unmarshal(body, &countries)
	if err != nil || len(countries) == 0 {
		return "", fmt.Errorf("invalid country code")
	}

	return countries[0].Name.Long, nil
}

func Address(addresses []string) ([][]Coord, error) {
	var results [][]Coord

	for _, address := range addresses {
		encodedAddress := url.QueryEscape(address)
		addr := strings.Split(encodedAddress, "-")
		apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?format=json&q=%s", addr[0])

		resp, err := http.Get(apiURL)
		if err != nil {
			fmt.Printf("Failed to get address %s: %v\n", address, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error fetching %s: %s\n", address, resp.Status)
			continue
		}

		var res []Coord
		err = json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			fmt.Printf("Error decoding response for %s: %v\n", address, err)
			continue
		}
		for i := 0; i < len(res); i++ {
			a := res[i].FullDisp

			country, err := GetCountry(addr[1])
			if len(addr[1]) == 3 {
				if err != nil {
					fmt.Println("Error: Invalid country adress")
					continue
				}
			} else {
				country = addr[1]
			}
			if !(strings.Contains(a, country) || ValidAddr(res[i].Type)) {
				res = append(res[:i], res[i+1:]...)
				i--
			}
		}
		res[0].Name = addr[0] + ", " + addr[1]
		if len(res) == 0 {
			fmt.Printf("No results found for address: %s\n", address)
			continue
		}

		results = append(results, res)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no valid results found")
	}
	return results, nil
}

func ValidAddr(addr string) bool {
	return addr == "state" || addr == "city" || addr == "province" || addr == "railway"
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

func List(location string, table, membe []string) []artists {
	result := make([]artists, 0)
	result1 := make([]artists, 0)
	indice := make([]int, 0)
	for _, c := range table {
		nb, _ := strconv.Atoi(c)
		indice = append(indice, nb)
	}
	a := make(map[string]int)
	b := make(map[string]int)
	for _, c := range Artists() {
		fAlb, _ := strconv.Atoi(c.FirstAlbum[len(c.FirstAlbum)-4:])
		if c.CreationDate >= indice[0] && c.CreationDate <= indice[1] && fAlb >= indice[2] && fAlb <= indice[3] && a[c.Name] == 0 {
			result = append(result, c)
			a[c.Name]++
		}
	}
	for _, v := range membe {
		for _, d := range result {
			h, _ := strconv.Atoi(v)
			if h == len(d.Members) && b[d.Name] == 0 {
				result1 = append(result1, d)
				b[d.Name]++
			}
		}
	}
	if len(result1) != 0 {
		result = result1
	}
	for i := 0; i < len(result); i++ {
		for j, v := range result[i].Location.Loca {
			if !strings.Contains(strings.ToLower(v), location) && j == len(result[i].Location.Loca)-1 {
				result = append(result[:i], result[i+1:]...)
				i--
				break
			}
			if strings.Contains(strings.ToLower(v), location) {
				break
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
