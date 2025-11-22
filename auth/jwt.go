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
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(userUUID string, userName string, email string, ttl time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		UserUUID: userUUID,
		Username: userName,
		Email:    email,
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
	claims := &Claims{}
	//check that the algorithm in the token matches what we expect
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		//return the secret key that veriifies signature
		return jwtSecret, nil
	}
	//decode and verify the token
	token, err := jwt.ParseWithClaims(tokenStr, claims, keyFunc)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil

}
