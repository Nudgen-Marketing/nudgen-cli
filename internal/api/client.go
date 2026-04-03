package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const DefaultBaseURL = "http://localhost:3000/api/v1"

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

type UserSession struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewClientWithToken(token string) *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *Client) delete(path string) error {
	req, err := c.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error string `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		if errorResp.Error != "" {
			return fmt.Errorf("API error: %s", errorResp.Error)
		}
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) post(path string, body interface{}, result interface{}) error {
	req, err := c.newRequest("POST", path, body)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errorResp struct {
			Error string `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		if errorResp.Error != "" {
			return fmt.Errorf("API error: %s", errorResp.Error)
		}
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *Client) patch(path string, body interface{}, result interface{}) error {
	req, err := c.newRequest("PATCH", path, body)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error string `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		if errorResp.Error != "" {
			return fmt.Errorf("API error: %s", errorResp.Error)
		}
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

// GetMe fetches the currently logged in user based on the PAT.
func (c *Client) GetMe() (*UserSession, error) {
	req, err := c.newRequest("GET", "/user/me", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var user UserSession
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Campaigns CRUD

func (c *Client) GetCampaigns() ([]Campaign, error) {
	req, err := c.newRequest("GET", "/campaigns", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result struct {
		Campaigns []Campaign `json:"campaigns"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Campaigns, nil
}

func (c *Client) CreateCampaign(body interface{}) (*Campaign, error) {
	var result struct {
		Campaign Campaign `json:"campaign"`
	}
	err := c.post("/campaigns", body, &result)
	return &result.Campaign, err
}

func (c *Client) DeleteCampaign(id string) error {
	return c.delete("/campaigns/" + id)
}

// Contacts CRUD

func (c *Client) GetContacts() ([]Contact, error) {
	req, err := c.newRequest("GET", "/contacts", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result struct {
		Contacts []Contact `json:"contacts"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Contacts, nil
}

func (c *Client) CreateContact(body interface{}) (string, error) {
	var result struct {
		Contact struct {
			ID string `json:"id"`
		} `json:"contact"`
	}
	err := c.post("/contacts/add", body, &result)
	return result.Contact.ID, err
}

func (c *Client) UpdateContact(id string, body interface{}) (string, error) {
	var result struct {
		ContactID string `json:"contactId"`
	}
	err := c.patch("/contacts/"+id, body, &result)
	return result.ContactID, err
}

func (c *Client) DeleteContact(id string) error {
	return c.delete("/contacts/" + id)
}

// Stats

func (c *Client) GetStats() (*Stats, error) {
	req, err := c.newRequest("GET", "/dashboard/overview", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result struct {
		Stats Stats `json:"stats"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Stats, nil
}

// Teams CRUD

func (c *Client) GetTeams() ([]Team, error) {
	req, err := c.newRequest("GET", "/teams", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result struct {
		Teams []Team `json:"teams"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Teams, nil
}

func (c *Client) SwitchTeam(teamID string) error {
	body := struct {
		TeamID string `json:"teamId"`
	}{
		TeamID: teamID,
	}

	req, err := c.newRequest("POST", "/teams/switch", body)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) CreateTeam(name string) (*Team, error) {
	var result struct {
		Team Team `json:"team"`
	}
	err := c.post("/teams", map[string]string{"name": name}, &result)
	return &result.Team, err
}

func (c *Client) DeleteTeam(id string) error {
	return c.delete("/teams/" + id)
}

// Brand CRUD

func (c *Client) GetBrandConfigs() ([]BrandConfig, error) {
	req, err := c.newRequest("GET", "/brand", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result struct {
		Configs []BrandConfig `json:"configs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Configs, nil
}

func (c *Client) CreateBrandConfig(body interface{}) (*BrandConfig, error) {
	var result struct {
		Config BrandConfig `json:"config"`
	}
	err := c.post("/brand", body, &result)
	return &result.Config, err
}

func (c *Client) UpdateBrandConfig(id string, body interface{}) (*BrandConfig, error) {
	var result struct {
		Config BrandConfig `json:"config"`
	}
	err := c.patch("/brand/"+id, body, &result)
	return &result.Config, err
}

func (c *Client) DeleteBrandConfig(id string) error {
	return c.delete("/brand/" + id)
}

// GetReferralInfo fetches affiliate info.
func (c *Client) GetReferralInfo() (*ReferralInfo, error) {
	req, err := c.newRequest("GET", "/affiliate/me", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result ReferralInfo
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAffiliateAnalytics fetches clicks, conversions, and referrals.
func (c *Client) GetAffiliateAnalytics() (*AffiliateAnalytics, error) {
	req, err := c.newRequest("GET", "/affiliate/analytics", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result AffiliateAnalytics
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAffiliatePayouts fetches the history of withdraw requests.
func (c *Client) GetAffiliatePayouts() ([]PayoutRecord, error) {
	req, err := c.newRequest("GET", "/affiliate/payouts", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result struct {
		Payouts []PayoutRecord `json:"payouts"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Payouts, nil
}
