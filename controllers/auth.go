package controllers

import (
	"backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthController(r *gin.Engine, db *gorm.DB) {
	r.POST("/signup", func(c *gin.Context) {
		services.SignUp(c, db)
	})
	r.GET("/login", func(c *gin.Context) {
		services.Login(c, db)
	})
}
