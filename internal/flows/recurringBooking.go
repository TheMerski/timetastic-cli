package flows

import (
	"log/slog"

	"github.com/themerski/timetastic-cli/internal/api"
	"github.com/themerski/timetastic-cli/internal/forms"
)

func BookRecurringLeave(client *api.TimetasticClient) {
	dep, err := client.GetDepartments()
	if err != nil {
		slog.Error("Failed to get departments", "error", err)
		return
	}
	department := forms.SelectDepartment(*dep)
	slog.Debug("Got department", "departmentId", department)

	leaves, err := client.GetLeaveTypes()
	if err != nil {
		slog.Error("Failed to get departments", "error", err)
		return
	}
	leaveType := forms.SelectLeaveType(*leaves)
	slog.Debug("Got leave type", "leaveTypeId", leaveType)

	when := forms.WhenToCreate()
	forms.BookRecurringLeave(client, department, leaveType, when)

	slog.Debug("Booked all leaves")
}
