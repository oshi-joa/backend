package controllers

import (
	"backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewBoardController(r *gin.Engine, db *gorm.DB) {
	board := r.Group("board")
	{
		board.POST("post", func(c *gin.Context) {
			services.BoardPost(c, db)
		})
		board.GET("get", func(c *gin.Context) {
			services.BoardAllGet(c, db)
		})
		board.GET("get/:id", func(c *gin.Context) {
			services.BoardGetByID(c, db)
		})
		board.DELETE("delete/:id", func(c *gin.Context) {
			services.BoardDeleteByID(c, db)
		})
	}
}
