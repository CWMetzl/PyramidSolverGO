package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// ExposedCard represents an exposed pyramid card with its coordinates, raw string, and numeric value.
type ExposedCard struct {
	row, col int
	card     string
	value    int
}

// GameState represents the current state of the game.
type GameState struct {
	Pyramid [][]string // Pyramid: 7 rows; a removed card is represented by ""
	Deck    []string   // Remaining cards in the stock (draw pile)
	Waste   []string   // Cards drawn from the deck (waste pile)
	Moves   []string   // Log of moves taken so far
}

// Result represents a (partial or complete) solution.
// RemovedCount is the number of pyramid cards removed (max 28).
type Result struct {
	Moves        []string
	RemovedCount int
}

// getCardValue returns the numeric value of a card.
// Ace counts as 1; numbers as themselves; j=11, q=12, k=13.
func getCardValue(card string) int {
	if strings.HasPrefix(card, "10") {
		return 10
	}
	rank := strings.ToLower(string(card[0]))
	switch rank {
	case "a":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "5":
		return 5
	case "6":
		return 6
	case "7":
		return 7
	case "8":
		return 8
	case "9":
		return 9
	case "j":
		return 11
	case "q":
		return 12
	case "k":
		return 13
	}
	return 0
}

// formatCard converts an input card (e.g. "7s", "10c", "kh") into a full name (e.g. "7 of Spades").
func formatCard(card string) string {
	if card == "" || card == "XX" {
		return "Empty"
	}
	var rank, suitChar string
	if strings.HasPrefix(card, "10") {
		rank = "10"
		if len(card) > 2 {
			suitChar = string(card[2])
		}
	} else {
		rank = strings.ToLower(string(card[0]))
		suitChar = string(card[len(card)-1])
	}

	// Map rank to full name.
	rankName := rank
	switch rank {
	case "a":
		rankName = "Ace"
	case "2":
		rankName = "2"
	case "3":
		rankName = "3"
	case "4":
		rankName = "4"
	case "5":
		rankName = "5"
	case "6":
		rankName = "6"
	case "7":
		rankName = "7"
	case "8":
		rankName = "8"
	case "9":
		rankName = "9"
	case "10":
		rankName = "10"
	case "j":
		rankName = "Jack"
	case "q":
		rankName = "Queen"
	case "k":
		rankName = "King"
	}

	// Map suit letter to full name.
	suitName := ""
	switch strings.ToLower(suitChar) {
	case "c":
		suitName = "Clubs"
	case "d":
		suitName = "Diamonds"
	case "h":
		suitName = "Hearts"
	case "s":
		suitName = "Spades"
	default:
		suitName = suitChar
	}

	return fmt.Sprintf("%s of %s", rankName, suitName)
}

// checkDeck verifies that the provided card list represents a full 52‑card deck.
// If a card is missing or appears more than once, it returns an error indicating the issue.
func checkDeck(cards []string) error {
	if len(cards) != 52 {
		return fmt.Errorf("deck must contain 52 cards, but found %d", len(cards))
	}
	cardCount := make(map[string]int)
	for _, card := range cards {
		cLower := strings.ToLower(card)
		cardCount[cLower]++
	}
	// Standard deck: each rank with each suit.
	ranks := []string{"a", "2", "3", "4", "5", "6", "7", "8", "9", "10", "j", "q", "k"}
	suits := []string{"c", "d", "h", "s"}
	for _, rank := range ranks {
		for _, suit := range suits {
			card := rank + suit
			if count, ok := cardCount[card]; !ok {
				return fmt.Errorf("missing card: %s", formatCard(card))
			} else if count > 1 {
				return fmt.Errorf("duplicate card: %s (appears %d times)", formatCard(card), count)
			}
		}
	}
	return nil
}

// buildPyramid constructs the pyramid using the first 28 cards.
func buildPyramid(cards []string) [][]string {
	pyramid := make([][]string, 7)
	index := 0
	for row := 0; row < 7; row++ {
		pyramid[row] = make([]string, row+1)
		for col := 0; col <= row; col++ {
			pyramid[row][col] = cards[index]
			index++
		}
	}
	return pyramid
}

// getExposedCards returns a slice of all exposed pyramid cards.
// A card is exposed if it is not removed ("") and either is on the bottom row
// or both of the cards directly beneath it have been removed.
func getExposedCards(pyramid [][]string) []ExposedCard {
	var exposed []ExposedCard
	for r := 0; r < len(pyramid); r++ {
		for c := 0; c < len(pyramid[r]); c++ {
			if pyramid[r][c] != "" {
				// Bottom row is always exposed.
				if r == len(pyramid)-1 || (pyramid[r+1][c] == "" && pyramid[r+1][c+1] == "") {
					exposed = append(exposed, ExposedCard{
						row:   r,
						col:   c,
						card:  pyramid[r][c],
						value: getCardValue(pyramid[r][c]),
					})
				}
			}
		}
	}
	return exposed
}

