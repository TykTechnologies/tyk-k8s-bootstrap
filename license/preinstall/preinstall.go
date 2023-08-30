// Package preinstall exposes an API to run necessary operations required in pre-install hook job of bootstrapping.
// While bootstrapping Tyk Stack, users need to provide a valida Tyk License key.
// In the pre-hook installation, the helper functions defined in this package verifies the validity of the license.
package preinstall

import (
	"errors"
	"tyk/tyk/bootstrap/license"
)

// PreHookInstall runs all required License validity operations that are required in pre-install hook.
func PreHookInstall() error {
	dashboardLicenseKey, err := license.GetDashboardLicense()
	if err != nil {
		return err
	}

	licenseIsValid, err := license.ValidateDashboardLicense(dashboardLicenseKey)
	if err != nil {
		return err
	}

	if !licenseIsValid {
		return errors.New("provided license is invalid")
	}

	return nil
}
