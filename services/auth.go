package services

import (
	"backend/entities"
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
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// 이메일이 이미 존재하는지 확인
	var existingUser entities.User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		// 이미 존재하는 이메일인 경우
		c.JSON(400, gin.H{"error": "Email already exists"})
		return
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	// 6자리 임의의 숫자 (code) 생성
	code := generateUUID()
	// coupleCode 입력이 비어있으면 code로 설정
	// 6자리 임의의 숫자 (code) 생성
	// coupleCode 입력이 비어있으면 code로 설정
	if user.CoupleCode == 0 {
		user.CoupleCode = code // 여기를 수정했습니다. 직접 할당합니다.
	}

	// 사용자 정보 생성
	upload := &entities.User{
		Name:       user.Name,
		Email:      user.Email,
		Password:   hashedPassword,
		CoupleCode: user.CoupleCode,
	}

	db.Create(&upload)
	c.JSON(200, gin.H{
		"status": "success",
		"data":   upload,
	})
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	c.BindJSON(&loginInfo)

	// 이메일을 기반으로 사용자 정보 조회
	var user entities.User
	if err := db.Where("email = ?", loginInfo.Email).First(&user).Error; gorm.ErrRecordNotFound == err {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// 비밀번호 검증
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid password"})
		return
	}

	// JWT 토큰 생성
	token, err := CreateToken(user.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create token"})
		return
	}

	// 로그인 성공 및 토큰 반환
	c.JSON(200, gin.H{
		"status": "success",
		"token":  token,
	})
}

var tokenBlacklist = make(map[string]bool)

// 로그아웃 함수
func Logout(c *gin.Context) {
	// 클라이언트로부터 토큰을 받음
	token := c.Request.Header.Get("Authorization")

	// 토큰이 블랙리스트에 있는지 확인하고 추가함
	tokenBlacklist[token] = true

	// 로그아웃 성공 메시지 반환
	c.JSON(200, gin.H{"message": "Successfully logged out"})
}

func CheckCode(c *gin.Context, db *gorm.DB) {
	var loginInfo struct {
		Email string `json:"email"`
	}
	c.BindJSON(&loginInfo)

	var user entities.User
	if err := db.Where("email = ?", loginInfo.Email).First(&user).Error; gorm.ErrRecordNotFound == err {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"code":   user.CoupleCode,
	})
}
