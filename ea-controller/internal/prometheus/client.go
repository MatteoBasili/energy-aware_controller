package prometheus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	BaseURL string
}

type QueryResult struct {
        Metric map[string]string `json:"metric"`
        Value  []interface{}     `json:"value"`
}

func New(url string) *Client {
	return &Client{BaseURL: url}
}

func (c *Client) Query(q string) ([]QueryResult, error) {

	endpoint := fmt.Sprintf("%s/api/v1/query", c.BaseURL)
	params := url.Values{}
	params.Add("query", q)

	resp, err := http.Get(endpoint + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Result []QueryResult `json:"result"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data.Result, nil
}

func (c *Client) Scalar(q string) float64 {
	res, _ := c.Query(q)
	if len(res) == 0 {
		return 0
	}
	val, _ := strconv.ParseFloat(res[0].Value[1].(string), 64)
	return val
}
