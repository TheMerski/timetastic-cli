package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Define the struct to match the JSON structure
type DepartmentResponse struct {
	ID             int    `json:"id"`
	OrganisationID int    `json:"organisationId"`
	Name           string `json:"name"`
	ManagerID      int    `json:"managerId"`
	BossID         int    `json:"bossId"`
	Archived       bool   `json:"archived"`
	UserCount      int    `json:"userCount"`
	MaxOff         int    `json:"maxOff"`
}

func (c *TimetasticClient) GetDepartments() (*[]DepartmentResponse, error) {
	url := "https://app.timetastic.co.uk/api/departments"

	req, _ := http.NewRequest("GET", url, nil)

	resp, err := c.makeRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	// Decode the response body
	var response []DepartmentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response, nil
}
