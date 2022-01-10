package twitch

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Proxy https://proxy.webshare.io/docs/#the-proxy-list-object
type Proxy struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"proxy_address"`
	Ports    struct {
		Http   int `json:"http"`
		Socks5 int `json:"socks5"`
	} `json:"ports"`
	LastVerification      string  `json:"last_verification"`
	CountryCode           string  `json:"country_code"`
	CountryCodeConfidence float32 `json:"country_code_confidence"`
	CityName              string  `json:"city_name"`
}

// ReplacementInfo https://proxy.webshare.io/docs/#the-replacement-info-object
type ReplacementInfo struct {
	RefreshLastAt string `json:"automatic_refresh_last_at,omitempty"`
	RefreshNextAt string `json:"automatic_refresh_next_at,omitempty"`
}

// GetNextRefresh gets the duration between the next refresh and the local time
func (r *ReplacementInfo) GetNextRefresh() (*time.Duration, error) {
	next, err := time.Parse(time.RFC3339, r.RefreshNextAt)
	if err != nil {
		log.Printf("[GetNextRefresh @ time.Parse] %v\n", err)
		return nil, fmt.Errorf("failed to parse next refresh")
	}

	when := next.Sub(time.Now().In(next.Location()))

	return &when, nil
}

type List struct {
	Proxies map[Proxy]bool
	Info    ReplacementInfo
}

// ListFromAPI collects the Proxy list and the replacement info
func ListFromAPI() (*List, error) {
	listURL := "https://proxy.webshare.io/api/proxy/list/"
	replacementURL := "https://proxy.webshare.io/api/proxy/replacement/info/"

	// wrapper for Proxy results
	results := struct {
		Proxies []Proxy `json:"results"`
	}{}

	// get Proxy list
	b, err := doRequest(listURL)
	if err != nil {
		log.Printf("[Proxy.ListFromAPI @ doRequest(listURL)] %v\n", err)
		return nil, fmt.Errorf("failed to retrieve Proxy list")
	}

	// unmarshal Proxy list
	err = json.Unmarshal(b, &results)
	if err != nil {
		log.Printf("[Proxy.ListFromAPI @ json.Unmarshal] %v\n", err)
		return nil, fmt.Errorf("failed to unmarshal Proxy list")
	}

	// get replacement info
	b, err = doRequest(replacementURL)
	if err != nil {
		log.Printf("[Proxy.ListFromAPI @ doRequest(replacementURL)] %v\n", err)
		return nil, fmt.Errorf("failed to retrieve Proxy replacement info")
	}

	var info ReplacementInfo

	// unmarshal Proxy info
	err = json.Unmarshal(b, &info)
	if err != nil {
		log.Printf("[Proxy.ListFromAPI @ json.Unmarshal] %v\n", err)
		return nil, fmt.Errorf("failed to unmarshal Proxy replacement info")
	}

	// setup proxies map
	proxies := make(map[Proxy]bool)
	for _, proxy := range results.Proxies {
		proxies[proxy] = false
	}

	return &List{
		Proxies: proxies,
		Info:    info,
	}, nil
}

// Does the requests with the strapped up authorization required to make the requests
func doRequest(uri string) ([]byte, error) {
	client := http.DefaultClient

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Printf("[Proxy.doRequest @ http.NewRequest] %v\n", err)
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	token := fmt.Sprintf("Token %s", os.Getenv("WEBSHARE_API_KEY"))
	req.Header.Add("Authorization", token)

	resp, err := client.Do(req)
	if resp.StatusCode == http.StatusTooManyRequests {
		log.Printf("[Proxy.doRequest @ client.Do] %v %v\n", err, resp.StatusCode)
		return nil, fmt.Errorf("api rate-limit exceeded")
	}

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("[Proxy.doRequest @ client.Do] %v %v\n", err, resp.StatusCode)
		return nil, fmt.Errorf("failed to do request")
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[Proxy.doRequest @ io.ReadAll] %v\n", err)
		return nil, fmt.Errorf("failed to read request body")
	}

	return b, nil
}
