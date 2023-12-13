package electrictymaps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	*http.Client

	AuthToken string
	BaseURL   string
}

func New(token string) *Client {
	return &Client{
		AuthToken: token,
		BaseURL:   "https://api-access.electricitymaps.com/free-tier",
	}
}

type CarbonIntesnityResponse struct {
	CarbonIntensity int `json:"carbonIntensity"`
}

func (c *Client) GetCarbonIntesnity(zone string) (*CarbonIntesnityResponse, error) {
	var r CarbonIntesnityResponse
	url := fmt.Sprintf("%s/carbon-intensity/latest?zone=%s", c.BaseURL, zone)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("client: could not create request: %s\n", err)
	}

	req.Header.Set("auth-token", c.AuthToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client: failed getting response: %s\n", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("client: could not read response body: %s\n", err)
	}

	if err := json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("Can not unmarshal JSON:%s\n", err)
	}

	return &r, nil
}
