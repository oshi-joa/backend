package entities

type AnswerPerson struct {
	Name       string `json:"name"`
	CoupleCode string `json:"couple_code"`
	Reason     string `json:"reason"`
}

type Game struct {
	ID            string
	Title         string
	Answer1       string
	Answer2       string
	Answer1People []AnswerPerson
	Answer2People []AnswerPerson
}

type GameDTO struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Answer1 string `json:"answer1"`
	Answer2 string `json:"answer2"`
}
