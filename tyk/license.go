package tyk

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
	"strconv"
	"strings"
	"time"
	"tyk/tyk/bootstrap/data"
)

func readDashboardLicense() (string, error) {
	license, ok := os.LookupEnv(data.TykDashboardLicenseEnvVarName)
	if !ok {
		return "", errors.New("license env var is not present")
	}

	if license == "" {
		return "", errors.New("empty dashboard license")
	}

	return license, nil
}

func ValidateDashboardLicense() (bool, error) {
	license, err := readDashboardLicense()
	if err != nil {
		return false, err
	}

	token, _ := jwt.Parse(license, func(token *jwt.Token) (interface{}, error) {
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
