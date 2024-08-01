package auth

import (
    "testing"
)

func TestHashPassword(t *testing.T) {
    password := "mysecretpassword"
    hash, err := HashPassword(password)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if hash == "" {
        t.Fatalf("expected non-empty hash")
    }
}

func TestComparePasswordHash(t *testing.T) {
    password := "mysecretpassword"
    wrongPassword := "wrongpassword"

    hash, err := HashPassword(password)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    match, err := ComparePasswordHash(password, hash)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if !match {
        t.Fatalf("expected passwords to match")
    }

    match, err = ComparePasswordHash(wrongPassword, hash)
    if err == nil {
        t.Fatalf("expected error, got none")
    }
    if match {
        t.Fatalf("expected passwords not to match")
    }
}
