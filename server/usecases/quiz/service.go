package quiz

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/KyleBanks/ably-test/server/drivers/logger"
	"github.com/KyleBanks/ably-test/server/drivers/transport"
	"github.com/KyleBanks/ably-test/server/models"
)

const (
	messageNameQuestion = "question"
	messageNameResults  = "results"
)

type Config struct {
	MinPlayers       int
	RoundDuration    time.Duration
	TimeBetweenGames time.Duration
	QuestionsPath    string
}

type Service struct {
	logger    logger.Logger
	transport transport.Transporter
	config    Config

	questions []models.QuestionDefinition
}

func NewService(
	logger logger.Logger,
	transport transport.Transporter,
	config Config,
) (*Service, error) {
	questionsData, err := os.ReadFile(config.QuestionsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to real questions file: %w", err)
	}

	var questions []models.QuestionDefinition
	err = json.Unmarshal(questionsData, &questions)
	if err != nil {
		return nil, fmt.Errorf("failed to parse questions as JSON: %w", err)
	}

	return &Service{
		logger:    logger,
		transport: transport,
		config:    config,
		questions: questions,
	}, nil
}

// Run the quiz game loop, starting a new quiz after each game completes.
func (s *Service) Run(ctx context.Context) error {
	for {
		game := models.NewGame(s.questions)
		if err := s.runGame(ctx, game); err != nil {
			return fmt.Errorf("error during round: %w", err)
		}
		time.Sleep(s.config.TimeBetweenGames)
	}
}
