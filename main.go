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

	if r.URL.Path != "/" {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("Error getting artists: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = templates.ExecuteTemplate(w, "artists.html", artists)
	if err != nil {
		log.Printf("Error executing template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
func main() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", ArtistHandler)
	log.Println("Serving on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Println("Error starting server:", err)
	}
}