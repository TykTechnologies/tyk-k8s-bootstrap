package tyk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"tyk/tyk/bootstrap/tyk/api"
)

const (
	AdminOrganisationsEndpoint     = "/admin/organisations"
	ApiUsersActionsResetEndpoint   = "%s/api/users/%s/actions/reset"
	ApiPortalCatalogueEndpoint     = "/api/portal/catalogue"
	ApiPortalPagesEndpoint         = "/api/portal/pages"
	ApiPortalConfigurationEndpoint = "/api/portal/configuration"
	ApiPortalCnameEndpoint         = "/api/portal/cname"

	TykModePro = "pro"

	TykAuth = "TYK_AUTH"
	TykOrg  = "TYK_ORG"
	TykMode = "TYK_MODE"
	TykUrl  = "TYK_URL"
)

func (s *Service) CheckForExistingOrganisation() error {
	fmt.Println("Checking for existing organisations")

	orgsApiEndpoint := s.appArgs.DashboardUrl + AdminOrganisationsEndpoint
	req, err := http.NewRequest("GET", orgsApiEndpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set("admin-auth", s.appArgs.TykAdminSecret)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	orgs := api.OrgResponse{}
	err = json.Unmarshal(bodyBytes, &orgs)
	if err != nil {
		return err
	}

	if len(orgs.Organisations) > 0 {
		for _, organisation := range orgs.Organisations {
			if organisation["owner_name"] == s.appArgs.CurrentOrgName ||
				organisation["cname"] == s.appArgs.Cname {
				return errors.New("there shouldn't be any organisations, please " +
					"disable bootstrapping to avoid losing data or delete " +
					"already existing organisations")
			}
		}
	} else {
		fmt.Println("No organisations have been detected, we can proceed")
		return nil
	}

	return nil
}

func (s *Service) createOrganisation(dashBoardUrl string) (string, error) {
	createOrgData := api.CreateOrgRequest{
		OwnerName:    s.appArgs.CurrentOrgName,
		CnameEnabled: true,
		Cname:        s.appArgs.Cname,
	}

	reqBodyBytes, err := json.Marshal(createOrgData)
	if err != nil {
		return "", err
	}

	reqBody := bytes.NewReader(reqBodyBytes)
	req, err := http.NewRequest("POST", dashBoardUrl+AdminOrganisationsEndpoint, reqBody)
	if err != nil {
		return "", err
	}

	req.Header.Set("admin-auth", s.appArgs.TykAdminSecret)
	req.Header.Set("Content-Type", "application/json")

	res, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	createOrgResponse := api.DashboardGeneralResponse{}
	err = json.Unmarshal(bodyBytes, &createOrgResponse)
	if err != nil {
		return "", err
	}

	return createOrgResponse.Meta, nil
}
