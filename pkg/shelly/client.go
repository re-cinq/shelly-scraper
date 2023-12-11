package shelly

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	*http.Client

	addr string
}

func New(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

type SwitchStatusResponse struct {
	ID      int     `json:"id"`
	Source  string  `json:"source"`
	Output  bool    `json:"ouput"`
	APower  float32 `json:"apower"`
	Voltage float32 `json:"voltage"`
	Current float32 `json:"current"`
	AEnergy AEnergy `json:"aenergy"`
}

type AEnergy struct {
	Total           float32   `json:"total"`
	ByMinute        []float32 `json:"by_minute"`
	MinuteTimestamp int       `json:"minute_ts"`
}

func (c *Client) GetSwitchStatus(id string) (*SwitchStatusResponse, error) {
	url := fmt.Sprintf("http://%s/rpc/Switch.GetStatus?id=%s", c.addr, id)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("client: could not create request: %s\n", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client: failed getting response: %s\n", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("client: could not read response body: %s\n", err)
	}

	var r SwitchStatusResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("Can not unmarshal JSON:%s\n", err)
	}
	return &r, nil
}
