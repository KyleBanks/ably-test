package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/KyleBanks/ably-test/server/drivers/logger/stdout"
	"github.com/KyleBanks/ably-test/server/drivers/transport/ably"
	"github.com/KyleBanks/ably-test/server/usecases/quiz"
)

const ablyChannel = "quiz"

// TODO: externalise config
var quizConfig = quiz.Config{
	MinPlayers:       2,
	QuestionsPath:    "./server/config/questions.json",
	RoundDuration:    time.Second * 10,
	TimeBetweenGames: time.Second * 10,
}

func main() {
	ctx := context.Background()
	var logger stdout.Logger

	ablyAPIKey, ok := os.LookupEnv("ABLY_API_KEY")
	if !ok {
		fmt.Println("you must set the ABLY_API_KEY environment variable")
		os.Exit(1)
		return
	}

	transport, err := ably.NewTransport(ctx, ablyAPIKey, ablyChannel)
	if err != nil {
		panic(err)
	}

	quizService, err := quiz.NewService(logger, transport, quizConfig)
	if err != nil {
		panic(err)
	}

	if err := quizService.Run(ctx); err != nil {
		panic(err)
	}
}
