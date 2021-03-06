package horoscope

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Reading struct {
	Summary string `json:"Horoscope"`
	Date    string `json:"Date"`
	Sign    string `json:"Sign"`
}

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL: "https://ohmanda.com/api/horoscope",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c Client) FormatURL(sign string) string {
	sign = url.QueryEscape(sign)
	return fmt.Sprintf("%s/%s", c.BaseURL, sign)
}

func (c *Client) GetReading(sign string) (Reading, error) {
	URL := c.FormatURL(sign)
	resp, err := c.HTTPClient.Get(URL)
	if err != nil {
		return Reading{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return Reading{}, fmt.Errorf("could not find sign: %s ", sign)
	}
	if resp.StatusCode != http.StatusOK {
		return Reading{}, fmt.Errorf("unexpected response status %q", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Reading{}, err
	}
	readings, err := ParseResponse(data)
	if err != nil {
		return Reading{}, err
	}
	return readings, nil
}

func ParseResponse(data []byte) (Reading, error) {
	var resp Reading
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return Reading{}, fmt.Errorf("invalid api response %s: %w", data, err)
	}

	reading := Reading{
		Summary: resp.Summary,
		Date:    resp.Date,
		Sign:    resp.Sign,
	}
	return reading, nil
}

func FormatURL(baseURL, sign string) string {
	return fmt.Sprintf("%s/%s", baseURL, sign)
}

func Get(sign string) (Reading, error) {
	c := NewClient()
	readings, err := c.GetReading(sign)
	if err != nil {
		return Reading{}, err
	}
	return readings, nil
}

func RunCLI(sign string) (Reading, error) {
	readings, err := Get(sign)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return readings, nil
}
