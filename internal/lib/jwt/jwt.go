package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserProvider interface {
	GetID() string
	GetEmail() string
}

type AppProvider interface {
	GetID() string
	GetSecret() []byte
}

func NewToken(
	user UserProvider,
	app AppProvider,
	expiration time.Duration,
) string {

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userID": user.GetID(),
			"email":  user.GetEmail(),
			"exp":    time.Now().Add(expiration).Unix(),
			"appID":  app.GetID(),
		})

	// No need to handle error - it returns only if @param key is not []byte which app.GetSecret() guarantees
	tokeString, _ := token.SignedString(app.GetSecret())

	return tokeString
}
