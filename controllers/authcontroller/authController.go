package authController

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jeypc/go-jwt-mux/helper"
	"golang.org/x/crypto/bcrypt"

	"github.com/jeypc/go-jwt-mux/models"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": "Invalid request body"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// Validate email
	if !helper.IsValidEmail(userInput.Email) {
		response := map[string]string{"message": "Invalid email format"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Validate password length
	if len(userInput.Password) < 6 {
		response := map[string]string{"message": "Password must be at least 6 characters long"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Hash password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		response := map[string]string{"message": "Error hashing password"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}
	userInput.Password = string(hashPassword)

	// Create user
	if err := models.DB.Create(&userInput).Error; err != nil {
		response := map[string]string{"message": "Error creating user"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "User registered successfully"}
	helper.ResponseJSON(w, http.StatusCreated, response)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": "Invalid request body"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// Retrieve user by email
	var user models.User
	if err := models.DB.Where("email = ?", userInput.Email).First(&user).Error; err != nil {
		response := map[string]string{"message": "Invalid email or password"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		response := map[string]string{"message": "Invalid email or password"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	// Generate JWT token
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours
	claims := &config.JWTClaim{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "go-jwt-mux",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWTKey)
	if err != nil {
		response := map[string]string{"message": "Error generating token"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Set JWT token as cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
	})

	response := map[string]string{"message": "Login successful"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	response := map[string]string{"message": "Logout successful"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	// Implement update user logic here based on the user ID
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": "Invalid request body"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// Update user in the database based on userID
	if err := models.DB.Model(&models.User{}).Where("id = ?", userID).Updates(&userInput).Error; err != nil {
		response := map[string]string{"message": "Error updating user"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "User updated successfully"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	// Implement delete user logic here based on the user ID
	if err := models.DB.Where("id = ?", userID).Delete(&models.User{}).Error; err != nil {
		response := map[string]string{"message": "Error deleting user"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "User deleted successfully"}
	helper.ResponseJSON(w, http.StatusOK, response)
}
