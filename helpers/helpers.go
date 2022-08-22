package helpers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/http"
	"os"
	"time"
	"tyk/tyk/bootstrap/constants"
	"tyk/tyk/bootstrap/data"
)

const AdminOrganisationsEndpoint = "/admin/organisations"
const ApiUsersActionsResetEndpoint = "%s/api/users/%s/actions/reset"
const ApiPortalCatalogueEndpoint = "/api/portal/catalogue"
const ApiPortalPagesEndpoint = "/api/portal/pages"
const ApiPortalConfigurationEndpoint = "/api/portal/configuration"

const TykModePro = "pro"

const TykAuth = "TYK_AUTH"
const TykOrg = "TYK_ORG"
const TykMode = "TYK_MODE"
const TykUrl = "TYK_URL"

func CheckForExistingOrganisation(client http.Client) error {
	fmt.Println("Checking for existing organisations")

	orgsApiEndpoint := data.AppConfig.DashboardUrl + AdminOrganisationsEndpoint
	req, err := http.NewRequest("GET", orgsApiEndpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("admin-auth", data.AppConfig.TykAdminSecret)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

	orgs := OrgResponse{}
	err = json.Unmarshal(bodyBytes, &orgs)
	if err != nil {
		return err
	}
	if len(orgs.Organisations) > 0 {
		for _, organisation := range orgs.Organisations {
			if organisation["owner_name"] == data.AppConfig.CurrentOrgName ||
				organisation["cname"] == data.AppConfig.Cname {
				return errors.New("there shouldn't be any organisations, please " +
					"disable bootstrapping to avoid losing data or delete" +
					"already existing organisations")
			}
		}
	} else {
		fmt.Println("No organisations have been detected, we can proceed")
		return nil
	}
	return nil
}

func CreateUser(client http.Client, dashboardUrl string, orgId string) (string, error) {
	userData, err := GetUserData(client, dashboardUrl, orgId)

	if err != nil {
		return "", err
	}

	err = SetUserPassword(client, userData.UserId, userData.AuthCode, dashboardUrl)
	if err != nil {
		return "", err
	}
	return userData.AuthCode, nil
}

type ResetPasswordStruct struct {
	NewPassword     string            `json:"new_password"`
	UserPermissions map[string]string `json:"user_permissions"`
}

func SetUserPassword(client http.Client, userId string, authCode string, dashboardUrl string) error {
	newPasswordData := ResetPasswordStruct{
		NewPassword:     data.AppConfig.TykAdminPassword,
		UserPermissions: map[string]string{"IsAdmin": "admin"},
	}
	reqBody, err := json.Marshal(newPasswordData)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf(ApiUsersActionsResetEndpoint, dashboardUrl, userId),
		bytes.NewReader(reqBody))
	req.Header.Set("authorization", authCode)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("resetting password did not work")
	}
	return nil
}

func GenerateCredentials(client http.Client) error {
	orgId, err := CreateOrganisation(client, data.AppConfig.DashboardUrl)
	if err != nil {
		return err
	}

	data.AppConfig.OrgId = orgId

	userAuth, err := CreateUser(client, data.AppConfig.DashboardUrl, orgId)
	if err != nil {
		return err
	}

	data.AppConfig.UserAuth = userAuth

	return nil
}

func BoostrapPortal(client http.Client) error {
	err := CreatePortalDefaultSettings(client)
	if err != nil {
		return err
	}

	err = InitialiseCatalogue(client)
	if err != nil {
		return err
	}

	err = CreatePortalHomepage(client)
	if err != nil {
		return err
	}

	fmt.Println("finished bootstrapping portal")

	return nil
}

type InitCatalogReq struct {
	OrgId string `json:"org_id"`
}

func InitialiseCatalogue(client http.Client) error {
	fmt.Println("Initialising Catalogue")

	initCatalog := InitCatalogReq{OrgId: data.AppConfig.OrgId}
	reqBody, err := json.Marshal(initCatalog)
	if err != nil {
		return err
	}
	reqData := bytes.NewReader(reqBody)
	req, err := http.NewRequest("POST", data.AppConfig.DashboardUrl+ApiPortalCatalogueEndpoint, reqData)
	req.Header.Set("Authorization", data.AppConfig.UserAuth)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return err
	}
	resp := DashboardGeneralResponse{}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(bodyBytes, &resp)
	if err != nil {
		return err
	}
	data.AppConfig.CatalogId = resp.Message

	fmt.Println(string(bodyBytes))
	return nil
}

