package models

import (
	"regexp"
	"unicode"
)

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

type RegisterData struct {
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
	Terms           bool
}

type RegisterErrors struct {
	UsernameError string
	EmailError    string
	PasswordError string
	ConfirmError  string
	TermsError    string
}

// ValidateRegister valide les données d'inscription
func ValidateRegister(data RegisterData) *RegisterErrors {
	errors := &RegisterErrors{}

	// Validation du nom d'utilisateur
	if len(data.Username) < 3 {
		errors.UsernameError = "Le nom d'utilisateur doit contenir au moins 3 caractères"
	} else if len(data.Username) > 20 {
		errors.UsernameError = "Le nom d'utilisateur ne doit pas dépasser 20 caractères"
	} else if !isValidUsername(data.Username) {
		errors.UsernameError = "Le nom d'utilisateur ne peut contenir que des lettres, chiffres, _ et -"
	}

	// Validation de l'email
	if !isValidEmail(data.Email) {
		errors.EmailError = "Veuillez entrer une adresse email valide"
	}

	// Validation du mot de passe
	if len(data.Password) < 8 {
		errors.PasswordError = "Le mot de passe doit contenir au moins 8 caractères"
	} else if !isStrongPassword(data.Password) {
		errors.PasswordError = "Le mot de passe doit contenir des majuscules, minuscules et des chiffres"
	}

	// Vérification de la confirmation du mot de passe
	if data.Password != data.ConfirmPassword {
		errors.ConfirmError = "Les mots de passe ne correspondent pas"
	}

	// Vérification des conditions d'utilisation
	if !data.Terms {
		errors.TermsError = "Vous devez accepter les conditions d'utilisation"
	}

	return errors
}

func isValidUsername(username string) bool {
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", username)
	return matched
}

func isValidEmail(email string) bool {
	matched, _ := regexp.MatchString("^[^\\s@]+@[^\\s@]+\\.[^\\s@]+$", email)
	return matched
}

func isStrongPassword(password string) bool {
	hasUppercase := false
	hasLowercase := false
	hasNumbers := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUppercase = true
		}
		if unicode.IsLower(char) {
			hasLowercase = true
		}
		if unicode.IsDigit(char) {
			hasNumbers = true
		}
	}

	return hasUppercase && hasLowercase && hasNumbers
}

// HasErrors vérifie s'il y a des erreurs de validation
func (e *RegisterErrors) HasErrors() bool {
	return e.UsernameError != "" || e.EmailError != "" || e.PasswordError != "" || e.ConfirmError != "" || e.TermsError != ""
}
