package handlers

import (
	"database/sql"
	"hospital-management/config"
	"hospital-management/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if user.Role != "receptionist" && user.Role != "doctor" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	var existingUser models.User
	err := config.DB.QueryRow("SELECT username FROM users WHERE username = $1", user.Username).Scan(&existingUser.Username)
	if err != sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	_, err = config.DB.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3)", user.Username, string(hashedPassword), user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var dbUser models.User
	err := config.DB.QueryRow("SELECT id, username, password, role FROM users WHERE username = $1", user.Username).
		Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password, &dbUser.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not registered. Please sign up first."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if user.Role != dbUser.Role {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Role mismatch"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": dbUser.Username,
		"role":     dbUser.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(config.JWTSecret))
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
