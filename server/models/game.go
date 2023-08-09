package models

type QuestionDefinition struct {
	PublicQuestionDefinition
	CorrectAnswerID string `json:"correct_answer"`
}

// PublicQuestionDefinition omits the correct answer of a question and is safe to broadcast to players
type PublicQuestionDefinition struct {
	ID      string             `json:"id"`
	Query   string             `json:"query"`
	Answers []AnswerDefinition `json:"answers"`
}

type AnswerDefinition struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

// ResponseMap for a question maps a clientID to their correct/incorrect state
type ResponseMap map[string]bool

type Game struct {
	Questions []QuestionDefinition

	responses       []ResponseMap
	startingPlayers []string
}

func NewGame(
	questions []QuestionDefinition,
) *Game {
	return &Game{
		Questions: questions,
		responses: make([]ResponseMap, len(questions)),
	}
}

func (g *Game) SetStartingPlayers(players []string) {
	g.startingPlayers = players
}

func (g *Game) AddResponses(questionIdx int, res ResponseMap) {
	g.responses[questionIdx] = res
}
