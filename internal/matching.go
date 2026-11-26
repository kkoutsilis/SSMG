package internal

import "math/rand"

type Match struct {
	From Participant
	To   Participant
}

func GenerateMatches(participants []Participant) []Match {
	shuffled := append([]Participant(nil), participants...)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return CircularMatching(shuffled)
}

func CircularMatching(participants []Participant) []Match {
	matches := make([]Match, 0, len(participants))
	for i, from := range participants {
		to := participants[(i+1)%len(participants)]
		matches = append(matches, Match{From: from, To: to})
	}
	return matches
}
