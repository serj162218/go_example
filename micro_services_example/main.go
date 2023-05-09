package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BlackList struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

var jwtKey = []byte("shh!it's_secret_key")
var db *sql.DB
var rdb *redis.Client

func main() {
	var err error
	db, err = sql.Open("mysql", "micro_services_example:micro_services_example@tcp(127.0.0.1:3306)/micro_services_example")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	generateRandomSecretKey()
	router := gin.Default()
	router.POST("/register", register)
	router.POST("/login", login)
	router.POST("/logout", logout)
	router.GET("/protected", authJWTMiddleware(protectedEndpoint))

	log.Fatal(http.ListenAndServe(":8080", router))
}
func generateRandomSecretKey() {
	secretKey := make([]byte, 32)
	_, err := rand.Read(secretKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(secretKey))
	jwtKey = secretKey
}

func register(c *gin.Context) {
	var user User
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
	if _, err := db.Exec("INSERT INTO users (id, email, password) VALUES (?, ?, ?)", user.ID, user.Email, hashedPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}

func login(c *gin.Context) {
	// Handle user login
	var user User
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

	var dbUser User
	if err := db.QueryRow("SELECT id, email, password FROM users WHERE id = ?", user.ID).Scan(&dbUser.ID, &dbUser.Email, &dbUser.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	// Generate JWT token
	token, err := generateToken(dbUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return token
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func logout(c *gin.Context) {
	// Handle user logout
	tokenString := c.GetHeader("Authorization")
	// Add this token to black list which can be store in database or redis
	err := addTokenToBlacklist(tokenString)
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

func addTokenToBlacklist(token string) error {
	//store this token to the database
	err := rdb.SAdd(context.TODO(), "black_list", token).Err()
	if err != nil {
		return err
	}
	return nil
}

func protectedEndpoint(c *gin.Context) {
	// Handle protected endpoint
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"message": "protected endpoint"})
}

func authJWTMiddleware(next gin.HandlerFunc) gin.HandlerFunc {
	// Handle JWT authentication middleware
	return func(c *gin.Context) {
		// Read token from request Header
		tokenString := c.GetHeader("Authorization")
		token, err := verifyToken(tokenString)
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
		if isTokenInBlackList(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "Invalid JWT token"})
			return
		}

		// JWT token vaildation passed, call next function
		next(c)
	}
}

func isTokenInBlackList(token string) bool {
	//check if the token is in redis
	isExist, err := rdb.SIsMember(context.TODO(), "black_list", token).Result()
	if err != nil {
		log.Fatal(err.Error())
		return true
	}
	return isExist
}

func generateToken(user User) (string, error) {
	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	// Verify JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check that the JWT token is valid
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// return JWT token
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