// isPyramidEmpty returns true if all cards in the pyramid have been removed.
func isPyramidEmpty(pyramid [][]string) bool {
	for _, row := range pyramid {
		for _, card := range row {
			if card != "" {
				return false
			}
		}
	}
	return true
}

// countRemoved counts how many pyramid cards have been removed.
func countRemoved(pyramid [][]string) int {
	count := 0
	for _, row := range pyramid {
		for _, card := range row {
			if card == "" {
				count++
			}
		}
	}
	return count
}

// clonePyramid creates a deep copy of the pyramid.
func clonePyramid(pyramid [][]string) [][]string {
	newPyramid := make([][]string, len(pyramid))
	for i, row := range pyramid {
		newRow := make([]string, len(row))
		copy(newRow, row)
		newPyramid[i] = newRow
	}
	return newPyramid
}

// cloneState creates a deep copy of a GameState.
func cloneState(state GameState) GameState {
	newState := GameState{
		Pyramid: clonePyramid(state.Pyramid),
		Deck:    make([]string, len(state.Deck)),
		Waste:   make([]string, len(state.Waste)),
		Moves:   make([]string, len(state.Moves)),
	}
	copy(newState.Deck, state.Deck)
	copy(newState.Waste, state.Waste)
	copy(newState.Moves, state.Moves)
	return newState
}

// serializeState creates a string key representing the current state.
// It serializes the pyramid (using "XX" for removed cards), deck, and waste.
func serializeState(state GameState) string {
	var sb strings.Builder
	for i, row := range state.Pyramid {
		for j, card := range row {
			if card == "" {
				sb.WriteString("XX")
			} else {
				sb.WriteString(card)
			}
			if j < len(row)-1 {
				sb.WriteString(",")
			}
		}
		if i < len(state.Pyramid)-1 {
			sb.WriteString(";")
		}
	}
	sb.WriteString("|")
	sb.WriteString(strings.Join(state.Deck, ","))
	sb.WriteString("|")
	sb.WriteString(strings.Join(state.Waste, ","))
	return sb.String()
}

