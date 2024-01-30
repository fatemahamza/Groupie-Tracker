package apimanager
import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)
const (
	artistsEndpoint  = "https://groupietrackers.herokuapp.com/api/artists"
	relationEndpoint = "https://groupietrackers.herokuapp.com/api/relation/"
)
type Artists struct {
	ID             int                 `json:"id"`
	Image          string              `json:"image"`
	Name           string              `json:"name"`
	Members        []string            `json:"members"`
	DatesLocations map[string][]string `json:"datesLocations"`
	CreationDate   int                 `json:"creationDate"`
	FirstAlbum     string              `json:"firstAlbum"`
}
func GetJson(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d for URL: %s", resp.StatusCode, url)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}
func GetArtistDetails(artistName string) (Artists, error) {
	url := artistsEndpoint
	var artists []Artists
	err := GetJson(url, &artists)
	if err != nil {
		return Artists{}, fmt.Errorf("error getting Artist info: %s", err.Error())
	}
	for _, artist := range artists {
		if artist.Name == artistName {
			// Retrieve additional details for the found artist
			relationUrl := relationEndpoint + strconv.Itoa(artist.ID)
			var tempVal map[string]interface{}
			err := GetJson(relationUrl, &tempVal)
			if err != nil {
				return Artists{}, fmt.Errorf("error getting Artist Relation info: %s", err.Error())
			}
			datesLocations, ok := tempVal["datesLocations"].(map[string]interface{})
			if !ok {
				return Artists{}, fmt.Errorf("datesLocations type assertion failed")
			}
			convertedDates := make(map[string][]string)
			for key, value := range datesLocations {
				strippedValue := strings.TrimSuffix(fmt.Sprint(value)[1:], "]")
				convertedDates[key] = strings.Split(strippedValue, " ")
			}
			artist.DatesLocations = convertedDates
			return artist, nil
		}
	}
	return Artists{}, fmt.Errorf("artist not found: %s", artistName)
}
func GetArtists() ([]Artists, error) {
	url := artistsEndpoint
	var artists []Artists
	err := GetJson(url, &artists)
	if err != nil {
		return nil, fmt.Errorf("error getting Artist info: %s", err.Error())
	}
	for i, artist := range artists {
		relationUrl := relationEndpoint + strconv.Itoa(artist.ID)
		var tempVal map[string]interface{}
		err := GetJson(relationUrl, &tempVal)
		if err != nil {
			return nil, fmt.Errorf("error getting Artist Relation info: %s", err.Error())
		}
		datesLocations, ok := tempVal["datesLocations"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("datesLocations type assertion failed")
		}
		convertedDates := make(map[string][]string)
		for key, value := range datesLocations {
			// value looks like this "[12-10-2004 17-6-2004 24-2-2004]"
			// remove brackets
			strippedValue := strings.TrimSuffix(fmt.Sprint(value)[1:], "]")
			// split by space
			convertedDates[key] = strings.Split(strippedValue, " ")
		}
		artists[i].DatesLocations = convertedDates
	}
	return artists, nil
}