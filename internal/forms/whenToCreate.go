package forms

import (
	"errors"
	"strconv"
	"time"

	"github.com/charmbracelet/huh"
)

type CreationData struct {
	Day           time.Weekday
	StartDate     string
	EndDate       *string
	WeeksToCreate int
	WeeksToAdd    int
}

func WhenToCreate() CreationData {
	var endDate string
	var data CreationData

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[time.Weekday]().
				Title("For which day do you want to create an entry").
				Options(
					huh.NewOption("Monday", time.Monday),
					huh.NewOption("Tuesday", time.Tuesday),
					huh.NewOption("Wednesday", time.Wednesday),
					huh.NewOption("Thursday", time.Thursday),
					huh.NewOption("Friday", time.Friday),
				).
				Value(&data.Day),
		),
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("How often should it be created").
				Options(
					huh.NewOption("Weekly", 1),
					huh.NewOption("Bi-weekly", 2),
				).
				Value(&data.WeeksToAdd),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("From which date do you want to start").
				Prompt("?").
				Validate(validateDate).
				Value(&data.StartDate),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("When do you want to end, or how many weeks do you want to create").
				Prompt("?").
				Validate(validateDateOrInt).
				Value(&endDate),
		),
	)

	form.Run()

	if validateInt(endDate) == nil {
		data.WeeksToCreate, _ = strconv.Atoi(endDate)
	} else {
		data.EndDate = &endDate
	}

	return data
}

func validateDateOrInt(input string) error {
	if err := validateDate(input); err != nil {
		if intErr := validateInt(input); intErr != nil {
			return errors.New("Please enter a valid date (YYYY-MM-DD) or a number")
		}
	}
	return nil
}

func validateDate(input string) error {
	_, err := time.Parse(time.DateOnly, input)
	if err != nil {
		return errors.New("Please enter a valid date (YYYY-MM-DD)")
	}
	return nil
}

func validateInt(input string) error {
	_, err := strconv.Atoi(input)
	if err != nil {
		return errors.New("Please enter a valid number")
	}
	return nil
}
