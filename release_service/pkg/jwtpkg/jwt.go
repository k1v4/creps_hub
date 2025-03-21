package jwtpkg

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"release_service/internal/entity"
	"strings"
	"time"
)

const secret = "secret"

func ExtractToken(c echo.Context) string {
	bearerToken := c.Request().Header.Get("Authorization")

	if bearerToken == "" {
		return ""
	}

	return strings.TrimPrefix(bearerToken, "Bearer ")
}

func NewAccessToken(user entity.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken Функция для валидации токена
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("невалидный токен")
	}
}

// RefreshAccessToken Функция для обновления Access Token с помощью Refresh Token
func RefreshAccessToken(refreshToken string, duration time.Duration) (string, error) {
	// Валидируем Refresh Token
	claims, err := ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("невалидный Refresh Token: %v", err)
	}

	// Извлекаем данные пользователя из claims
	user := entity.User{
		ID:    int(claims["id"].(float64)), // JWT числа возвращает как float64
		Email: claims["email"].(string),
	}

	// Создаем новый Access Token
	newAccessToken, err := NewAccessToken(user, duration)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании Access Token: %v", err)
	}

	return newAccessToken, nil
}

func ValidateTokenAndGetUserId(tokenString string) (int, int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId, okUser := claims["id"].(float64)
		if !okUser {
			return 0, 0, fmt.Errorf("userId not found in token")
		}

		accessLevelId, okAccess := claims["access_level_id"].(float64)
		if !okAccess {
			return 0, 0, fmt.Errorf("access_level_id not found in token")
		}

		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return 0, 0, fmt.Errorf("token expired")
			}
		}

		return int(userId), int(accessLevelId), nil
	}

	return 0, 0, errors.New("invalid token")
}
