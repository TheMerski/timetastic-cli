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
	leaves, err := client.GetLeaveTypes()
	if err != nil {
		slog.Error("Failed to get departments", "error", err)
		return
	}

	recurringBookingData := forms.RecurringBookingForm(*dep, *leaves)
	forms.BookRecurringLeave(client, recurringBookingData)

	slog.Debug("Booked all leaves")
}
