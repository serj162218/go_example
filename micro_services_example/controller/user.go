package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/serj162218/go_example/micro_services_example/helper"
	"github.com/serj162218/go_example/micro_services_example/initializer"
	"github.com/serj162218/go_example/micro_services_example/model"
	"golang.org/x/crypto/bcrypt"
)

func UserRegister(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(user.Email) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	if len(user.Password) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}

	if len(user.ID) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if _, err := initializer.DB.Exec("INSERT INTO users (id, email, password) VALUES (?, ?, ?)", user.ID, user.Email, hashedPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}

func UserLogin(c *gin.Context) {
	// Handle user login
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(user.ID) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	if len(user.Password) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}

	var dbUser model.User
	if err := initializer.DB.QueryRow("SELECT id, email, password FROM users WHERE id = ?", user.ID).Scan(&dbUser.ID, &dbUser.Email, &dbUser.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	// Generate JWT token
	token, err := helper.GenerateToken(dbUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return token
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func UserLogout(c *gin.Context) {
	// Handle user logout
	tokenString := c.GetHeader("Authorization")
	// Add this token to black list which can be store in database or redis
	err := model.AddTokenToBlacklist(tokenString)
	if err != nil {
		// Handle error
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Successful logout
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successfully",
	})
}

func UserProtectedEndpoint(c *gin.Context) {
	// Handle protected endpoint
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"message": "protected endpoint"})
}

func UserAuthJWTMiddleware(next gin.HandlerFunc) gin.HandlerFunc {
	// Handle JWT authentication middleware
	return func(c *gin.Context) {
		// Read token from request Header
		tokenString := c.GetHeader("Authorization")
		token, err := helper.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Check if the token is valid
		if !token.Valid {
			// JWT token is invalid, return error
			c.JSON(http.StatusUnauthorized, gin.H{"err": "Invalid JWT token"})
			return
		}
		//Check if the token is in the black list
		if model.IsTokenInBlackList(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "Invalid JWT token"})
			return
		}

		// JWT token vaildation passed, call next function
		next(c)
	}
}
