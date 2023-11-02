// Package preinstallation exposes an API to run necessary operations required in pre-install hook job of bootstrapping.
// While bootstrapping Tyk Stack, users need to provide a valida Tyk License key.
// In the pre-hook installation, the helper functions defined in this package verifies the validity of the license.
package preinstallation

import (
	"errors"
	"tyk/tyk/bootstrap/data"
	"tyk/tyk/bootstrap/license"
)

// PreHookInstall runs all required license validation operations that are required in pre-install hook.
func PreHookInstall() error {
	licenseIsValid, err := license.ValidateDashboardLicense(data.BootstrapConf.Tyk.DashboardLicense)
	if err != nil {
		return err
	}

	if !licenseIsValid {
		return errors.New("provided license is invalid")
	}

	return nil
}