func CreatePortalHomepage(client http.Client) error {
	fmt.Println("Creating portal homepage")

	homepageContents := GetPortalHomepage()
	reqBody, err := json.Marshal(homepageContents)
	if err != nil {
		return err
	}
	reqData := bytes.NewReader(reqBody)
	req, err := http.NewRequest("POST", data.AppConfig.DashboardUrl+ApiPortalPagesEndpoint, reqData)
	req.Header.Set("Authorization", data.AppConfig.UserAuth)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return err
	}
	resp := DashboardGeneralResponse{}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(bodyBytes, &resp)
	if err != nil {
		return err
	}

	fmt.Println(string(bodyBytes))

	return nil
}

func GetPortalHomepage() PortalHomepageRequest {
	return PortalHomepageRequest{
		IsHomepage:   true,
		TemplateName: "",
		Title:        "Developer portal name",
		Slug:         "/",
		Fields: PortalFields{
			JumboCTATitle:       "Tyk Developer Portal",
			SubHeading:          "Sub Header",
			JumboCTALink:        "#cta",
			JumboCTALinkTitle:   "Your awesome APIs, hosted with Tyk!",
			PanelOneContent:     "Panel 1 content.",
			PanelOneLink:        "#panel1",
			PanelOneLinkTitle:   "Panel 1 Button",
			PanelOneTitle:       "Panel 1 Title",
			PanelThereeContent:  "",
			PanelThreeContent:   "Panel 3 content.",
			PanelThreeLink:      "#panel3",
			PanelThreeLinkTitle: "Panel 3 Button",
			PanelThreeTitle:     "Panel 3 Title",
			PanelTwoContent:     "Panel 2 content.",
			PanelTwoLink:        "#panel2",
			PanelTwoLinkTitle:   "Panel 2 Button",
			PanelTwoTitle:       "Panel 2 Title",
		},
	}

}

type PortalHomepageRequest struct {
	IsHomepage   bool         `json:"is_homepage"`
	TemplateName string       `json:"template_name"`
	Title        string       `json:"title"`
	Slug         string       `json:"slug"`
	Fields       PortalFields `json:"fields"`
}

type PortalFields struct {
	JumboCTATitle       string `json:"JumboCTATitle"`
	SubHeading          string `json:"SubHeading"`
	JumboCTALink        string `json:"JumboCTALink"`
	JumboCTALinkTitle   string `json:"JumboCTALinkTitle"`
	PanelOneContent     string `json:"PanelOneContent"`
	PanelOneLink        string `json:"PanelOneLink"`
	PanelOneLinkTitle   string `json:"PanelOneLinkTitle"`
	PanelOneTitle       string `json:"PanelOneTitle"`
	PanelThereeContent  string `json:"PanelThereeContent"`
	PanelThreeContent   string `json:"PanelThreeContent"`
	PanelThreeLink      string `json:"PanelThreeLink"`
	PanelThreeLinkTitle string `json:"PanelThreeLinkTitle"`
	PanelThreeTitle     string `json:"PanelThreeTitle"`
	PanelTwoContent     string `json:"PanelTwoContent"`
	PanelTwoLink        string `json:"PanelTwoLink"`
	PanelTwoLinkTitle   string `json:"PanelTwoLinkTitle"`
	PanelTwoTitle       string `json:"PanelTwoTitle"`
}

func CreatePortalDefaultSettings(client http.Client) error {
	fmt.Println("Creating bootstrap default settings")

	req, err := http.NewRequest("POST", data.AppConfig.DashboardUrl+ApiPortalConfigurationEndpoint, nil)
	req.Header.Set("Authorization", data.AppConfig.UserAuth)

	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(resBytes))

	return nil
}

type CreateUserRequest struct {
	OrganisationId  string            `json:"org_id"`
	FirstName       string            `json:"first_name"`
	LastName        string            `json:"last_name"`
	EmailAddress    string            `json:"email_address"`
	Active          bool              `json:"active"`
	UserPermissions map[string]string `json:"user_permissions"`
}