// solveState recursively explores moves from the current state and returns the best Result (partial or complete).
// The visited map keys states to the maximum number of pyramid cards removed seen so far.
func solveState(state GameState, visited map[string]int) Result {
	curRemoved := countRemoved(state.Pyramid)
	bestResult := Result{
		Moves:        state.Moves,
		RemovedCount: curRemoved,
	}
	// If pyramid is completely cleared, return immediately.
	if isPyramidEmpty(state.Pyramid) {
		return bestResult
	}

	key := serializeState(state)
	if val, ok := visited[key]; ok && val >= curRemoved {
		// Already seen this state with at least as many removals.
		return bestResult
	}
	visited[key] = curRemoved

	// 1. Try pyramid removal moves.
	exposed := getExposedCards(state.Pyramid)
	// Sort exposed cards so that lower rows (bottom-up) come first.
	sort.Slice(exposed, func(i, j int) bool {
		return exposed[i].row > exposed[j].row
	})

	// 1a. Remove any exposed King (value 13) from the pyramid.
	for _, exp := range exposed {
		if exp.value == 13 {
			newState := cloneState(state)
			newState.Pyramid[exp.row][exp.col] = ""
			newState.Moves = append(newState.Moves,
				fmt.Sprintf("Remove %s from pyramid at (%d,%d)", formatCard(exp.card), exp.row, exp.col))
			res := solveState(newState, visited)
			if res.RemovedCount > bestResult.RemovedCount {
				bestResult = res
			}
			if bestResult.RemovedCount == 28 {
				return bestResult
			}
		}
	}

	// 1b. Remove any two exposed pyramid cards that add to 13.
	for i := 0; i < len(exposed); i++ {
		for j := i + 1; j < len(exposed); j++ {
			if exposed[i].value+exposed[j].value == 13 {
				newState := cloneState(state)
				newState.Pyramid[exposed[i].row][exposed[i].col] = ""
				newState.Pyramid[exposed[j].row][exposed[j].col] = ""
				newState.Moves = append(newState.Moves,
					fmt.Sprintf("Remove pair from pyramid: %s at (%d,%d) and %s at (%d,%d)",
						formatCard(exposed[i].card), exposed[i].row, exposed[i].col,
						formatCard(exposed[j].card), exposed[j].row, exposed[j].col))
				res := solveState(newState, visited)
				if res.RemovedCount > bestResult.RemovedCount {
					bestResult = res
				}
				if bestResult.RemovedCount == 28 {
					return bestResult
				}
			}
		}
	}

	// 2. Try removal moves using the waste.
	if len(state.Waste) > 0 {
		topWaste := state.Waste[len(state.Waste)-1]
		wasteValue := getCardValue(topWaste)
		for _, exp := range exposed {
			if wasteValue+exp.value == 13 {
				newState := cloneState(state)
				// Remove the top waste card.
				newState.Waste = newState.Waste[:len(newState.Waste)-1]
				newState.Pyramid[exp.row][exp.col] = ""
				newState.Moves = append(newState.Moves,
					fmt.Sprintf("Remove waste card %s and pyramid card %s at (%d,%d)",
						formatCard(topWaste), formatCard(exp.card), exp.row, exp.col))
				res := solveState(newState, visited)
				if res.RemovedCount > bestResult.RemovedCount {
					bestResult = res
				}
				if bestResult.RemovedCount == 28 {
					return bestResult
				}
			}
		}
	}

	// 3. Draw a card from the deck if available.
	if len(state.Deck) > 0 {
		newState := cloneState(state)
		drawn := newState.Deck[0]
		newState.Deck = newState.Deck[1:]
		// If the drawn card is a King, remove it immediately.
		if getCardValue(drawn) == 13 {
			newState.Moves = append(newState.Moves,
				fmt.Sprintf("Draw and remove King from deck: %s", formatCard(drawn)))
			res := solveState(newState, visited)
			if res.RemovedCount > bestResult.RemovedCount {
				bestResult = res
			}
			if bestResult.RemovedCount == 28 {
				return bestResult
			}
		} else {
			newState.Waste = append(newState.Waste, drawn)
			newState.Moves = append(newState.Moves,
				fmt.Sprintf("Draw card from deck: %s", formatCard(drawn)))
			res := solveState(newState, visited)
			if res.RemovedCount > bestResult.RemovedCount {
				bestResult = res
			}
			if bestResult.RemovedCount == 28 {
				return bestResult
			}
		}
	}

	// 4. If deck is empty but waste is not, recycle the waste.
	if len(state.Deck) == 0 && len(state.Waste) > 0 {
		newState := cloneState(state)
		// Recycle: new deck becomes the waste in reverse order.
		newDeck := make([]string, len(newState.Waste))
		for i, card := range newState.Waste {
			newDeck[len(newState.Waste)-1-i] = card
		}
		newState.Deck = newDeck
		newState.Waste = []string{}
		newState.Moves = append(newState.Moves, "Recycle waste into deck")
		res := solveState(newState, visited)
		if res.RemovedCount > bestResult.RemovedCount {
			bestResult = res
		}
		if bestResult.RemovedCount == 28 {
			return bestResult
		}
	}

	return bestResult
}

// solvePyramidSolitaire sets up the initial game state and returns the best move sequence found.
func solvePyramidSolitaire(cards []string) Result {
	pyramid := buildPyramid(cards)
	deck := make([]string, len(cards)-28)
	copy(deck, cards[28:])
	initialState := GameState{
		Pyramid: pyramid,
		Deck:    deck,
		Waste:   []string{},
		Moves:   []string{},
	}
	visited := make(map[string]int)
	return solveState(initialState, visited)
}

func main() {
	// The card sequence (first 28 cards form the pyramid, remaining form the deck).
	//cardSequence := "10c kh 3s 5d 10h ks 2d 8d ac qh 8s qd 4h 6s 8s 4s 9s 7c 10d 2s 7s 3d jh 6h 7d 3c 5c kd 10s ah 4c 6s 2c ad qc 6c 2h 4d jc jd as 8h 3h js 7h qs 9h 5h 5s 9d 9c kc"
	cardSequence := "jd 6h 4c 6c ac 3h 7c 2h jh 10s 8c ah qh 3d qd 2d 8s qc jc 4h 5s js 2s 3c 4d 7h 9c 5h 8h as 6d kd 5c kc 10d 8d 3s 9h ad kh 9d qs 7d 4s 9s 10h 10c ks 6s 5d 7s 2c"
	cards := strings.Split(cardSequence, " ")

	// Check that the deck is valid.
	if err := checkDeck(cards); err != nil {
		fmt.Printf("Deck check error: %s\n", err)
		os.Exit(1)
	}

	result := solvePyramidSolitaire(cards)

	fmt.Printf("Best partial solution removed %d of 28 pyramid cards.\n", result.RemovedCount)
	fmt.Println("Moves:")
	for _, move := range result.Moves {
		fmt.Println(move)
	}
}
