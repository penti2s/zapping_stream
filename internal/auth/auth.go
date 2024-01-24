package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"zapping_stream/internal/db"
	"zapping_stream/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func Register(r *http.Request) (string, error) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return "", err
	}

	salt := os.Getenv("SALT")
	var user model.User
	if err := db.GetDB().Where("email = ?", req.Email).First(&user).Error; err == nil {
		return "", errors.New("un usuario con ese email ya existe") // Añade "" como primer valor
	}

	passwordWithSalt := fmt.Sprintf("%s%s", req.Password, salt)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordWithSalt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	newUser := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := db.GetDB().Create(&newUser).Error; err != nil {
		return "", err
	}

	token, err := GenerateToken(newUser)
	if err != nil {
		return "", err
	}

	return token, nil
}

func Login(r *http.Request) (string, error) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return "", err
	}

	var user model.User
	if err := db.GetDB().Where("email = ?", req.Email).First(&user).Error; err != nil {
		return "", errors.New("usuario no encontrado")
	}

	salt := os.Getenv("SALT")
	passwordWithSalt := fmt.Sprintf("%s%s", req.Password, salt)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordWithSalt)); err != nil {
		return "", errors.New("contraseña incorrecta")
	}

	token, err := GenerateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")

	claims, err := ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "Token inválido", http.StatusUnauthorized)
		return
	}

	var user model.User
	if err := db.GetDB().First(&user, claims.UserID).Error; err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	userResponse := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	jsonResponse, err := json.Marshal(userResponse)
	if err != nil {
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Token-Expired", "true")

	response := map[string]string{"message": "Has cerrado sesión correctamente."}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
