package models

type ScoreResult struct {
	TotalQuestions   int `json:"total_questions"`
	CorrectAnswers   int `json:"correct_answers"`
	IncorrectAnswers int `json:"incorrect_answers"`
	Score 			 int `json:"score"`
}
