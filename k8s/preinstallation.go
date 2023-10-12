package k8s

import (
	"errors"
	"tyk/tyk/bootstrap/tyk"
)

// PreHookInstall runs all required license validation operations that are required in pre-install hook.
func PreHookInstall() error {
	licenseIsValid, err := tyk.ValidateDashboardLicense()
	if err != nil {
		return err
	}

	if !licenseIsValid {
		return errors.New("provided license is invalid")
	}

	return nil
}
