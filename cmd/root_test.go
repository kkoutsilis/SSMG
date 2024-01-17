package cmd

import (
	"html/template"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckFileExistsReturnsTrueWhenFileExists(t *testing.T) {
	if !checkFileExists("../test_data.json") {
		t.Error("Expected checkFileExists to return true when file exists")
	}

}

func TestCheckFileExistsReturnsFalseWhenFileDoesNotExist(t *testing.T) {
	if checkFileExists("file_does_not_exist.json") {
		t.Error("Expected checkFileExists to return false when file does not exist")
	}
}

func TestCheckIsJsonReturnsTrueWhenFileIsJson(t *testing.T) {
	if !checkIsJson("../test_data.json") {
		t.Error("Expected checkIsJson to return true when file is json")
	}
}

func TestCheckIsJsonReturnsFalseWhenFileIsNotJson(t *testing.T) {
	if checkIsJson("../test_data.txt") {
		t.Error("Expected checkIsJson to return false when file is not json")
	}
}

func TestGenerateSecretSantaMachesReturnsAppropriateLenghtOfMatchPairs(t *testing.T) {
	payload := make([]Data, 0, 3)
	payload = append(payload, Data{Name: "A", Email: "testa@test.org"})
	payload = append(payload, Data{Name: "B", Email: "testb@test.org"})
	payload = append(payload, Data{Name: "C", Email: "testc@test.org"})

	matches := generateSecretSantaMatches(payload)

	if len(matches) != len(payload) {
		t.Error("Expected generateSecretSantaMatches to return the same number of matches as participants")
	}
}

func TestCircularMatchingAlgorithmReturnsExpectedMatchPairs(t *testing.T) {
	payload := make([]Data, 0, 3)
	payload = append(payload, Data{Name: "A", Email: "testa@test.org"})
	payload = append(payload, Data{Name: "B", Email: "testb@test.org"})
	payload = append(payload, Data{Name: "C", Email: "testc@test.org"})

	expected := make([]MatchPair, 0, 3)

	expected = append(expected, MatchPair{From: payload[0], To: payload[1]})
	expected = append(expected, MatchPair{From: payload[1], To: payload[2]})
	expected = append(expected, MatchPair{From: payload[2], To: payload[0]})

	result := circlularMatchingAlgorithm(payload)

	if !reflect.DeepEqual(result, expected) {
		t.Error("Expected circlularMatchingAlgorithm to return expected list of match pairs")
	}

}

func TestPopulateEmailBodyReturnsExpectedString(t *testing.T) {
	p1 := Data{Name: "A", Email: "testa@test.org"}
	p2 := Data{Name: "B", Email: "testb@test.org"}

	matchPair := MatchPair{From: p1, To: p2}

	tmpl, err := template.New("test").Parse("Hi {{.From.Name}} your Secret Santa match is {{.To.Name}}!")
	if err != nil {
		t.Error("Expected template to be parsed without error")
	}

	emailBody, err := populateEmailBody(matchPair, tmpl)

	if err != nil {
		t.Error("Expected populateEmailBody to not return an error")
	}

	assert.Equal(t, "Hi A your Secret Santa match is B!", emailBody)

}

func TestPopulateEmailBodyReturnsErrorWhenInvalidTemplateProvided(t *testing.T) {
	p1 := Data{Name: "A", Email: "testa@test.org"}
	p2 := Data{Name: "B", Email: "testb@test.org"}

	matchPair := MatchPair{From: p1, To: p2}

	tmpl, err := template.New("test").Parse("Hi {{.Something.Unexpected}}!")
	if err != nil {
		t.Error("Expected template to be parsed without error")
	}

	emailBody, err := populateEmailBody(matchPair, tmpl)

	assert.Equal(t, "", emailBody)
	if err == nil {
		t.Error("Expected populateEmailBody to return an error when invalid template is provided")
	}

}

func TestGenerateEmailMessagesReturnsExpectedNumberOfMessages(t *testing.T) {
	payload := make([]Data, 0, 3)
	payload = append(payload, Data{Name: "A", Email: "testa@test.org"})
	payload = append(payload, Data{Name: "B", Email: "testb@test.org"})
	payload = append(payload, Data{Name: "C", Email: "testc@test.org"})

	matchPairs := make([]MatchPair, 0, 3)

	matchPairs = append(matchPairs, MatchPair{From: payload[0], To: payload[1]})
	matchPairs = append(matchPairs, MatchPair{From: payload[1], To: payload[2]})
	matchPairs = append(matchPairs, MatchPair{From: payload[2], To: payload[0]})

	tmpl, err := template.New("test").Parse("Hi {{.From.Name}} your Secret Santa match is {{.To.Name}}!")
	if err != nil {
		t.Error("Expected template to be parsed without error")
	}
	messages, err := generateEmailMessages(matchPairs, tmpl)

	if err != nil {
		t.Error("Expected generateEmailMessages to not return an error")
	}

	assert.Equal(t, len(matchPairs), len(messages))
}
