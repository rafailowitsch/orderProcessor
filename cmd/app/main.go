package main

import (
	_ "github.com/swaggo/http-swagger"
	"orderProcessor/internal/app"
)

// @title Order Processor API
// @version 1.0
// @description This is a sample server for managing orders.
// @host localhost:8089
// @BasePath /
// @schemes http
func main() {
	app.Run()
}
