package auth

//get a secret require ('crypto').randomBytes(64).toString('hex')
import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func InitJWTSecret() { //must be instillized at app start
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("[jwt.go] JWT_SECRET environment variable not set")
	}
	jwtSecret = []byte(secret)
}

type Claims struct {
	UserUUID string `json:"uuid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(userUUID string, userName string, ttl time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		UserUUID: userUUID,
		Username: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userUUID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func ParseJWT(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, funct(t *jwt.Token)(interface{}, error){
		//Enforce HMAC SHA-256
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil{
		return nil, err
	}
	 claims, ok := token.Claims.(*UserClaims)
	 if !ok || !token.Valid{
		return nil, jwt.ErrInvalidKey
	 }
	 return claims, nil
}