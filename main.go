package main

import (
	"log"
	"net/http"

	"forum-go/handlers"
	"forum-go/models"
)

func main() {
	// Initialiser la base de données
	err := models.InitDB()
	if err != nil {
		log.Fatal("❌ Erreur lors de l'initialisation de la BD:", err)
	}
	defer models.DB.Close()

	// Route pour la page d'accueil
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "html/index.html")
			return
		}
		http.NotFound(w, r)
	})

	// Route pour l'inscription (GET et POST)
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.RegisterPage(w, r)
		} else if r.Method == http.MethodPost {
			handlers.Register(w, r)
		}
	})

	// Route pour la connexion (GET et POST)
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.LoginPage(w, r)
		} else if r.Method == http.MethodPost {
			handlers.Login(w, r)
		}
	})

	// Servir les fichiers statiques (CSS, images, etc.)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	addr := ":8080"
	log.Println("✅ Serveur lancé sur http://localhost:8080")
	log.Println("📝 Page d'accueil: http://localhost:8080")
	log.Println("📝 Page d'inscription: http://localhost:8080/register")
	log.Println("📝 Page de connexion: http://localhost:8080/login")
	log.Fatal(http.ListenAndServe(addr, nil))
}
