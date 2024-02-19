package main

import (
	"fmt"
)

func JustRunTheGame(shoeCount int) {

	for i := 0; i < shoeCount; i++ {
		fmt.Println("Starting a new game - #", i+1)
		fmt.Println("---------------------------------------------------------------------------")
		game := NewBaccaratGame(6)

		// Play the game until the shoe is empty
		for {
			_, _, err := game.PlayGame()

			// If there's an error, print it and break the loop
			if err != nil {
				fmt.Println(err)
				fmt.Printf("Total hands played: %d\n", game.Statistics.TotalHandsPlayed)
				fmt.Printf("Player wins: %d\n", game.Statistics.PlayerWinCount)
				fmt.Printf("Banker wins: %d\n", game.Statistics.BankerWinCount)
				fmt.Printf("Ties: %d\n", game.Statistics.TieCount)
				fmt.Printf("Player max win streak: %d\n", game.Statistics.PlayerMaxWinStreak)
				fmt.Printf("Banker max win streak: %d\n", game.Statistics.BankerMaxWinStreak)
				fmt.Printf("Tie max streak: %d\n", game.Statistics.TieMaxStreak)
				fmt.Printf("Player win percentage: %.2f%%\n", game.Statistics.GetPlayerWinPercentage()*100)
				fmt.Printf("Banker win percentage: %.2f%%\n", game.Statistics.GetBankerWinPercentage()*100)
				fmt.Printf("Tie percentage: %.2f%%\n", game.Statistics.GetTiePercentage()*100)

				break
			}
		}
		fmt.Printf("\n---------------------------------------------------------------------------\n\n\n")
	}
}

func main() {
	fmt.Printf("TotalHands,WinCount,LossCount,TieCount,MaxWinStreak,MaxLossStreak,MaxTieStreak,WinPercentage,LossPercentage,TiePercentage,MaxBet,Bankroll\n")
	Player := NewPlayer(nil, 500000, 10, 5000000)

	for i := 0; i < 1500; i++ {
		//fmt.Println("Starting a new game - #", i+1)
		//fmt.Println("---------------------------------------------------------------------------")
		baccaratGame := NewBaccaratGame(6)
		Player.StartNewGame(baccaratGame, true)

		Player.PlayTheGame()
		Player.PrintStats()

		//fmt.Printf("\n---------------------------------------------------------------------------\n\n\n")
	}
}
