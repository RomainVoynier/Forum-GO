package handlers

import (
	"net/http"
	"text/template"

	"forum-go/models"
)

// RegisterPageData pour passer les données au template
type RegisterPageData struct {
	Username      string
	Email         string
	UsernameError string
	EmailError    string
	PasswordError string
	ConfirmError  string
	TermsError    string
	Error         string
	Success       bool
}

// LoginPageData pour passer les données au template de connexion
type LoginPageData struct {
	Username      string
	UsernameError string
	PasswordError string
	Error         string
	Success       bool
}

// RegisterPage affiche la page d'inscription
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("html/account.html")
	if err != nil {
		http.Error(w, "Erreur de chargement de la page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

// Register traite la soumission du formulaire d'inscription
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// Parser le formulaire
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erreur lors du traitement du formulaire", http.StatusBadRequest)
		return
	}

	// Récupérer les données
	registerData := models.RegisterData{
		Username:        r.FormValue("username"),
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
		Terms:           r.FormValue("terms") == "on",
	}

	// Valider les données
	validationErrors := models.ValidateRegister(registerData)

	// S'il y a des erreurs de validation
	if validationErrors.HasErrors() {
		tmpl, _ := template.ParseFiles("html/account.html")
		data := RegisterPageData{
			Username:      registerData.Username,
			Email:         registerData.Email,
			UsernameError: validationErrors.UsernameError,
			EmailError:    validationErrors.EmailError,
			PasswordError: validationErrors.PasswordError,
			ConfirmError:  validationErrors.ConfirmError,
			TermsError:    validationErrors.TermsError,
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, data)
		return
	}

	// Vérifier si l'utilisateur existe déjà
	exists, _ := models.UserExists(registerData.Username, registerData.Email)
	if exists {
		tmpl, _ := template.ParseFiles("html/account.html")
		data := RegisterPageData{
			Username: registerData.Username,
			Email:    registerData.Email,
			Error:    "Cet utilisateur ou cet email existe déjà",
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, data)
		return
	}

	// Créer l'utilisateur dans la base de données
	user, err := models.CreateUser(registerData.Username, registerData.Email, registerData.Password)
	if err != nil {
		tmpl, _ := template.ParseFiles("html/account.html")
		data := RegisterPageData{
			Username: registerData.Username,
			Email:    registerData.Email,
			Error:    "Erreur lors de la création du compte",
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, data)
		return
	}

	// Afficher le succès et connecter l'utilisateur
	setSessionCookie(w, user.Username)

	tmpl, _ := template.ParseFiles("html/account.html")
	data := RegisterPageData{
		Username: user.Username,
		Success:  true,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

// LoginPage affiche la page de connexion
func LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("html/login.html")
	if err != nil {
		http.Error(w, "Erreur de chargement de la page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

// Login traite la soumission du formulaire de connexion
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erreur lors du traitement du formulaire", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	data := LoginPageData{
		Username: username,
	}

	// Valider les entrées
	if len(username) == 0 {
		data.UsernameError = "Veuillez entrer un nom d'utilisateur ou un email"
	}
	if len(password) == 0 {
		data.PasswordError = "Veuillez entrer votre mot de passe"
	}

	if data.UsernameError != "" || data.PasswordError != "" {
		tmpl, _ := template.ParseFiles("html/login.html")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, data)
		return
	}

	// Chercher l'utilisateur par nom d'utilisateur
	user, err := models.GetUserByUsername(username)
	if err != nil {
		// Essayer par email
		user, err = models.GetUserByEmail(username)
		if err != nil {
			data.Error = "Nom d'utilisateur ou email incorrect"
			tmpl, _ := template.ParseFiles("html/login.html")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			tmpl.Execute(w, data)
			return
		}
	}

	// Vérifier le mot de passe
	// ⚠️ À améliorer avec bcrypt en production
	if user.Password != password {
		data.Error = "Mot de passe incorrect"
		tmpl, _ := template.ParseFiles("html/login.html")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, data)
		return
	}

	// Succès ! Connexion réussie
    setSessionCookie(w, user.Username)
    data.Success = true
    data.Username = user.Username
    tmpl, _ := template.ParseFiles("html/login.html")
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    tmpl.Execute(w, data)
}

func setSessionCookie(w http.ResponseWriter, username string) {
    cookie := &http.Cookie{
        Name:     "forum_username",
        Value:    username,
        Path:     "/",
        HttpOnly: true,
        MaxAge:   86400,
    }
    http.SetCookie(w, cookie)
}
