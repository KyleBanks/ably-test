package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResults(t *testing.T) {
	type testCase struct {
		in     *Game
		expect Results
	}

	tests := map[string]testCase{
		"correctly tallies scores": {
			in: &Game{
				responses: []ResponseMap{
					{
						"P1": true,
						"P2": true,
					},
					{
						"P1": false,
						"P2": true,
					},
					{
						"P1": true,
						"P2": true,
					},
				},
				startingPlayers: []string{"P1", "P2"},
			},
			expect: Results{
				Scores: PlayerScore{
					"P1": 2,
					"P2": 3,
				},
			},
		},
		"returns empty on empty input": {
			in: &Game{},
			expect: Results{
				Scores: make(PlayerScore),
			},
		},
		"always includes starting players": {
			in: &Game{
				responses: []ResponseMap{
					{
						"P1": true,
						"P2": true,
					},
					{
						"P1": false,
					},
					{
						"P1": true,
					},
				},
				startingPlayers: []string{"P1", "P2", "P3"},
			},
			expect: Results{
				Scores: PlayerScore{
					"P1": 2,
					"P2": 1,
					"P3": 0,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res := NewResults(tc.in)
			assert.Equal(t, tc.expect, *res)
		})
	}
}
