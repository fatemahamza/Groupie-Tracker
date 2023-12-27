package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var client *http.Client

type Artists struct {
	ID      int      `json:"id"`
	Image   string   `json:"image"`
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

func GetArtists() {
	url := "https://groupietrackers.herokuapp.com/api/artists"
	var artists []Artists
	err := GetJson(url, &artists)
	if err != nil {
		fmt.Printf("Error getting Artist info: %s\n", err.Error())
		return
	}

	fmt.Println("Here are the artist members:")
	for _, artist := range artists {
		fmt.Printf("Artist: %s\n", artist.Name)
		fmt.Println("Members:")
		for _, member := range artist.Members {
			fmt.Println(member)
		}
		fmt.Println()
	}
}

func GetJson(url string, target interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	GetArtists()
	//json.NewEncoder(w).Encode(artists)
}

func main() {
	client = &http.Client{} // Initialize the HTTP client

	http.HandleFunc("/artists", ArtistHandler)

	fmt.Println("Server is running on: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("error starting server")
	}
	GetArtists()
}