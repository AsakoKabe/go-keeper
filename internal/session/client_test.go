package session

import "testing"

func TestNewClientSession(t *testing.T) {
	session := NewClientSession()
	if session == nil {
		t.Fatalf("expected new session, got nil")
	}

	if session.Token != "" {
		t.Fatalf("expected empty token, got '%s'", session.Token)
	}
}

func TestSetToken(t *testing.T) {
	session := NewClientSession()
	token := "test-token"

	session.SetToken(token)

	if session.Token != token {
		t.Fatalf("expected token to be '%s', got '%s'", token, session.Token)
	}
}

func TestIsAuth(t *testing.T) {
	session := NewClientSession()

	if session.IsAuth() {
		t.Fatalf("expected IsAuth to return false for empty token, got true")
	}

	session.SetToken("valid-token")

	if !session.IsAuth() {
		t.Fatalf("expected IsAuth to return true for non-empty token, got false")
	}
}
