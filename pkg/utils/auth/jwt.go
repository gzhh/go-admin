package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
	jwt.RegisteredClaims
	User map[string]interface{} `json:"user,omitempty"`
}

func GenerateToken(secretKey string, userInfo map[string]interface{}) (string, error) {
	// Create claims with multiple fields populated
	claims := MyCustomClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "admin server",
			Subject:   "admin",
			ID:        userInfo["username"].(string),
		},
		userInfo,
	}
	// fmt.Printf("user: %+v\n", claims.User)

	mySigningKey := []byte(secretKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	// fmt.Printf("%v %v", ss, err)
	if err != nil {
		return "", err
	}
	return ss, nil
}
