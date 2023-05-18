package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/serj162218/go_example/micro_services_example/controller"
	"github.com/serj162218/go_example/micro_services_example/initializer"
)

func main() {
	//initialize
	initializer.Initialize()

	//router
	router := gin.Default()
	router.POST("/register", controller.UserRegister)
	router.POST("/login", controller.UserLogin)
	router.POST("/logout", controller.UserLogout)
	router.GET("/protected", controller.UserAuthJWTMiddleware(controller.UserProtectedEndpoint))

	log.Fatal(http.ListenAndServe(":8080", router))
}
