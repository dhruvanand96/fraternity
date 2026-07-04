package auth

import "testing"

func TestCreateValidateAndRevoke(t *testing.T) {
	service := NewService("test-secret")
	token, err := service.Create(42)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	claims, err := service.Validate("Bearer " + token)
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if claims.UserID != 42 {
		t.Fatalf("Validate() user id = %d, want 42", claims.UserID)
	}

	service.Revoke(claims.SessionID)
	if _, err := service.Validate("Bearer " + token); err == nil {
		t.Fatal("Validate() accepted a revoked token")
	}
}
