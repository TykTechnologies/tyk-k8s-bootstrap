package tyk

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"tyk/tyk/bootstrap/tyk/api"
)

func (s *Service) BoostrapPortal() error {
	err := s.CreatePortalDefaultSettings()
	if err != nil {
		return err
	}

	err = s.InitialiseCatalogue()
	if err != nil {
		return err
	}

	err = s.CreatePortalHomepage()
	if err != nil {
		return err
	}

	err = s.SetPortalCname()
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) SetPortalCname() error {
	fmt.Println("Setting portal cname")

	cnameReq := api.CnameRequest{Cname: s.appArgs.Cname}
	reqBody, err := json.Marshal(cnameReq)
	if err != nil {
		return err
	}
	reqData := bytes.NewReader(reqBody)

	req, err := http.NewRequest("PUT", s.appArgs.DashboardUrl+ApiPortalCnameEndpoint, reqData)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", s.appArgs.UserAuth)

	res, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("failed to set portal cname")
	}

	return nil
}

func (s *Service) InitialiseCatalogue() error {
	fmt.Println("Initialising Catalogue")

	initCatalog := api.InitCatalogReq{OrgId: s.appArgs.OrgId}
	reqBody, err := json.Marshal(initCatalog)
	if err != nil {
		return err
	}
	reqData := bytes.NewReader(reqBody)
	req, err := http.NewRequest("POST", s.appArgs.DashboardUrl+ApiPortalCatalogueEndpoint, reqData)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", s.appArgs.UserAuth)

	res, err := s.httpClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return err
	}
	resp := api.DashboardGeneralResponse{}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(bodyBytes, &resp)
	if err != nil {
		return err
	}
	s.appArgs.CatalogId = resp.Message

	return nil
}

func (s *Service) CreatePortalHomepage() error {
	fmt.Println("Creating portal homepage")

	homepageContents := GetPortalHomepage()
	reqBody, err := json.Marshal(homepageContents)
	if err != nil {
		return err
	}

	reqData := bytes.NewReader(reqBody)
	req, err := http.NewRequest("POST", s.appArgs.DashboardUrl+ApiPortalPagesEndpoint, reqData)
	req.Header.Set("Authorization", s.appArgs.UserAuth)
	if err != nil {
		return err
	}
	res, err := s.httpClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return err
	}

	resp := api.DashboardGeneralResponse{}

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

func GetPortalHomepage() api.PortalHomepageRequest {
	return api.PortalHomepageRequest{
		IsHomepage:   true,
		TemplateName: "",
		Title:        "Developer portal name",
		Slug:         "/",
		Fields: api.PortalFields{
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

func (s *Service) CreatePortalDefaultSettings() error {
	fmt.Println("Creating bootstrap default settings")

	req, err := http.NewRequest("POST", s.appArgs.DashboardUrl+ApiPortalConfigurationEndpoint, nil)
	req.Header.Set("Authorization", s.appArgs.UserAuth)

	if err != nil {
		return err
	}
	res, err := s.httpClient.Do(req)
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
