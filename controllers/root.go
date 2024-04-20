package controllers

import (
	"backend/repositories"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func NewController(port string) {
	r := gin.Default()

	r.Use(cors.New(
		cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
			MaxAge:       24 * time.Hour,
		}))

	db := repositories.MySQLInit()
	AuthController(r, db)
	NewBoardController(r, db)

	err := r.Run(port)
	if err != nil {
		panic(err)
	}
}
