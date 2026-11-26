package internal

import (
	"html/template"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadParticipants(t *testing.T) {
	path := filepath.Join(t.TempDir(), "participants.json")
	if err := os.WriteFile(path, []byte(`[{"name":"A","email":"a@example.com","wishlist":["book"]}]`), 0o600); err != nil {
		t.Fatal(err)
	}

	participants, err := LoadParticipants(path)
	if err != nil {
		t.Fatal(err)
	}
	want := []Participant{{Name: "A", Email: "a@example.com", WishList: []string{"book"}}}
	if !reflect.DeepEqual(participants, want) {
		t.Fatalf("LoadParticipants() = %#v, want %#v", participants, want)
	}
}

func TestRenderEmailBody(t *testing.T) {
	tmpl, err := template.New("email").Parse("Hi {{.From.Name}}, your match is {{.To.Name}}")
	if err != nil {
		t.Fatal(err)
	}
	match := Match{From: Participant{Name: "A"}, To: Participant{Name: "B"}}
	body, err := renderEmailBody(match, tmpl)
	if err != nil {
		t.Fatal(err)
	}
	if want := "Hi A, your match is B"; body != want {
		t.Fatalf("renderEmailBody() = %q, want %q", body, want)
	}
}

func TestLoadParticipantsRejectsEmptyInput(t *testing.T) {
	path := filepath.Join(t.TempDir(), "participants.json")
	if err := os.WriteFile(path, []byte("[]"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadParticipants(path); err == nil {
		t.Fatal("LoadParticipants() expected an error for empty input")
	}
}

func TestCircularMatching(t *testing.T) {
	participants := []Participant{{Name: "A"}, {Name: "B"}, {Name: "C"}}
	want := []Match{
		{From: participants[0], To: participants[1]},
		{From: participants[1], To: participants[2]},
		{From: participants[2], To: participants[0]},
	}
	if got := CircularMatching(participants); !reflect.DeepEqual(got, want) {
		t.Fatalf("CircularMatching() = %#v, want %#v", got, want)
	}
}

func TestGenerateMatchesPreservesParticipants(t *testing.T) {
	participants := []Participant{{Name: "A"}, {Name: "B"}, {Name: "C"}}
	matches := GenerateMatches(participants)
	if len(matches) != len(participants) {
		t.Fatalf("GenerateMatches() returned %d matches, want %d", len(matches), len(participants))
	}
	if !reflect.DeepEqual(participants, []Participant{{Name: "A"}, {Name: "B"}, {Name: "C"}}) {
		t.Fatal("GenerateMatches() modified the input participants")
	}
}
