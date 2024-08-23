package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Define the struct to match the JSON structure
type LeaveType struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	OrganisationID     int    `json:"organisationId"`
	Deducted           bool   `json:"deducted"`
	RequiresApproval   bool   `json:"requiresApproval"`
	IncludeMaxOff      bool   `json:"includeMaxOff"`
	IsPrivate          bool   `json:"isPrivate"`
	Active             bool   `json:"active"`
	IsInUse            bool   `json:"isInUse"`
	Color              string `json:"color"`
	Icon               string `json:"icon"`
	CalendarVisibility int    `json:"calendarVisibility"`
	LimitHours         *int   `json:"limitHours"` // Use pointer to handle null values
	LimitDays          *int   `json:"limitDays"`  // Use pointer to handle null values
}

func (c *TimetasticClient) GetLeaveTypes() (*[]LeaveType, error) {
	url := "https://app.timetastic.co.uk/api/leavetypes"

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
	var response []LeaveType
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response, nil
}
