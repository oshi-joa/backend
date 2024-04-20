package controllers

import (
	"backend/repositories"
	"github.com/gin-gonic/gin"
)

func NewController(port string) {
	r := gin.Default()

	db := repositories.MySQLInit()
	AuthController(r, db)

	err := r.Run(port)
	if err != nil {
		panic(err)
	}
}
