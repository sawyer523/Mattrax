package azuread

import (
	"io"
	"net/http"
	"time"

	mattrax "github.com/mattrax/Mattrax/internal"
)

const microsoftGraphURL = "https://graph.microsoft.com/v1.0"

// TODO: Global for Mattrax
var client = &http.Client{
	Timeout: 5 * time.Second,
}

type Service struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"` // todo go time
	AccessToken string `json:"access_token"`
}

func (aad *Service) RefreshAccessToken(server *mattrax.Server) error {
	// formData := url.Values{}
	// formData.Add("client_id", server.Settings.Windows.AADMDMApplicationID)
	// formData.Add("scope", "https://graph.microsoft.com/.default")
	// formData.Add("client_secret", server.Settings.Windows.AADMDMClientSecret)
	// formData.Add("grant_type", "client_credentials")

	// resp, err := client.Post("https://login.microsoftonline.com/"+url.QueryEscape(server.Settings.Windows.AADTenantName)+"/oauth2/v2.0/token", "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	// if err != nil {
	// 	return err
	// }

	// var aadUpdated *Service // TODO: Handle errors response
	// if err := json.NewDecoder(resp.Body).Decode(aadUpdated); err != nil {
	// 	return err
	// }

	// aad = aadUpdated

	return nil
}

func (aad Service) NewRequest(method string, path string, body io.Reader) (*http.Client, *http.Request, error) {
	req, err := http.NewRequest("GET", microsoftGraphURL+"/devices", body)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("Authorization", "Bearer "+aad.AccessToken)

	return client, req, nil
}
