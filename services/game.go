package services

import (
	"backend/entities"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GamePOST(c *gin.Context, mongo *mongo.Client) {
	var game entities.GameDTO
	if err := c.ShouldBindJSON(&game); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	upload := entities.Game{
		ID:      game.ID,
		Title:   game.Title,
		Answer1: game.Answer1,
		Answer2: game.Answer2,
	}

	collection := mongo.Database("neoga-joa").Collection("game")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // 중요: 함수가 끝날 때 컨텍스트를 취소하여 리소스 누출 방지

	result, err := collection.InsertOne(ctx, upload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data into MongoDB"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

func GameALLGET(c *gin.Context, db *mongo.Client) {
	collection := db.Database("neoga-joa").Collection("game")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	var results []entities.Game
	if err := cur.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
	})
}

func GameGETByID(c *gin.Context, mongo *mongo.Client) {
	id := c.Param("id")

	// MongoDB 컬렉션 선택
	collection := mongo.Database("neoga-joa").Collection("game")

	// 컨텍스트 생성
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 카테고리에 대한 필터 생성
	filter := bson.D{{"id", id}}

	// 컬렉션에서 데이터 가져오기
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(ctx)

	// 결과를 담을 슬라이스 선언
	var results []entities.Game

	// 결과를 슬라이스에 저장
	if err := cur.All(ctx, &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 결과를 JSON으로 응답
	c.JSON(http.StatusOK, gin.H{"results": results})
}
func Answer1(c *gin.Context, mongo *mongo.Client, db *gorm.DB) {
	var userInfo struct {
		Email  string `json:"email"`
		Reason string `json:"reason"`
	}
	if err := c.BindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user information"})
		return
	}

	id := c.Param("id")

	// MongoDB 컬렉션 및 쿼리용 컨텍스트 생성
	collection := mongo.Database("neoga-joa").Collection("game")
	ctx := context.TODO()

	// MongoDB에서 해당 ID에 해당하는 객체 검색
	filter := bson.M{"id": id} // 클라이언트에서 전달받은 id를 기준으로 검색
	var game entities.Game
	err := collection.FindOne(ctx, filter).Decode(&game)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var user entities.User
	if err := db.Where("email = ?", userInfo.Email).First(&user).Error; gorm.ErrRecordNotFound == err {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 사용자의 이름과 커플 코드를 구조체로 생성하여 Answer1People에 추가
	answerPerson := entities.AnswerPerson{
		Name:       user.Name,
		CoupleCode: utils.ToString(user.CoupleCode),
		Reason:     userInfo.Reason,
	}
	game.Answer1People = append(game.Answer1People, answerPerson)

	// MongoDB에 업데이트된 객체 저장
	_, err = collection.ReplaceOne(ctx, filter, game)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User added to Answer1People"})
}

func Answer2(c *gin.Context, mongo *mongo.Client, db *gorm.DB) {
	var userInfo struct {
		Email  string `json:"email"`
		Reason string `json:"reason"`
	}
	if err := c.BindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user information"})
		return
	}

	id := c.Param("id")

	// MongoDB 컬렉션 및 쿼리용 컨텍스트 생성
	collection := mongo.Database("neoga-joa").Collection("game")
	ctx := context.TODO()

	// MongoDB에서 해당 ID에 해당하는 객체 검색
	filter := bson.M{"id": id} // 클라이언트에서 전달받은 id를 기준으로 검색
	var game entities.Game
	err := collection.FindOne(ctx, filter).Decode(&game)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var user entities.User
	if err := db.Where("email = ?", userInfo.Email).First(&user).Error; gorm.ErrRecordNotFound == err {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 사용자의 이름과 커플 코드를 구조체로 생성하여 Answer1People에 추가
	answerPerson := entities.AnswerPerson{
		Name:       user.Name,
		CoupleCode: utils.ToString(user.CoupleCode),
		Reason:     userInfo.Reason,
	}
	game.Answer2People = append(game.Answer2People, answerPerson)

	// MongoDB에 업데이트된 객체 저장
	_, err = collection.ReplaceOne(ctx, filter, game)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User added to Answer2People"})
}

func CheckAnswer(c *gin.Context, mongo *mongo.Client, db *gorm.DB) {
	// 클라이언트로부터 전달받은 ID를 파라미터에서 추출
	id := c.Param("id")

	var userInfo struct {
		Email string `json:"email"`
	}
	if err := c.BindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user information"})
		return
	}

	// MongoDB 컬렉션 선택
	collection := mongo.Database("neoga-joa").Collection("game")
	ctx := context.TODO()

	// 이메일로부터 사용자 정보 가져오기
	var user entities.User
	if err := db.Where("email = ?", userInfo.Email).First(&user).Error; gorm.ErrRecordNotFound == err {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// MongoDB에서 해당 ID에 해당하는 게임 객체 검색
	filter := bson.D{{"id", id}}
	var game entities.Game
	if err := collection.FindOne(ctx, filter).Decode(&game); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 사용자의 이름과 커플 코드
	userCoupleCode := strconv.Itoa(user.CoupleCode)

	// 답변1에 커플 코드가 들어있는 사람들의 이름과 이유를 검사
	var answer1DifferentNames []gin.H
	for _, person := range game.Answer1People {
		if person.CoupleCode != userCoupleCode {
			// 다른 사람의 이름과 이유를 반환
			answer1DifferentNames = append(answer1DifferentNames, gin.H{"name": person.Name, "reason": person.Reason})
		}
	}

	// 답변2에 커플 코드가 들어있는 사람들의 이름을 검사
	var answer2DifferentNames []gin.H
	for _, person := range game.Answer2People {
		if person.CoupleCode != userCoupleCode {
			// 다른 사람의 이름과 이유를 반환
			answer2DifferentNames = append(answer2DifferentNames, gin.H{"name": person.Name, "reason": person.Reason})
		}
	}

	// 다른 이름이 없으면 메시지 추가
	message := ""
	if len(answer1DifferentNames) == 0 && len(answer2DifferentNames) == 0 {
		message = "No names found for both answers"
	}

	// 결과 반환
	responseData := gin.H{
		"answer1_different_names": answer1DifferentNames,
		"answer2_different_names": answer2DifferentNames,
		"message":                 message,
	}
	c.JSON(http.StatusOK, responseData)
}
