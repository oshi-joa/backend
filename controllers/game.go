package controllers

import (
	"backend/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func NewGameController(r *gin.Engine, mongo *mongo.Client, db *gorm.DB) {
	game := r.Group("game")
	{
		game.POST("/post", func(c *gin.Context) {
			services.GamePOST(c, mongo)
		})
		game.POST("/answer1/:id", func(c *gin.Context) {
			services.Answer1(c, mongo, db)
		})
		game.POST("/answer2/:id", func(c *gin.Context) {
			services.Answer2(c, mongo, db)
		})
		game.GET("/get", func(c *gin.Context) {
			services.GameALLGET(c, mongo)
		})
		game.GET("/get/:id", func(c *gin.Context) {
			services.GameGETByID(c, mongo)
		})
		game.POST("/get/couple/:id", func(c *gin.Context) {
			services.CheckAnswer(c, mongo, db)
		})
	}
}
