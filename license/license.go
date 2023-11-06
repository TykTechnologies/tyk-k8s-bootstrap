package license

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func ValidateDashboardLicense(license string) (bool, error) {
	token, _ := jwt.Parse(license, func(token *jwt.Token) (interface{}, error) { // nolint:errcheck
		return []byte(""), nil
	})

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
