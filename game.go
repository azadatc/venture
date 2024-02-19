package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Suit represents a card suit
type Suit int

const (
	Spades Suit = iota
	Hearts
	Diamonds
	Clubs
)

// Value represents a card value
type Value int

const (
	Ace Value = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

// Card represents a single card in the deck
type Card struct {
	Suit  Suit
	Value Value
	Name  string
}

// Deck represents a deck of cards
type Deck struct {
	Cards []Card
}

// Shoe represents a shoe consisting of multiple decks
type Shoe struct {
	Decks []Deck
}

// NewDeck creates a new deck of cards
func NewDeck() Deck {
	var cards []Card
	cardSuits := []string{"Spades", "Diamonds", "Hearts", "Clubs"}
	cardValues := map[string]int{
		"Ace": 1, "Two": 2, "Three": 3, "Four": 4, "Five": 5, "Six": 6, "Seven": 7, "Eight": 8, "Nine": 9,
		"Ten": 10, "Jack": 10, "Queen": 10, "King": 10,
	}
	for _, suit := range cardSuits {
		for cardName, cardValue := range cardValues {
			cards = append(cards, Card{Suit: SuitFromString(suit), Value: Value(cardValue), Name: fmt.Sprintf("%s of %s", cardName, suit)})
		}
	}
	return Deck{Cards: cards}
}

func NewShoe(numDecks int) Shoe {
	shoe := Shoe{}
	for i := 0; i < numDecks; i++ {
		shoe.Decks = append(shoe.Decks, NewDeck())
	}

	return shoe
}

// SuitFromString converts a string to a Suit constant
func SuitFromString(s string) Suit {
	switch s {
	case "Spades":
		return Spades
	case "Diamonds":
		return Diamonds
	case "Hearts":
		return Hearts
	case "Clubs":
		return Clubs
	default:
		return -1
	}
}

// Shuffle shuffles the deck of cards
func (d *Deck) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

// Shuffle shuffles all decks in the shoe
func (s *Shoe) Shuffle() {
	// Collect all cards into a single slice
	var allCards []Card
	for _, deck := range s.Decks {
		allCards = append(allCards, deck.Cards...)
	}

	// Shuffle the combined slice of cards
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(allCards), func(i, j int) {
		allCards[i], allCards[j] = allCards[j], allCards[i]
	})

	// Distribute the shuffled cards back into decks
	cardIndex := 0
	for i := range s.Decks {
		deckSize := len(s.Decks[i].Cards)
		copy(s.Decks[i].Cards, allCards[cardIndex:cardIndex+deckSize])
		cardIndex += deckSize
	}
}

// Draw takes the top card from the shoe, returns it, and removes it from the shoe
func (s *Shoe) Draw() Card {
	// Get the top card from the first deck
	topCard := s.Decks[0].Cards[0]
	// Remove the top card from the deck
	s.Decks[0].Cards = s.Decks[0].Cards[1:]
	// If the deck is empty, remove it from the shoe
	if len(s.Decks[0].Cards) == 0 {
		s.Decks = s.Decks[1:]
	}
	return topCard
}

// SuitToString converts suit constants to string
func SuitToString(s Suit) string {
	suits := [...]string{"Spades", "Hearts", "Diamonds", "Clubs"}
	return suits[s]
}

type Stat struct {
	TotalHandsPlayed    int
	PlayerWinCount      int
	BankerWinCount      int
	TieCount            int
	PlayerMaxWinStreak  int
	PlayerCurrentStreak int
	BankerMaxWinStreak  int
	BankerCurrentStreak int
	TieMaxStreak        int
	TieCurrentStreak    int
	PlayerWonLastHand   bool
	BankerWonLastHand   bool
	LastGameWasTied     bool
}

func (s Stat) GetPlayerWinPercentage() float64 {
	return float64(s.PlayerWinCount) / float64(s.TotalHandsPlayed)
}

func (s Stat) GetBankerWinPercentage() float64 {
	return float64(s.BankerWinCount) / float64(s.TotalHandsPlayed)
}

func (s Stat) GetTiePercentage() float64 {
	return float64(s.TieCount) / float64(s.TotalHandsPlayed)
}

// Define constants for the outcomes of the game
const (
	PlayerWins = iota
	BankerWins
	Tie
)

type BaccaratGame struct {
	Shoe        Shoe
	Player      []Card
	Banker      []Card
	Statistics  *Stat
	canContinue bool
}

// NewBaccaratGame creates a new game of Baccarat with a shuffled shoe
func NewBaccaratGame(numDecks int) *BaccaratGame {
	shoe := NewShoe(numDecks)
	shoe.Shuffle()
	topCard := shoe.Draw()
	//fmt.Printf("Top card is: %d\n", topCard.Value)
	//fmt.Printf("Discarding %d cards\n", topCard.Value)
	for i := 0; i < int(topCard.Value); i++ {
		shoe.Draw()
	}
	return &BaccaratGame{
		Shoe:        shoe,
		canContinue: true,
		Statistics: &Stat{
			TotalHandsPlayed:   0,
			PlayerWinCount:     0,
			BankerWinCount:     0,
			TieCount:           0,
			PlayerMaxWinStreak: 0,
			BankerMaxWinStreak: 0,
			TieMaxStreak:       0,
			PlayerWonLastHand:  false,
			BankerWonLastHand:  false,
		},
	}
}

