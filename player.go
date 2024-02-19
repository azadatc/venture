package main

import (
	"fmt"
	"math/rand"
	"time"
)

type PlayerStats struct {
	TotalHandsPlayed  int
	WinCount          int
	LossCount         int
	TieCount          int
	MaxBet			int
	CurrentWinStreak  int
	CurrentLossStreak int
	CurrentTieStreak  int
	MaxWinStreak      int
	MaxLossStreak     int
	MaxTieStreak      int
}

func (p *PlayerStats) GetWinPercentage() float64 {
	return float64(p.WinCount) / float64(p.TotalHandsPlayed)
}

func (p *PlayerStats) GetLossPercentage() float64 {
	return float64(p.LossCount) / float64(p.TotalHandsPlayed)
}

func (p *PlayerStats) GetTiePercentage() float64 {
	return float64(p.TieCount) / float64(p.TotalHandsPlayed)
}

func (s *PlayerStats) ResetStats() {
	s.CurrentLossStreak = 0
	s.CurrentWinStreak = 0
	s.CurrentTieStreak = 0
	s.TotalHandsPlayed = 0
	s.WinCount = 0
	s.LossCount = 0
	s.TieCount = 0
	s.MaxWinStreak = 0
	s.MaxLossStreak = 0
	s.MaxTieStreak = 0
	s.MaxBet = 0
}

type Player struct {
	BaccaratGame      *BaccaratGame
	InitialInvestment int
	CurrentBankroll   int
	MinimumBet        int
	CurrentBet        int
	WonLastHand       bool
	TiedLastHand      bool
	MaxBet            int
	Stats             *PlayerStats
}

func NewPlayer(game *BaccaratGame, initialInvestment, minimumBet, maxBet int) *Player {
	return &Player{
		BaccaratGame:      game,
		InitialInvestment: initialInvestment,
		CurrentBankroll:   initialInvestment,
		MinimumBet:        minimumBet,
		MaxBet:            maxBet,
		CurrentBet:        minimumBet,
		WonLastHand:       true,
		Stats: &PlayerStats{
			TotalHandsPlayed: 0,
			WinCount:         0,
			LossCount:        0,
			TieCount:         0,
			MaxWinStreak:     0,
			MaxLossStreak:    0,
			MaxTieStreak:     0,
			MaxBet:           0,
		},
	}
}

func (p *Player) StartNewGame(game *BaccaratGame, startWithMinimumBet bool) {
	p.BaccaratGame = game
	p.Stats.ResetStats()
	if startWithMinimumBet {
		p.WonLastHand = true
	}
}

func (p Player) PrintStats() {
	// fmt.Printf("Total hands played: %d\n", p.Stats.TotalHandsPlayed)
	// fmt.Printf("Player wins: %d\n", p.Stats.WinCount)
	// fmt.Printf("Player loses: %d\n", p.Stats.LossCount)
	// fmt.Printf("Player ties: %d\n", p.Stats.TieCount)
	// fmt.Printf("Player max win streak: %d\n", p.Stats.MaxWinStreak)
	// fmt.Printf("Player max loss streak: %d\n", p.Stats.MaxLossStreak)
	// fmt.Printf("Player max tie streak: %d\n", p.Stats.MaxTieStreak)
	// fmt.Printf("Player win percentage: %.2f%%\n", p.Stats.GetWinPercentage()*100)
	// fmt.Printf("Player loss percentage: %.2f%%\n", p.Stats.GetLossPercentage()*100)
	// fmt.Printf("Player tie percentage: %.2f%%\n", p.Stats.GetTiePercentage()*100)

	// fmt.Printf("Current bankroll: %d\n", p.CurrentBankroll)

	//fmt.Printf("TotalHands\tWinCount\tLossCount\tTieCount\tMaxWinStreak\tMaxLossStreak\tMaxTieStreak\tWinPercentage\tLossPercentage\tTiePercentage\tBankroll\n")
	fmt.Printf("%d,%d,%d,%d,%d,%d,%d,%.2f%%,%.2f%%,%.2f%%,%d,%d\n",
		p.Stats.TotalHandsPlayed,
		p.Stats.WinCount,
		p.Stats.LossCount,
		p.Stats.TieCount,
		p.Stats.MaxWinStreak,
		p.Stats.MaxLossStreak,
		p.Stats.MaxTieStreak,
		p.Stats.GetWinPercentage()*100,
		p.Stats.GetLossPercentage()*100,
		p.Stats.GetTiePercentage()*100,
		p.Stats.MaxBet,
		p.CurrentBankroll)
}

func (p *Player) Guess() int {
	// Generate a random int, either 0 or 1
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return r.Intn(2)
}

func (p *Player) PlayTheGame() {
	for {
		if !p.BaccaratGame.canContinue {
			break
		}

		if p.CurrentBankroll < p.MinimumBet {
			//fmt.Println("You're out of money!")
			break
		}

		// calculate what the next bet should be
		if p.WonLastHand {
			p.CurrentBet = p.MinimumBet
		} else if p.TiedLastHand {
		} else {
			p.CurrentBet = p.CurrentBet * 2
		}

		if p.CurrentBet > p.MaxBet {
			p.CurrentBet = p.MaxBet
		}

		if p.CurrentBet > p.CurrentBankroll {
			p.CurrentBet = p.CurrentBankroll
		}

		if p.CurrentBet > p.Stats.MaxBet {
			p.Stats.MaxBet = p.CurrentBet
		}

		// Place the bet
		//fmt.Printf("Placing bet of %d\n", p.CurrentBet)

		p.CurrentBankroll -= p.CurrentBet

		// Play the game
		p.Stats.TotalHandsPlayed++
		playerGuess := p.Guess()
		res, _, _ := p.BaccaratGame.PlayGame()
		//fmt.Printf("Current bankroll: %d\n", p.CurrentBankroll)
		//fmt.Printf("Player guessed %d, result was %d\n", playerGuess, res)

		if res == playerGuess {
			// player wins
			//fmt.Printf("Player wins %d\n", p.CurrentBet*2)
			p.CurrentBankroll += p.CurrentBet * 2
			p.WonLastHand = true
			p.TiedLastHand = false
			p.Stats.WinCount++
			p.Stats.CurrentWinStreak++
			if p.Stats.CurrentWinStreak > p.Stats.MaxWinStreak {
				p.Stats.MaxWinStreak = p.Stats.CurrentWinStreak
			}
			p.Stats.CurrentLossStreak = 0
			p.Stats.CurrentTieStreak = 0
		} else if res == Tie {
			//fmt.Println("Player ties")
			p.CurrentBankroll += p.CurrentBet
			p.WonLastHand = false
			p.TiedLastHand = true
			p.Stats.TieCount++
			p.Stats.CurrentTieStreak++
			if p.Stats.CurrentTieStreak > p.Stats.MaxTieStreak {
				p.Stats.MaxTieStreak = p.Stats.CurrentTieStreak
			}
		} else {
			// player loses
			//fmt.Printf("Player loses %d\n", p.CurrentBet)
			p.WonLastHand = false
			p.TiedLastHand = false
			p.Stats.LossCount++
			p.Stats.CurrentLossStreak++
			if p.Stats.CurrentLossStreak > p.Stats.MaxLossStreak {
				p.Stats.MaxLossStreak = p.Stats.CurrentLossStreak
			}
			p.Stats.CurrentWinStreak = 0
			p.Stats.CurrentTieStreak = 0
		}

		//fmt.Printf("Current bankroll: %d\n", p.CurrentBankroll)
	}
}
