package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type NonWorkingDay struct {
	ID        interface{} `json:"id"`
	DayOffset int         `json:"dayOffset"`
	DayPart   int         `json:"dayPart"`
	Name      interface{} `json:"name"`
	Icon      interface{} `json:"icon"`
}

type Holiday struct {
	Day      int    `json:"day"`
	AmHolID  int    `json:"am_hol_id"`
	PmHolID  int    `json:"pm_hol_id"`
	AmLt     string `json:"am_lt"`
	PmLt     string `json:"pm_lt"`
	AmPend   bool   `json:"am_pend"`
	PmPend   bool   `json:"pm_pend"`
	AmLtName string `json:"am_lt_name"`
	PmLtName string `json:"pm_lt_name"`
	Icon     string `json:"icon"`
	IsStart  bool   `json:"isStart"`
	IsEnd    bool   `json:"isEnd"`
}

type User struct {
	Days            int             `json:"days"`
	Name            string          `json:"name"`
	Remaining       string          `json:"remaining"`
	DeptName        string          `json:"deptName"`
	UserID          int             `json:"userId"`
	GravatarURL     string          `json:"gravatarUrl"`
	NonWorkingDays  []NonWorkingDay `json:"nonworkingdays"`
	Holidays        []Holiday       `json:"holidays"`
	LockedDays      []interface{}   `json:"lockeddays"`
	DepartmentID    int             `json:"departmentId"`
	Birthday        string          `json:"birthday"`
	WorkAnniversary string          `json:"workAnniversary"`
	IsDeptManager   bool            `json:"isDeptManager"`
	CanViewCalendar bool            `json:"canViewCalendar"`
	CanManage       bool            `json:"canManage"`
	CanEdit         bool            `json:"canEdit"`
	FirstName       string          `json:"firstname"`
	Surname         string          `json:"surname"`
	IsFavourite     bool            `json:"isFavourite"`
	Initials        string          `json:"initials"`
	StartDate       string          `json:"startDate"`
	Year            int             `json:"year"`
}

type Day struct {
	DayChar   string `json:"dayChar"`
	DayNumber int    `json:"dayNumber"`
	Month     int    `json:"month"`
	IsToday   bool   `json:"isToday"`
}

type WallchartData struct {
	Start          string      `json:"start"`
	End            string      `json:"end"`
	Users          []User      `json:"users"`
	Days           []Day       `json:"days"`
	EndText        interface{} `json:"endText"`
	StartText      interface{} `json:"startText"`
	MonthStart     int         `json:"monthstart"`
	Year           int         `json:"year"`
	CanPageForward bool        `json:"canPageForward"`
	CanPageBack    bool        `json:"canPageBack"`
	IsRestricted   bool        `json:"isRestricted"`
}

func (c *TimetasticClient) GetWallchart(departmentId int) (*WallchartData, error) {
	reqUrl := "https://app.timetastic.co.uk/Wallchart/GetWallChartPage"

	data := url.Values{
		"offset":          {"0"},
		"start":           {time.Now().Format(time.DateOnly)},
		"end":             {time.Now().AddDate(0, 1, 0).Format(time.DateOnly)},
		"accountTypeId":   {"0"},
		"departmentIds[]": {strconv.Itoa(departmentId)},
		"sortOrder":       {"0"},
	}

	slog.Info("data", "data", data.Encode())

	req, _ := http.NewRequest("POST", reqUrl, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := c.makeRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		slog.Error("Received non-200 response", "statusCode", resp.StatusCode, "status", resp.Status)
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	// Decode the response body
	var response WallchartData
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response, nil
}