// DealHands deals two cards to the Player and two cards to the Banker using the following order: Player, Banker, Player, Banker
func (g *BaccaratGame) DealHands() {
	// Check if the shoe has only 1 deck left
	if len(g.Shoe.Decks) == 1 {
		g.canContinue = false
		return
	}
	g.Player = append(g.Player, g.Shoe.Draw())
	g.Banker = append(g.Banker, g.Shoe.Draw())
	g.Player = append(g.Player, g.Shoe.Draw())
	g.Banker = append(g.Banker, g.Shoe.Draw())
}

// EvaluateHand evaluates the value of a hand
func EvaluateHand(hand []Card) int {
	total := 0
	for _, card := range hand {
		total += int(card.Value)
	}
	return total % 10
}

// PlayGame plays a single game of Baccarat according to the specified rules
func (g *BaccaratGame) PlayGame() (int, int, error) {
	// Check if it's possible to continue playing
	if !g.canContinue {
		return -1, -1, fmt.Errorf("the shoe is empty")
	}

	g.Statistics.TotalHandsPlayed++

	// Deal initial hands
	g.DealHands()

	// Evaluate hands
	playerScore := EvaluateHand(g.Player)
	bankerScore := EvaluateHand(g.Banker)

	// Check for naturals
	if playerScore == 8 || playerScore == 9 || bankerScore == 8 || bankerScore == 9 {
		// No action needed for naturals
		// Return the game results
		who, score := g.EvaluateGameResults()
		return who, score, nil
	}

	playerDraws := false
	playerDrew := -1

	// Player draws rules
	if playerScore <= 5 {
		// Player draws a third card
		g.Player = append(g.Player, g.Shoe.Draw())
		playerDraws = true
		playerDrew = int(g.Player[len(g.Player)-1].Value)
	}

	if (bankerScore <= 2) || (bankerScore == 3 && playerDraws && playerDrew != 8) || (bankerScore == 4 && playerDraws && playerDrew >= 2 && playerDrew <= 7) || (bankerScore == 5 && playerDraws && playerDrew >= 4 && playerDrew <= 7) || (bankerScore == 6 && playerDraws && playerDrew >= 6 && playerDrew <= 7) {
		g.Banker = append(g.Banker, g.Shoe.Draw())
	}
	who, score := g.EvaluateGameResults()
	return who, score, nil
}

func (g *BaccaratGame) EvaluateGameResults() (int, int) {
	playerScore := EvaluateHand(g.Player)
	bankerScore := EvaluateHand(g.Banker)

	if playerScore > bankerScore {
		g.Statistics.PlayerWinCount++
		// calculate player win streak and update max win streak
		if g.Statistics.PlayerWonLastHand {
			g.Statistics.PlayerCurrentStreak++
			if g.Statistics.PlayerCurrentStreak > g.Statistics.PlayerMaxWinStreak {
				g.Statistics.PlayerMaxWinStreak = g.Statistics.PlayerCurrentStreak
			}
		} else {
			g.Statistics.PlayerCurrentStreak = 1
			g.Statistics.PlayerWonLastHand = true
			g.Statistics.BankerWonLastHand = false
			g.Statistics.LastGameWasTied = false
			g.Statistics.BankerCurrentStreak = 0
			g.Statistics.TieCurrentStreak = 0
		}

		return PlayerWins, playerScore
	} else if bankerScore > playerScore {
		g.Statistics.BankerWinCount++
		// calculate player win streak and update max win streak
		if g.Statistics.BankerWonLastHand {
			g.Statistics.BankerCurrentStreak++
			if g.Statistics.BankerCurrentStreak > g.Statistics.BankerMaxWinStreak {
				g.Statistics.BankerMaxWinStreak = g.Statistics.BankerCurrentStreak
			}
		} else {
			g.Statistics.BankerCurrentStreak = 1
			g.Statistics.PlayerWonLastHand = false
			g.Statistics.BankerWonLastHand = true
			g.Statistics.LastGameWasTied = false
			g.Statistics.PlayerCurrentStreak = 0
			g.Statistics.TieCurrentStreak = 0
		}
		return BankerWins, bankerScore
	} else {
		g.Statistics.TieCount++
		// calculate player win streak and update max win streak
		if g.Statistics.LastGameWasTied {
			g.Statistics.TieCurrentStreak++
			if g.Statistics.TieCurrentStreak > g.Statistics.TieMaxStreak {
				g.Statistics.TieMaxStreak = g.Statistics.TieCurrentStreak
			}
		} else {
			g.Statistics.BankerCurrentStreak = 0
			g.Statistics.PlayerWonLastHand = false
			g.Statistics.BankerWonLastHand = false
			g.Statistics.LastGameWasTied = true
			g.Statistics.PlayerCurrentStreak = 0
			g.Statistics.TieCurrentStreak = 1
		}
		return Tie, playerScore
	}
}
