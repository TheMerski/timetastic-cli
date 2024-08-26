package forms

import (
	"github.com/charmbracelet/huh"
	"github.com/themerski/timetastic-cli/internal/api"
)

func GetDepartmentGroup(departments []api.DepartmentResponse, value *int) *huh.Group {
	options := make([]huh.Option[int], len(departments))
	for i, d := range departments {
		options[i] = huh.NewOption(d.Name, d.ID)
	}

	return huh.NewGroup(
		huh.NewSelect[int]().
			Title("Select the department").
			Options(
				options...,
			).
			Value(value),
	)
}
