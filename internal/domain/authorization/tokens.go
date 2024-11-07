package authorization

import (
	"github.com/gbooksoft/backend/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type CustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func CreateTokens(userId, email string) (string, string, error) {
	issuedAt := time.Now()
	atExp := issuedAt.Add(2 * time.Hour)
	rtExp := issuedAt.Add(7 * 24 * time.Hour)

	claims := CustomClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			Audience:  jwt.ClaimStrings{"web-front"},
			ExpiresAt: jwt.NewNumericDate(time.Unix(atExp.Unix(), 0)),
			IssuedAt:  jwt.NewNumericDate(time.Unix(issuedAt.Unix(), 0)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessTokenStr, err := accessToken.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", "", err
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userId,
		ExpiresAt: jwt.NewNumericDate(time.Unix(rtExp.Unix(), 0)),
	})
	refreshTokenStr, err := refreshToken.SignedString([]byte(os.Getenv("JWT_KEY")))

	if err != nil {
		return "", "", err
	}

	db := database.GetDBPool()

	_, err = db.Exec(`
INSERT INTO auth_token (expired_at, access, refresh) VALUES (?, ?, ?)`,
		atExp.Unix(),
		accessTokenStr,
		refreshTokenStr,
	)

	if err != nil {
		return "", "", err
	}

	return accessTokenStr, refreshTokenStr, err
}
