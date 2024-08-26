package forms

import (
	"github.com/charmbracelet/huh"
	"github.com/themerski/timetastic-cli/internal/api"
)

func GetLeaveTypesGroup(leavetypes []api.LeaveType, value *int) *huh.Group {
	options := make([]huh.Option[int], len(leavetypes))
	for i, d := range leavetypes {
		options[i] = huh.NewOption(d.Name, d.ID)
	}

	return huh.NewGroup(
		huh.NewSelect[int]().
			Title("Select the leave type").
			Options(
				options...,
			).
			Value(value),
	)
}