type CreateUserResponse struct {
	Status  string `json:"Status"`
	Message string `json:"Message"`
	Meta    struct {
		APIModel struct {
		} `json:"api_model"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		EmailAddress    string `json:"email_address"`
		OrgID           string `json:"org_id"`
		Active          bool   `json:"active"`
		ID              string `json:"id"`
		AccessKey       string `json:"access_key"`
		UserPermissions struct {
			IsAdmin string `json:"IsAdmin"`
		} `json:"user_permissions"`
		GroupID         string        `json:"group_id"`
		PasswordMaxDays int           `json:"password_max_days"`
		PasswordUpdated time.Time     `json:"password_updated"`
		PWHistory       []interface{} `json:"PWHistory"`
		CreatedAt       time.Time     `json:"created_at"`
	} `json:"Meta"`
}

type NeededUserData struct {
	AuthCode string
	UserId   string
}

func GetUserData(client http.Client, dashboardUrl string, orgId string) (NeededUserData, error) {
	reqBody := CreateUserRequest{
		OrganisationId:  orgId,
		FirstName:       data.AppConfig.TykAdminFirstName,
		LastName:        data.AppConfig.TykAdminLastName,
		EmailAddress:    data.AppConfig.TykAdminEmailAddress,
		Active:          true,
		UserPermissions: map[string]string{"IsAdmin": "admin"},
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return NeededUserData{}, err
	}
	reqData := bytes.NewReader(reqBytes)
	req, err := http.NewRequest("POST", dashboardUrl+"/admin/users", reqData)

	req.Header.Set("admin-auth", data.AppConfig.TykAdminSecret)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return NeededUserData{}, err
	}

	fmt.Println(res.StatusCode)
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

	getUserResponse := CreateUserResponse{}
	err = json.Unmarshal(bodyBytes, &getUserResponse)
	if err != nil {
		return NeededUserData{}, err
	}
	return NeededUserData{UserId: getUserResponse.Meta.ID, AuthCode: getUserResponse.Message}, nil
}

type CreateOrgStruct struct {
	OwnerName    string `json:"owner_name"`
	CnameEnabled bool   `json:"cname_enabled"`
	Cname        string `json:"cname"`
}

func CreateOrganisation(client http.Client, dashBoardUrl string) (string, error) {
	createOrgData := CreateOrgStruct{
		OwnerName:    data.AppConfig.CurrentOrgName,
		CnameEnabled: true,
		Cname:        data.AppConfig.Cname,
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
	req.Header.Set("admin-auth", data.AppConfig.TykAdminSecret)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

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

func BootstrapTykOperatorSecret() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	secrets, err := clientset.CoreV1().Secrets(os.Getenv(constants.TykPodNamespaceEnvVarName)).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	found := false
	for _, value := range secrets.Items {
		if value.Name == os.Getenv(constants.OperatorSecretNameEnvVarName) {
			err = clientset.CoreV1().Secrets(os.Getenv(constants.TykPodNamespaceEnvVarName)).
				Delete(context.TODO(), value.Name, metav1.DeleteOptions{})
			if err != nil {
				return err
			}
			found = true
			break
		}
	}

	if found == false {
		fmt.Println("A previously created operator secret has not been identified")
		err = CreateTykOperatorSecret(clientset)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("A previously created operator secret was identified and deleted")
	}
	return nil
}

func CreateTykOperatorSecret(clientset *kubernetes.Clientset) error {
	secretData := map[string][]byte{
		TykAuth: []byte(data.AppConfig.UserAuth),
		TykOrg:  []byte(data.AppConfig.OrgId),
		TykMode: []byte(TykModePro),
		TykUrl:  []byte(data.AppConfig.DashboardUrl),
	}

	objectMeta := metav1.ObjectMeta{Name: os.Getenv(constants.OperatorSecretNameEnvVarName)}

	secret := v1.Secret{
		ObjectMeta: objectMeta,
		Data:       secretData,
	}
	_, err := clientset.CoreV1().Secrets(os.Getenv(constants.TykPodNamespaceEnvVarName)).
		Create(context.TODO(), &secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
