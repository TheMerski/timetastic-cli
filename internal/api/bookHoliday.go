package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type LeaveRequest struct {
	UserID               int    `json:"userId"`
	DepartmentID         int    `json:"departmentId"`
	LeaveTypeID          int    `json:"leaveTypeId"`
	Start                string `json:"start"` // Date string
	End                  string `json:"end"`   // Date string
	FromDayPart          int    `json:"fromDayPart"`
	ToDayPart            int    `json:"toDayPart"`
	FromTime             int    `json:"fromTime"`
	ToTime               int    `json:"toTime"`
	Reason               string `json:"reason"`
	AttachmentID         *int   `json:"attachmentId"` // Use pointer to handle null values
	AllDay               bool   `json:"allDay"`
	SpaceBeforeFractions bool   `json:"spaceBeforeFractions"`
}

type BookingResponse struct {
	Success                    bool   `json:"success"`
	Response                   string `json:"response"`
	Rejections                 bool   `json:"rejections"`
	AutoApproved               bool   `json:"autoApproved"`
	HolidayID                  int    `json:"holidayId"`
	OverrideRequired           bool   `json:"overrideRequired"`
	OverrideLockedDateRequired bool   `json:"overrideLockedDateRequired"`
}

func (c *TimetasticClient) BookLeave(department int, leaveType int, start string, end string) (*BookingResponse, error) {
	reqUrl := "https://app.timetastic.co.uk/api/holidays/BookMobile"

	data := LeaveRequest{
		UserID:               c.userData.UserId,
		DepartmentID:         department,
		LeaveTypeID:          leaveType,
		Start:                start,
		End:                  end,
		FromDayPart:          1,
		ToDayPart:            2,
		FromTime:             0,
		AllDay:               false,
		SpaceBeforeFractions: false,
	}
	databytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %v", err)
	}

	req, _ := http.NewRequest("POST", reqUrl, bytes.NewBuffer(databytes))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.makeRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	// Decode the response body
	var response BookingResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response, nil
}
