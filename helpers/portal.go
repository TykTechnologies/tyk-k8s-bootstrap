package helpers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"tyk/tyk/bootstrap/data"

	"k8s.io/apimachinery/pkg/util/json"
)

func BootstrapPortal(client http.Client) error {
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
