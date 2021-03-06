package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/atmiguel/cerealnotes/models"
	"github.com/dgrijalva/jwt-go"
)

var InvalidJWTokenError = errors.New("Token was invalid or unreadable")

func parseTokenFromString(tokenAsString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(
		strings.TrimSpace(tokenAsString),
		&JwtTokenClaim{},
		func(*jwt.Token) (interface{}, error) {
			return tokenSigningKey, nil
		})
}

func createTokenAsString(
	userId models.UserId,
	durationTilExpiration time.Duration,
) (string, error) {
	claims := JwtTokenClaim{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(durationTilExpiration).Unix(),
			Issuer:    "CerealNotes",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tokenSigningKey)
}

func getUserIdFromJwtToken(request *http.Request) (models.UserId, error) {
	cookie, err := request.Cookie(cerealNotesCookieName)
	if err != nil {
		return 0, err
	}

	token, err := parseTokenFromString(cookie.Value)
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*JwtTokenClaim); ok && token.Valid {
		return claims.UserId, nil
	}

	return 0, InvalidJWTokenError
}

func tokenTest1() {
	var num models.UserId = 32
	bob, err := createTokenAsString(num, 1)
	if err != nil {
		fmt.Println("create error")
		log.Fatal(err)
	}

	token, err := parseTokenFromString(bob)
	if err != nil {
		fmt.Println("parse error")
		log.Fatal(err)
	}
	fmt.Println(bob)
	if claims, ok := token.Claims.(*JwtTokenClaim); ok && token.Valid {
		if claims.UserId != 32 {
			log.Fatal("error in token")
		}
		fmt.Printf("%v %v", claims.UserId, claims.StandardClaims.ExpiresAt)
	} else {
		fmt.Println("Token claims could not be read")
	}
}
