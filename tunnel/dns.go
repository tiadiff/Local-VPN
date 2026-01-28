package tunnel

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

// DoHResponse is the simplified structure of Cloudflare's JSON response
type DoHResponse struct {
	Status int `json:"Status"` // 0 = Success
	Answer []struct {
		Name string `json:"name"`
		Type int    `json:"type"` // 1 = A (IPv4)
		Data string `json:"data"` // The IP address
	} `json:"Answer"`
}

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

// ResolveDoH resolves a domain name to an IPv4 address using Cloudflare DoH
func ResolveDoH(domain string) (string, error) {
	// If it's already an IP, return it
	if net.ParseIP(domain) != nil {
		return domain, nil
	}

	url := fmt.Sprintf("https://1.1.1.1/dns-query?name=%s&type=A", domain)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/dns-json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("DoH request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("DoH bad status: %d", resp.StatusCode)
	}

	var doh DoHResponse
	if err := json.NewDecoder(resp.Body).Decode(&doh); err != nil {
		return "", fmt.Errorf("DoH decode failed: %v", err)
	}

	if doh.Status != 0 {
		return "", fmt.Errorf("DoH DNS error status: %d", doh.Status)
	}

	for _, answer := range doh.Answer {
		if answer.Type == 1 { // A Record (IPv4)
			return answer.Data, nil
		}
	}

	return "", fmt.Errorf("no A record found for %s", domain)
}
