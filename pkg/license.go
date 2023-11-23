package pkg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

// ValidateDashboardLicense validates if the given license argument is a valid and not expired
// JWT token.
func ValidateDashboardLicense(license string) (bool, error) {
	if license == "" {
		return false, fmt.Errorf("empty license")
	}

	token, _ := jwt.Parse(license, func(token *jwt.Token) (interface{}, error) { // nolint:errcheck
		return []byte(""), nil
	})

	if token == nil {
		return false, fmt.Errorf("failed to parse license %v\n", license)
	}

	if strings.ToLower(fmt.Sprint(token.Header["typ"])) == "jwt" {
		exp := strings.Split(fmt.Sprintf("%f", token.Claims.(jwt.MapClaims)["exp"]), ".")[0]

		expDate, err := strconv.ParseInt(exp, 10, 64)
		if err != nil {
			return false, errors.New("impossible to parse expiration date")
		}

		if time.Unix(expDate, 0).Before(time.Now()) {
			return false, errors.New("expired dashboard license")
		}
	} else {
		return false, errors.New("token is not of jwt type")
	}

	return true, nil
}
