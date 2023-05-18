package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/serj162218/go_example/micro_services_example/controller"
	"github.com/serj162218/go_example/micro_services_example/helper"
	"github.com/serj162218/go_example/micro_services_example/initializer"
	"github.com/serj162218/go_example/micro_services_example/model"
)

func TestUserRegister(t *testing.T) {
	//set a router
	router := gin.Default()
	router.POST("/register", controller.UserRegister)

	//create a request with a body with JSON
	user := model.User{
		ID:       "test",
		Email:    "test@example.com",
		Password: "123456",
	}
	requestBody, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	// Check the status code is what we expect.
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expectedResponse := `{"message":"User registered"}`
	if recorder.Body.String() != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v want %v",
			recorder.Body.String(), expectedResponse)
	}

	//Remember that this function is for register. So we need to delete user after testing.
	initializer.DB.Exec("DELETE FROM users WHERE id = 'test'")
}
func TestUserLogin(t *testing.T) {
	router := gin.Default()
	router.POST("/login", controller.UserLogin)

	//Which is already in database.
	user := model.User{
		ID:       "testLogin",
		Email:    "testLogin@example.com",
		Password: "testLogin",
	}
	requestBody, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	// Check the status code is what we expect.
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestUserLogout(t *testing.T) {
	router := gin.Default()
	router.POST("/logout", controller.UserLogout)

	token := "testToken"
	req, _ := http.NewRequest("POST", "/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	// Check the status code is what we expect.
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expectedResponse := `{"message":"Logout successfully"}`
	if recorder.Body.String() != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v want %v",
			recorder.Body.String(), expectedResponse)
	}

	//This function will add the token into black list.
	isExist, err := initializer.RDB.SIsMember(context.Background(), "black_list", token).Result()
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
	if isExist != true {
		t.Errorf("expected %v, but got %v", true, isExist)
	}

	//Remember to delete the token after testing.
	initializer.RDB.SRem(context.Background(), "black_list", token)
}

func TestUserProtectedEndpoint(t *testing.T) {
	// bypass the middleware.
	router := gin.Default()
	router.GET("/protected", controller.UserProtectedEndpoint)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expectedResponse := `{"message":"protected endpoint"}`
	if recorder.Body.String() != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v want %v",
			recorder.Body.String(), expectedResponse)
	}
}

func TestUserAuthJWTMiddleware(t *testing.T) {
	//Custom the handler next middleware.
	router := gin.Default()
	router.GET("/userAuthJWTMiddleware", controller.UserAuthJWTMiddleware(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok!"})
	}))

	//the JwtKey should be as same as in the middleware.
	token, err := helper.GenerateToken(model.User{}, initializer.JwtKey)
	if err != nil {
		t.Errorf("generate token got error: %v", err)
	}

	req, _ := http.NewRequest("GET", "/userAuthJWTMiddleware", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expectedResponse := `{"message":"ok!"}`
	if recorder.Body.String() != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v want %v",
			recorder.Body.String(), expectedResponse)
	}
}
