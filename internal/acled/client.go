package acled

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	Email    string
	Password string
	Token    string
	Expiry   time.Time
}

func (c *Client) Authenticate() error {
	data := url.Values{
		"grant_type": {"password"},
		"username":   {c.Email},
		"password":   {c.Password},
		"client_id":  {"acled"},
	}

	req, err := http.NewRequest("POST", "https://acleddata.com/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (acled-archiver-bot)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("AUTH FAILED (%d): %s", resp.StatusCode, string(body))
	}

	var res struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	c.Token = res.AccessToken
	c.Expiry = time.Now().Add(time.Duration(res.ExpiresIn) * time.Second)
	log.Println("🔑 OAuth Token refreshed (valid for 14 hours)")
	return nil
}

func (c *Client) FetchPage(year int, page int) ([]Event, error) {
	// Proactive token refresh
	if c.Token == "" || time.Now().After(c.Expiry.Add(-5*time.Minute)) {
		if err := c.Authenticate(); err != nil {
			return nil, err
		}
	}

	apiURL := fmt.Sprintf("https://acleddata.com/api/acled/read?year=%d&page=%d&limit=5000&population=full", year, page)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("User-Agent", "Mozilla/5.0 (acled-archiver-bot)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiRes APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiRes); err != nil {
		return nil, err
	}

	return apiRes.Data, nil
}
