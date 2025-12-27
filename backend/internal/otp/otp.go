// Package otp provides OTP (One-Time Password) generation, storage, and verification.
package otp

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type OTP struct {
	ID         string
	Email      string
	OTPHash    string
	ExpiresAt  time.Time
	Attempts   int
	MaxAttempts int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// GenerateOTP generates a random 6-digit OTP
func GenerateOTP() (string, error) {
	// Generate a random number between 100000 and 999999
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	
	// Convert bytes to uint32 and ensure it's in the range [100000, 999999]
	randomNum := uint32(bytes[0])<<24 | uint32(bytes[1])<<16 | uint32(bytes[2])<<8 | uint32(bytes[3])
	otp := 100000 + int(randomNum%900000)
	return fmt.Sprintf("%06d", otp), nil
}

// HashOTP hashes an OTP using bcrypt
func HashOTP(otp string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash OTP: %w", err)
	}
	return string(hashed), nil
}

// CreateOTP creates a new OTP record for an email
// It invalidates any existing OTPs for that email first
func (s *Store) CreateOTP(email string) (string, error) {
	// Generate 6-digit OTP
	otp, err := GenerateOTP()
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Hash the OTP
	otpHash, err := HashOTP(otp)
	if err != nil {
		return "", fmt.Errorf("failed to hash OTP: %w", err)
	}

	// Expires in 5 minutes
	expiresAt := time.Now().Add(5 * time.Minute)

	// Invalidate any existing OTPs for this email
	_, err = s.db.Exec(
		"DELETE FROM email_otps WHERE email = $1",
		email,
	)
	if err != nil {
		log.Printf("[OTP] WARNING - Failed to delete existing OTPs: %v", err)
		// Continue anyway, as this is not critical
	}

	// Create new OTP record
	id := uuid.New().String()
	_, err = s.db.Exec(
		"INSERT INTO email_otps (id, email, otp_hash, expires_at, attempts, max_attempts, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		id, email, otpHash, expiresAt, 0, 3, time.Now(), time.Now(),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create OTP record: %w", err)
	}

	log.Printf("[OTP] Created OTP for email: %s (expires at: %s)", email, expiresAt.Format(time.RFC3339))
	return otp, nil
}

// VerifyOTP verifies an OTP for an email
// Returns true if valid, false otherwise
// Increments attempts counter
func (s *Store) VerifyOTP(email, otp string) (bool, error) {
	var otpRecord OTP
	err := s.db.QueryRow(
		"SELECT id, email, otp_hash, expires_at, attempts, max_attempts, created_at, updated_at FROM email_otps WHERE email = $1 ORDER BY created_at DESC LIMIT 1",
		email,
	).Scan(
		&otpRecord.ID, &otpRecord.Email, &otpRecord.OTPHash,
		&otpRecord.ExpiresAt, &otpRecord.Attempts, &otpRecord.MaxAttempts,
		&otpRecord.CreatedAt, &otpRecord.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to query OTP: %w", err)
	}

	// Check if expired
	if time.Now().After(otpRecord.ExpiresAt) {
		log.Printf("[OTP] OTP expired for email: %s", email)
		return false, nil
	}

	// Check if max attempts reached
	if otpRecord.Attempts >= otpRecord.MaxAttempts {
		log.Printf("[OTP] Max attempts reached for email: %s", email)
		return false, nil
	}

	// Increment attempts
	_, err = s.db.Exec(
		"UPDATE email_otps SET attempts = attempts + 1, updated_at = $1 WHERE id = $2",
		time.Now(), otpRecord.ID,
	)
	if err != nil {
		log.Printf("[OTP] WARNING - Failed to increment attempts: %v", err)
	}

	// Verify OTP hash
	err = bcrypt.CompareHashAndPassword([]byte(otpRecord.OTPHash), []byte(otp))
	if err != nil {
		log.Printf("[OTP] Invalid OTP for email: %s (attempts: %d/%d)", email, otpRecord.Attempts+1, otpRecord.MaxAttempts)
		return false, nil
	}

	// OTP is valid - delete it to prevent reuse
	_, err = s.db.Exec("DELETE FROM email_otps WHERE id = $1", otpRecord.ID)
	if err != nil {
		log.Printf("[OTP] WARNING - Failed to delete OTP after verification: %v", err)
	}

	log.Printf("[OTP] OTP verified successfully for email: %s", email)
	return true, nil
}

// CleanupExpiredOTPs removes expired OTPs from the database
func (s *Store) CleanupExpiredOTPs() error {
	result, err := s.db.Exec("DELETE FROM email_otps WHERE expires_at < NOW()")
	if err != nil {
		return fmt.Errorf("failed to cleanup expired OTPs: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("[OTP] Cleaned up %d expired OTP(s)", rowsAffected)
	}
	return nil
}

