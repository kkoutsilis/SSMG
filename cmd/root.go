package cmd

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var version string = "0.0.1"

type Data struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	WishList []string `json:"wishlist omitempty"`
}

type MatchPair struct {
	Person1 string `json:"person1"`
	Person2 string `json:"person2"`
}

func checkFileExists(path string) bool {
	_, error := os.Stat(path)
	return !errors.Is(error, os.ErrNotExist)
}

func checkIsJson(path string) bool {
	fileExtension := filepath.Ext(path)
	return fileExtension == ".json"
}

func match(data []Data) []MatchPair {
	var matches []MatchPair
	rand.Seed(time.Now().UnixNano())

	shuffledData := make([]Data, len(data))
	copy(shuffledData, data)

	for i := range shuffledData {
		j := rand.Intn(i + 1)
		shuffledData[i], shuffledData[j] = shuffledData[j], shuffledData[i]
	}

	for i := 0; i < len(shuffledData); i++ {
		matches = append(matches, MatchPair{Person1: shuffledData[i].Name, Person2: shuffledData[(i+1)%len(shuffledData)].Name})
		// remove the matched pair from the list
		shuffledData = append(shuffledData[:i], shuffledData[i+1:]...)
	}

	return matches
}

var rootCmd = &cobra.Command{
	Use:        "run [path]",
	Short:      "A cli tool that generates secret santa matches",
	ArgAliases: []string{"path"},
	Version:    version,
	Run: func(cmd *cobra.Command, args []string) {
		var filePath string = args[0]
		if filePath == "" {
			filePath = "data.json"
		}

		if !checkFileExists(filePath) {
			log.Fatal("File ", filePath, " does not exist")
		}

		if !checkIsJson(filePath) {
			log.Fatal("File ", filePath, " is not a json file")
		}

		file, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal("Error when opening file: ", err)
		}

		var payload []Data
		err = json.Unmarshal(file, &payload)
		if err != nil {
			log.Fatal("Error reading file content: ", err)
		}

		if len(payload) == 0 {
			log.Fatal("File is empty")
		}

		if len(payload)%2 != 0 {
			log.Fatal("Number of participants is odd (", len(payload), "), please provide an even number of participants for the secret santa matches")
		}

		matches := match(payload)

		for _, match := range matches {
			log.Println(match.Person1, " -> ", match.Person2)
		}

		// TODO send email to each participant with their match
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
