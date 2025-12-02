package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	accessTokenExpire  = 24 * time.Hour
	refreshTokenExpire = 72 * time.Hour
	signingKey         = "AIOJaqwdeqp321392sad"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID int    `json:"id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

func GenerateAccessToken(userID int, email, role string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenExpire).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userID,
		email,
		role,
	})

	signedToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		logrus.Errorf("Error generating access token for user with %d id", userID)
		return ""
	}

	return signedToken
}

func GenerateRefreshToken(userID int, email, role string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTokenExpire).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userID,
		email,
		role,
	})

	signedToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		logrus.Errorf("Error generating refresh token for user with %d id", userID)
		return ""
	}

	return signedToken
}

func ParseToken(accessToken string) (int, string, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing token method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, "", "", err
	}

	if !token.Valid {
		return 0, "", "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, "", "", err
	}

	return claims.UserID, claims.Email, claims.Role, nil
}
