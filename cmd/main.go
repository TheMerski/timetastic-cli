package main

import (
	"log/slog"

	"github.com/charmbracelet/huh"
	"github.com/themerski/timetastic-cli/internal/api"
	"github.com/themerski/timetastic-cli/internal/forms"
)

func main() {
	client := api.NewTimetasticClient()

	dep, err := client.GetDepartments()
	if err != nil {
		slog.Error("Failed to get departments", "error", err)
		return
	}
	department := forms.SelectDepartment(*dep)
	slog.Info("Got department", "departmentId", department)

	leaves, err := client.GetLeaveTypes()
	if err != nil {
		slog.Error("Failed to get departments", "error", err)
		return
	}
	leaveType := forms.SelectLeaveType(*leaves)
	slog.Info("Got leave type", "leaveTypeId", leaveType)

	var startDate string
	input := huh.NewInput().
		Title("Which date").
		Prompt("?").
		Value(&startDate)

	input.Run()

	res, err := client.BookLeave(department, leaveType, startDate, startDate)
	if err != nil {
		slog.Error("Failed to book leave", "error", err)
		return
	}
	slog.Info("Leave booked", "response", res)
}
