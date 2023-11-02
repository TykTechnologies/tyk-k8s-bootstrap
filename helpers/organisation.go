package helpers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"tyk/tyk/bootstrap/data"

	"k8s.io/apimachinery/pkg/util/json"
)

const (
	AdminOrganisationsEndpoint     = "/admin/organisations"
	ApiUsersActionsResetEndpoint   = "%s/api/users/%s/actions/reset"
	ApiPortalCatalogueEndpoint     = "/api/portal/catalogue"
	ApiPortalPagesEndpoint         = "/api/portal/pages"
	ApiPortalConfigurationEndpoint = "/api/portal/configuration"
	ApiPortalCnameEndpoint         = "/api/portal/cname"

	TykModePro = "pro"
	TykAuth    = "TYK_AUTH"
	TykOrg     = "TYK_ORG"
	TykMode    = "TYK_MODE"
	TykUrl     = "TYK_URL"
)

func CheckForExistingOrganisation(client http.Client) error {
	fmt.Println("Checking for existing organisations")

	orgsApiEndpoint := data.BootstrapConf.K8s.DashboardSvcUrl + AdminOrganisationsEndpoint
	req, err := http.NewRequest("GET", orgsApiEndpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set("admin-auth", data.BootstrapConf.Tyk.Admin.Secret)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	orgs := OrgResponse{}
	err = json.Unmarshal(bodyBytes, &orgs)
	if err != nil {
		return err
	}
	if len(orgs.Organisations) > 0 {
		for _, organisation := range orgs.Organisations {
			if organisation["owner_name"] == data.BootstrapConf.Tyk.Org.Name ||
				organisation["cname"] == data.BootstrapConf.Tyk.Org.Cname {
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

type CreateOrgStruct struct {
	OwnerName    string `json:"owner_name"`
	CnameEnabled bool   `json:"cname_enabled"`
	Cname        string `json:"cname"`
}

func CreateOrganisation(client http.Client, dashBoardUrl string) (string, error) {
	createOrgData := CreateOrgStruct{
		OwnerName:    data.BootstrapConf.Tyk.Org.Name,
		CnameEnabled: true,
		Cname:        data.BootstrapConf.Tyk.Org.Cname,
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
	req.Header.Set("admin-auth", data.BootstrapConf.Tyk.Admin.Secret)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	createOrgResponse := DashboardGeneralResponse{}
	err = json.Unmarshal(bodyBytes, &createOrgResponse)
	if err != nil {
		return "", err
	}

	return createOrgResponse.Meta, nil
}

type DashboardGeneralResponse struct {
	Status  string `json:"Status"`
	Message string `json:"Message"`
	Meta    string `json:"Meta"`
}

type OrgResponse struct {
	Organisations []map[string]interface{} `json:"organisations"`
	Pages         int                      `json:"pages"`
}
