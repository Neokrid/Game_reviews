package service

import (
	"errors"
	"time"

	"github.com/Neokrid/game-review/pkg/model"
	"github.com/Neokrid/game-review/pkg/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	cost       = 14
	tokenTTL   = 12 * time.Hour
	signingKey = "qwetrq21324taqwf21"
)

type tokenClaims struct {
	jwt.StandardClaims
	Userid uuid.UUID `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}
func (s *AuthService) CreateUser(user model.User) (uuid.UUID, error) {
	user.PasswordHash = generatePasswordHash(user.PasswordHash)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", err
		}

		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})
	return token.SignedString([]byte(signingKey))
}

func generatePasswordHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		logrus.Fatalf("ошибка при хэширование пароля: %s", err.Error())
	}
	return string(hash)
}

func (s *AuthService) ParseToken(accesstoken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accesstoken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("недействительный метод подписи")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	calms, ok := token.Claims.(*tokenClaims)
	if !ok {
		return uuid.Nil, errors.New("заявки на токены не относятся к типу *tokenCtaims*")
	}
	return calms.Userid, nil
}
