package main
import (
	"groupie-tracker/apimanager"
	"log"
	"net/http"
	"text/template"
)
var templates *template.Template
func init() {
	// Initialize templates during package initialization
	templates = template.Must(template.ParseGlob("templates/*.html"))
}
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := apimanager.GetArtists()
	if err != nil {
		log.Printf("Error getting artists: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Set the Content-Type header before writing to the response writer
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.URL.Path == "/" {
		err = templates.ExecuteTemplate(w, "artists.html", artists)
		if err != nil {
			log.Printf("Error executing template: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		log.Printf("Error finding page: %s", r.URL.Path)
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
}
func DetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the artist name from the query parameters
	artistName := r.URL.Query().Get("artist")
	if artistName == "" {
		http.Error(w, "Artist name not provided", http.StatusBadRequest)
		return
	}
	// Retrieve details for the specified artist
	artistDetails, err := apimanager.GetArtistDetails(artistName)
	if err != nil {
		log.Printf("Error getting artist details for %s: %s", artistName, err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Set the Content-Type header before writing to the response writer
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Pass the artist details to the template
	if r.URL.Path == "/details.html" {
		err = templates.ExecuteTemplate(w, "details.html", artistDetails)
		if err != nil {
			log.Printf("Error executing template: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		log.Printf("Error finding page: %s", r.URL.Path)
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
}
func main() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", ArtistHandler)
	mux.HandleFunc("/details.html", DetailsHandler)
	log.Println("Serving on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Println("Error starting server:", err)
	}
}