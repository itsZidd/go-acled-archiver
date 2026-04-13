package acled

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings" // Ensure this is imported
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

	// Correctly encode the form data into the request body
	req, err := http.NewRequest("POST", "https://acleddata.com/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	// Explicitly set headers for the OAuth POST request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("AUTH FAILED (Status %d): %s", resp.StatusCode, string(body))
	}

	var res struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("auth decode error: %v", err)
	}

	c.Token = res.AccessToken
	c.Expiry = time.Now().Add(time.Duration(res.ExpiresIn) * time.Second)
	log.Println("Successfully authenticated with ACLED.")
	return nil
}

func (c *Client) FetchPage(year int, page int) ([]Event, error) {
	if c.Token == "" || time.Now().After(c.Expiry.Add(-5*time.Minute)) {
		if err := c.Authenticate(); err != nil {
			return nil, err
		}
	}

	apiURL := fmt.Sprintf("https://acleddata.com/api/acled/read?year=%d&page=%d&limit=5000", year, page)
	req, _ := http.NewRequest("GET", apiURL, nil)

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		body, _ := io.ReadAll(resp.Body)
		snippet := string(body)
		if len(snippet) > 150 {
			snippet = snippet[:150]
		}
		return nil, fmt.Errorf("HTML Error (Access Denied?): %s", snippet)
	}

	var apiRes APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiRes); err != nil {
		return nil, err
	}

	return apiRes.Data, nil
}
