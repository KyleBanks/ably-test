package models

// PlayerScore maps a clientID to their score
type PlayerScore map[string]int

type Results struct {
	Scores PlayerScore `json:"scores"`
}

func NewResults(
	g *Game,
) *Results {
	r := Results{
		Scores: make(PlayerScore),
	}

	for _, responses := range g.responses {
		for clientId, correct := range responses {
			if !correct {
				continue
			}

			score, _ := r.Scores[clientId]
			score++
			r.Scores[clientId] = score
		}
	}

	// starting players are always included in the results, even if they never submitted a response
	for _, clientId := range g.startingPlayers {
		if _, ok := r.Scores[clientId]; !ok {
			r.Scores[clientId] = 0
		}
	}

	return &r
}
