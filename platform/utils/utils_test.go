package utils

import "testing"

func TestCalculateSessionState(t *testing.T) {
	clientID := "iAmClientID"
	origin := "www.clientWeb.com"
	opbs := GenerateNewOPBS()
	salt := GenerateRandomString(20, true)
	want := CalculateSessionState(clientID, origin, opbs, salt)
	got := CalculateSessionState(clientID, origin, opbs, salt)
	if want != got {
		t.Fatalf("got %s, want %s", got, want)
	}
}
