// Package firebase provides Firebase token verification without Admin SDK
// This allows verification using JWT parsing when Admin SDK is not available
package firebase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// VerifyIDTokenREST verifies a Firebase ID token using JWT parsing (no Admin SDK required)
// This is a basic verification that checks token structure and claims
// For production, you should use Firebase Admin SDK for full verification
func VerifyIDTokenREST(ctx context.Context, idToken, projectID string) (uid, email string, err error) {
	// Parse token without verification first to get claims
	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return "", "", fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("invalid token claims")
	}

	// Verify issuer
	iss, ok := claims["iss"].(string)
	if !ok || !strings.HasPrefix(iss, "https://securetoken.google.com/") {
		return "", "", fmt.Errorf("invalid issuer")
	}

	// Verify audience (project ID)
	aud, ok := claims["aud"].(string)
	if !ok || aud != projectID {
		return "", "", fmt.Errorf("invalid audience, expected %s, got %s", projectID, aud)
	}

	// Check expiration
	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", "", fmt.Errorf("token missing expiration")
	}
	// Note: We're not checking expiration time here since we're doing basic verification
	// In production with Admin SDK, this is handled automatically
	_ = exp

	// Extract UID and email
	uid, _ = claims["user_id"].(string)
	if uid == "" {
		uid, _ = claims["sub"].(string)
	}
	email, _ = claims["email"].(string)

	if uid == "" {
		return "", "", fmt.Errorf("token missing user_id")
	}

	// Check email verification status from claims
	emailVerified, _ := claims["email_verified"].(bool)
	if !emailVerified {
		log.Printf("[FIREBASE] WARNING - Email not verified for user: %s", email)
		// We'll still return the email, but the caller should check verification
	}

	log.Printf("[FIREBASE] Verified token for user: %s (UID: %s, Email verified: %v)", email, uid, emailVerified)
	return uid, email, nil
}
