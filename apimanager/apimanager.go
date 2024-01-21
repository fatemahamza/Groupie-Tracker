package apimanager
import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)
type Artists struct {
	ID             int                 `json:"id"`
	Image          string              `json:"image"`
	Name           string              `json:"name"`
	Members        []string            `json:"members"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
func GetJson(url string, target any) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}
func GetArtists() ([]Artists, error) {
	url := "https://groupietrackers.herokuapp.com/api/artists"
	var artists []Artists
	err := GetJson(url, &artists)
	if err != nil {
		return nil, fmt.Errorf("error getting Artist info: %s", err.Error())
	}
	for i, artist := range artists {
		relationUrl := "https://groupietrackers.herokuapp.com/api/relation/" + strconv.Itoa(artist.ID)
		var tempVal map[string]any
		err := GetJson(relationUrl, &tempVal)
		if err != nil {
			return nil, fmt.Errorf("error getting Artist Relation info: %s", err.Error())
		}
		convertedDates := make(map[string][]string)
		datesLocations := tempVal["datesLocations"].(map[string]interface{})
		for key, value := range datesLocations {
			// value looks like this "[12-10-2004 17-6-2004 24-2-2004]"
			//remove brackets
			strippedValue := strings.TrimSuffix(fmt.Sprint(value)[1:], "]")
			// split by space
			convertedDates[key] = strings.Split(strippedValue, " ")
		}
		artists[i].DatesLocations = convertedDates
	}
	return artists, nil
}
