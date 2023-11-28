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

	var artists Artists

	err := GetJson(url, &artists)
	if err != nil {
		fmt.Printf("Error getting Artist info: %s\n", err.Error())
	} else {
		fmt.Printf("Here is your artist name: %s\n", artists.Name)
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

func main() {

	var artists Artists

	url := "https://groupietrackers.herokuapp.com/api/artists"
	err := GetJson(url, &artists)
	if err != nil {
		fmt.Printf("Error getting Artist info: %s\n", err.Error())
		return
	}

	fmt.Printf("Here is your artist name: %s\n", artists.Name)
}
