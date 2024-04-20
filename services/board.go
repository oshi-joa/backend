package services

import (
	"backend/entities"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func BoardPost(c *gin.Context, db *gorm.DB) {
	var board entities.BoardDTO
	if err := c.ShouldBindJSON(&board); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "KST 타임존을 설정하는 중 오류가 발생했습니다.",
		})
		return
	}

	var user entities.User
	if err := db.Where("email = ?", board.Email).First(&user).Error; gorm.ErrRecordNotFound == err {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	// 현재 KST 시간을 문자열로 변환
	currentTime := time.Now().In(location).Format("2006-01-02 15:04:05")

	upload := entities.Board{
		Author:      user.Name,
		Title:       board.Title,
		Description: board.Description,
		Date:        currentTime,
	}
	if err := db.Create(&upload).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"model":  upload,
	})
}

func BoardAllGet(c *gin.Context, db *gorm.DB) {
	var boards []entities.Board
	if err := db.Find(&boards).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"data":   boards,
	})
}

func BoardGetByID(c *gin.Context, db *gorm.DB) {
	var board entities.Board
	id := c.Param("id")
	if err := db.Where("id =?", id).First(&board).Error; gorm.ErrRecordNotFound == err {
		c.JSON(404, gin.H{"error": "Board not found"})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"model":  board,
	})
}

func BoardDeleteByID(c *gin.Context, db *gorm.DB) {
	var board entities.Board
	id := c.Param("id")
	if err := db.Where("id =?", id).First(&board).Error; gorm.ErrRecordNotFound == err {
		c.JSON(404, gin.H{"error": "Board not found"})
		return
	}
	if err := db.Delete(&board).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"model":  board,
	})
}
