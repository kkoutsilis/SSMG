package cmd

import (
	"reflect"
	"testing"
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
	var payload []Data
	payload = append(payload, Data{Name: "A", Email: "testa@test.org"})
	payload = append(payload, Data{Name: "B", Email: "testb@test.org"})
	payload = append(payload, Data{Name: "C", Email: "testb@test.org"})

	matches := generateSecretSantaMatches(payload)

	if len(matches) != len(payload) {
		t.Error("Expected generateSecretSantaMatches to return the same number of matches as participants")
	}
}

func TestCircularMatchingAlgorithmReturnsExpectedMatchPairs(t *testing.T) {
	var payload []Data
	payload = append(payload, Data{Name: "A", Email: "testa@test.org"})
	payload = append(payload, Data{Name: "B", Email: "testb@test.org"})
	payload = append(payload, Data{Name: "C", Email: "testb@test.org"})

	var expected []MatchPair

	expected = append(expected, MatchPair{From: payload[0], To: payload[1]})
	expected = append(expected, MatchPair{From: payload[1], To: payload[2]})
	expected = append(expected, MatchPair{From: payload[2], To: payload[0]})

	result := circlularMatchingAlgorithm(payload)

	if !reflect.DeepEqual(result, expected) {
		t.Error("Expected circlularMatchingAlgorithm to return expected list of match pairs")
	}

}
