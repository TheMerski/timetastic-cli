package forms

import (
	"github.com/charmbracelet/huh"
	"github.com/themerski/timetastic-cli/internal/api"
)

func SelectDepartment(departments []api.DepartmentResponse) int {
	var departmentId int
	options := make([]huh.Option[int], len(departments))
	for i, d := range departments {
		options[i] = huh.NewOption(d.Name, d.ID)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Select the department").
				Options(
					options...,
				).
				Value(&departmentId),
		),
	)

	form.Run()
	return departmentId
}
