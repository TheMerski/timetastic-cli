package forms

// A simple example that shows how to send messages to a Bubble Tea program
// from outside the program using Program.Send(Msg).

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/themerski/timetastic-cli/internal/api"
)

var (
	spinnerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(1, 0)
	dotStyle      = helpStyle.UnsetMargins()
	durationStyle = dotStyle
	appStyle      = lipgloss.NewStyle().Margin(1, 2, 0, 2)
)

type resultMsg struct {
	date    string
	result  string
	message string
	success bool
}

type doneMsg struct {
	success int
	failed  int
}

func (r resultMsg) String() string {
	return fmt.Sprintf("%s %s, %s", r.date, r.result,
		r.message)
}

type model struct {
	spinner       spinner.Model
	results       []resultMsg
	failedResults []resultMsg
	quitting      bool
	succeed       int
	failed        int
}

func newModel() model {
	const numLastResults = 5
	s := spinner.New()
	s.Style = spinnerStyle
	return model{
		spinner:       s,
		results:       make([]resultMsg, numLastResults),
		failedResults: make([]resultMsg, 5),
		succeed:       0,
		failed:        0,
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit
	case doneMsg:
		m.quitting = true
		m.succeed = msg.success
		m.failed = msg.failed
		return m, tea.Quit
	case resultMsg:
		m.results = append(m.results[1:], msg)
		if !msg.success {
			m.failedResults = append(m.failedResults, msg)
		}
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m model) View() string {
	var s string

	s += "\n\n"

	for _, res := range m.results {
		s += res.String() + "\n"
	}

	if !m.quitting {
		s += helpStyle.Render("Press any key to exit")
	}

	if m.quitting {
		s = "You're all set!\n\n"
		s += fmt.Sprintf("Successfully created %d bookings\n", m.succeed)
		s += fmt.Sprintf("Failed to create %d bookings:\n", m.failed)
		for _, res := range m.results {
			if res.result == "Failed to create" {
				s += "\t" + res.String() + "\n"
			}
		}
	} else {
		s += m.spinner.View() + " Creating bookings..."
	}

	return appStyle.Render(s)
}

func BookRecurringLeave(client *api.TimetasticClient, department int, leaveType int, data CreationData) {
	p := tea.NewProgram(newModel())
	date, err := time.Parse(time.DateOnly, data.StartDate)
	slog.Info("Starting bookings", "start", date, "startDay", date.Weekday(), "end", data.EndDate, "weeksToAdd", data.WeeksToAdd, "weeksToCreate", data.WeeksToCreate)
	if err != nil {
		slog.Error("Failed to parse start date, returning", "error", err)
		return
	}
	// forwards the date to the next weekday
	for date.Weekday() != data.Day {
		date = date.AddDate(0, 0, 1)
	}

	var endDate time.Time
	if data.EndDate == nil {
		endDate = date.Add(time.Hour * 24 * 7 * time.Duration(data.WeeksToCreate))
	} else {
		endDate, err = time.Parse(time.DateOnly, *data.EndDate)
		if err != nil {
			slog.Error("Failed to parse end date, returning", "error", err)
			return
		}
	}
	successCount := 0
	failedCount := 0

	go func() {
		for date.Before(endDate) || date.Equal(endDate) {
			slog.Info("Creating bookings", "start", date, "startDay", date.Weekday(), "end", data.EndDate, "weeksToAdd", data.WeeksToAdd, "weeksToCreate", data.WeeksToCreate)
			// Book leave
			res, err := client.BookLeave(department, leaveType, date.Format(time.DateOnly), date.Format(time.DateOnly))
			if err != nil {
				fmt.Println("Error booking leave:", err)
				continue
			}
			resMessage := "Created successfully"
			if !res.Success {
				resMessage = "Failed to create"
				failedCount++
			} else {
				successCount++
			}
			// Send the Bubble Tea program a message from outside the
			// tea.Program. This will block until it is ready to receive
			// messages.
			p.Send(resultMsg{date.Format(time.DateOnly), resMessage, res.Response, res.Success})

			// Move to the next date to book
			date = date.AddDate(0, 0, 7*data.WeeksToAdd)
		}
		p.Send(doneMsg{success: successCount, failed: failedCount})
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
