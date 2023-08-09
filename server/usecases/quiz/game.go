package quiz

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/KyleBanks/ably-test/server/drivers/transport"
	"github.com/KyleBanks/ably-test/server/models"
)

func (s *Service) runGame(
	ctx context.Context,
	game *models.Game,
) error {
	s.logger.Info(ctx, "waiting for %d players", s.config.MinPlayers)
	err := s.waitForPlayers(ctx, game)
	if err != nil {
		return fmt.Errorf("failed to wait for numPlayer count: %w", err)
	}

	s.logger.Info(ctx, "game started")
	for idx, question := range game.Questions {
		res, err := s.runQuestion(ctx, question, s.config.RoundDuration)
		if err != nil {
			return fmt.Errorf("failed to score responses: %w", err)
		}
		game.AddResponses(idx, res)
	}

	s.logger.Info(ctx, "game complete")
	err = s.publishResults(ctx, game)
	if err != nil {
		return fmt.Errorf("failed to publish results: %w", err)
	}

	return nil
}

func (s *Service) waitForPlayers(
	ctx context.Context,
	game *models.Game,
) error {
	var players []string
	var err error

	var lastCount int
	for lastCount < s.config.MinPlayers {
		players, err = s.transport.ClientIDs(ctx)
		if err != nil {
			return fmt.Errorf("failed to check player count of game: %w", err)
		}

		numPlayers := len(players)
		if numPlayers != lastCount {
			s.logger.Info(ctx, "player count changed; from=%d, to=%d", lastCount, numPlayers)
			lastCount = numPlayers
		}

		time.Sleep(time.Millisecond * 100)
	}

	game.SetStartingPlayers(players)
	return nil
}

func (s *Service) runQuestion(
	ctx context.Context,
	question models.QuestionDefinition,
	duration time.Duration,
) (models.ResponseMap, error) {

	// Broadcast the question to all players
	err := s.transport.Publish(ctx, messageNameQuestion, question.PublicQuestionDefinition)
	if err != nil {
		return nil, fmt.Errorf("failed to publish question: %w", err)
	}

	res := make(models.ResponseMap)
	// Subscribe for the duration of the question, using the question ID as the
	// name to filter on. When the question duration elapses any late responses
	// to this question will be ignored.
	done, err := s.transport.Subscribe(ctx, question.ID, func(m *transport.Message) {
		answer, ok := m.Data.(string)
		if !ok {
			s.logger.Error(ctx, "failed to interpret response '%v', discarding...", m.Data)
			return
		}

		correct := strings.EqualFold(answer, question.CorrectAnswerID)
		s.logger.Info(ctx, "answer '%s' from clientId=%s is correct=%v", answer, m.ClientID, correct)

		// Messages are ordered so its safe to overwrite an existing result if the player
		// changes their answer.
		//
		// Not checking if this ClientID exists in the game, just allowing new players to
		// join mid-game for now.
		//
		// Ably SDK guarantees this callback will be invoked sequentially rather than concurrently,
		// no need for a mutex.
		res[m.ClientID] = correct
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to responses: %w", err)
	}
	defer done()

	time.Sleep(duration)
	return res, nil
}

func (s *Service) publishResults(
	ctx context.Context,
	game *models.Game,
) error {
	results := models.NewResults(game)
	if err := s.transport.Publish(ctx, messageNameResults, results); err != nil {
		return fmt.Errorf("failed to publish results: %w", err)
	}
	return nil
}
