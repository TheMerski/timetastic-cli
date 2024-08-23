package forms

import (
	"github.com/charmbracelet/huh"
	"github.com/themerski/timetastic-cli/internal/api"
)

func SelectLeaveType(leavetypes []api.LeaveType) int {
	var leaveTypeId int
	options := make([]huh.Option[int], len(leavetypes))
	for i, d := range leavetypes {
		options[i] = huh.NewOption(d.Name, d.ID)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Select the leave type").
				Options(
					options...,
				).
				Value(&leaveTypeId),
		),
	)

	form.Run()
	return leaveTypeId
}
