package authorization

import (
	"github.com/gbooksoft/backend/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func CreateTokens(userId, email string) (string, int64, string, int64, error) {
	atExp := time.Now().Add(2 * time.Hour).Unix()
	rtExp := time.Now().Add(7 * 24 * time.Hour).Unix()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    userId,
		"email": email,
		"exp":   atExp,
	})

	accessTokenStr, err := accessToken.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", 0, "", 0, err
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userId,
		"exp": rtExp,
	})
	refreshTokenStr, err := refreshToken.SignedString([]byte(os.Getenv("JWT_KEY")))

	if err != nil {
		return "", 0, "", 0, err
	}

	db := database.GetDBPool()

	_, err = db.Exec(`
INSERT INTO auth_token (expired_at, access, refresh) VALUES (?, ?, ?)`,
		atExp,
		accessTokenStr,
		refreshTokenStr,
	)

	if err != nil {
		return "", 0, "", 0, err
	}

	return accessTokenStr, atExp, refreshTokenStr, rtExp, err
}
