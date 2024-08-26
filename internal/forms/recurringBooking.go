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
}

type doneMsg struct {
	created int
}

func (r resultMsg) String() string {
	return fmt.Sprintf("%s %s, %s", r.date, r.result,
		r.message)
}

type model struct {
	spinner  spinner.Model
	results  []resultMsg
	quitting bool
	created  int
}

func newModel() model {
	const numLastResults = 5
	s := spinner.New()
	s.Style = spinnerStyle
	return model{
		spinner: s,
		results: make([]resultMsg, numLastResults),
		created: 0,
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
		m.created = msg.created
		return m, tea.Quit
	case resultMsg:
		m.results = append(m.results[1:], msg)
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

	if m.quitting {
		s += "You're all set! Handled " + fmt.Sprintf("%d", m.created) + " bookings."
	} else {
		s += m.spinner.View() + " Creating bookings..."
	}

	s += "\n\n"

	for _, res := range m.results {
		s += res.String() + "\n"
	}

	if !m.quitting {
		s += helpStyle.Render("Press any key to exit")
	}

	if m.quitting {
		s += "\n"
	}

	return appStyle.Render(s)
}

func BookRecurringLeave(client *api.TimetasticClient, department int, leaveType int, data CreationData) {
	p := tea.NewProgram(newModel())
	date, err := time.Parse(time.DateOnly, data.StartDate)
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
	count := 0

	go func() {
		for date.Before(endDate) {
			// Book leave
			res, err := client.BookLeave(department, leaveType, date.Format(time.DateOnly), date.Format(time.DateOnly))
			if err != nil {
				fmt.Println("Error booking leave:", err)
				continue
			}
			date = date.AddDate(0, 0, 7*data.WeeksToAdd)
			count++

			resMessage := "Created successfully"
			if !res.Success {
				resMessage = "Failed to create"
			}

			// Send the Bubble Tea program a message from outside the
			// tea.Program. This will block until it is ready to receive
			// messages.
			p.Send(resultMsg{date.Format(time.DateOnly), resMessage, res.Response})
		}
		p.Send(doneMsg{created: count})
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
