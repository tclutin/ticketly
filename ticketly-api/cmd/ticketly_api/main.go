package main

import (
	"errors"
	"fmt"
	"github.com/centrifugal/centrifuge-go"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tclutin/ticketly/ticketly_api/internal/app"
	"log"
	"strconv"
	"time"
)

var (
	ErrSigningMethod      = errors.New("invalid signing method")
	ErrTokenExpired       = errors.New("token is expired")
	ErrInvalidClaimFormat = errors.New("invalid claim format")
	ErrSubClaimNotFound   = errors.New("sub claim not found")
	ErrInvalidSubFormat   = errors.New("invalid sub format")
)

type Manager interface {
	ParseToken(accessToken string) (string, error)
	NewAccessToken(userID uint64, ttl time.Duration) (string, error)
	NewRefreshToken() uuid.UUID
}

type TokenManager struct {
	signingKey string
}

func MustLoadTokenManager(secret string) Manager {
	if secret == "" {
		log.Fatalln("signingKey is empty")
	}
	return &TokenManager{signingKey: secret}
}

func (t TokenManager) ParseToken(jwtToken string) (string, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrSigningMethod
		}

		return []byte(t.signingKey), nil
	})

	if err != nil {
		return "", ErrTokenExpired
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidClaimFormat
	}

	sub, ok := claims["sub"]
	if !ok {
		return "", ErrSubClaimNotFound
	}

	subFloat, ok := sub.(string)
	if !ok {
		return "", ErrInvalidSubFormat
	}

	return subFloat, nil
}

func (t TokenManager) NewAccessToken(userID uint64, ttl time.Duration) (string, error) {
	claim := jwt.MapClaims{
		"exp": time.Now().UTC().Add(ttl).Unix(),
		"sub": strconv.FormatUint(userID, 10),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString([]byte(t.signingKey))
}

func (t TokenManager) NewRefreshToken() uuid.UUID {
	return uuid.New()
}

func main() {

	app.New().Run()
	test1 := MustLoadTokenManager("app.go")

	test, _ := test1.NewAccessToken(12, 1*time.Hour)

	fmt.Println(test)

	client := centrifuge.NewJsonClient("ws://localhost:8000/connection/websocket", centrifuge.Config{
		Token: test,
	})

	defer client.Close()
	// 4. Обработчики событий
	client.OnConnected(func(_ centrifuge.ConnectedEvent) {
		log.Println("Connected to Centrifugo")
	})

	client.OnDisconnected(func(_ centrifuge.DisconnectedEvent) {
		log.Println("Disconnected from Centrifugo")
	})

	// 5. Подключение
	if err := client.Connect(); err != nil {
		log.Fatalf("Connection error: %v", err)
	}

	// 6. Подписка на канал
	subscription, err := client.NewSubscription("general")
	if err != nil {
		log.Fatalf("Subscription error: %v", err)
	}

	subscription.OnPublication(func(e centrifuge.PublicationEvent) {
		log.Printf("New message: %s", string(e.Data))
	})

	if err := subscription.Subscribe(); err != nil {
		log.Fatalf("Subscribe error: %v", err)
	}

	for {

	}
}
