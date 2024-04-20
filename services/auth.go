package services

import (
	"backend/entities"
	"crypto/sha1"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

func CreateToken(userId string) (string, error) {
	var err error

	secretKey := "asdooinvzxcubuwebdcs" // 이 키는 보안을 위해 안전하게 보관해야 합니다.

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Hour * 24).Unix() // 토큰 유효 시간은 24시간으로 설정했습니다. 필요에 따라 변경하세요.
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	token, err := at.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func SignUp(c *gin.Context, db *gorm.DB) {
	var user entities.UserDTO
	c.BindJSON(&user)

	// 이메일이 이미 존재하는지 확인
	var existingUser entities.User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		// 이미 존재하는 이메일인 경우
		c.JSON(400, gin.H{"error": "Email already exists"})
		return
	}

	hashedPassword := hashPassword(user.Password)

	// 6자리 임의의 숫자 생성
	uuid := generateUUID()

	// 사용자 정보 생성
	upload := &entities.User{
		Uuid:       uuid,
		Name:       user.Name,
		Email:      user.Email,
		Password:   hashedPassword,
		CoupleCode: user.CoupleCode,
	}

	db.Create(&upload)
	c.JSON(200, gin.H{
		"status": "success",
		"data":   user,
	})
}

func hashPassword(password string) string {
	// 비밀번호를 해시화하여 반환
	hashed := sha1.New()
	hashed.Write([]byte(password))
	return fmt.Sprintf("%x", hashed.Sum(nil))
}

func generateUUID() int {
	// 랜덤 숫자 생성을 위한 시드 설정
	rand.Seed(time.Now().UnixNano())

	// 6자리의 임의의 숫자 생성
	min := 100000
	max := 999999
	return rand.Intn(max-min+1) + min
}

func Login(c *gin.Context, db *gorm.DB) {
	var loginInfo struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 클라이언트로부터 이메일과 비밀번호 받기
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 이메일로 사용자 찾기
	var user entities.User
	result := db.Where("email = ?", loginInfo.Email).First(&user)
	if result.Error != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// 비밀번호 일치 확인
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// 사용자 인증이 성공하면 JWT 토큰 생성
	token, err := CreateToken(user.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create token"})
		return
	}

	// 토큰을 클라이언트에게 반환
	c.JSON(200, gin.H{"token": token})
}
