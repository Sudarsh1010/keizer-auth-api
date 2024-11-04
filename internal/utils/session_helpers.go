package utils

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"keizer-auth-api/internal/models"
	"keizer-auth-api/internal/repositories"
	"time"

	"github.com/gofiber/fiber/v2"
)

const sessionExpiresIn = 30 * 24 * time.Hour

func GenerateSessionID() string {
	bytes := make([]byte, 15)
	rand.Read(bytes)
	return base32.StdEncoding.EncodeToString(bytes)
}

func ValidateSession(
	sessionID string,
	repo *repositories.SessionRepository,
) (*models.Session, error) {
	session, err := repo.GetSession(sessionID)
	if err != nil {
		return nil, errors.New("invalid session ID")
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, errors.New("expired session")
	}

	if time.Now().After(session.ExpiresAt.Add(-sessionExpiresIn / 2)) {
		session.ExpiresAt = time.Now().Add(sessionExpiresIn)
		if err := repo.UpdateSession(session); err != nil {
			return nil, err
		}
	}

	return session, nil
}

func SetSessionCookie(c *fiber.Ctx, sessionID string) {
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(sessionExpiresIn),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Path:     "/",
	})
}
