package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	cityFlag string
	daysFlag int

	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
)

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("Weather Forecast")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

// main is the entry point of the application. It loads environment variables,
// parses command-line flags, constructs a URL to call a weather API, and processes
// the response to display current weather conditions and a forecast.
//
// The function does not take any parameters and does not return any values.
func main() {

	// Load environment variables from .env file
	setupEnvironment()

	columns := setupColumns()

	var rows []table.Row

	// Get the number of days to forecast
	var noOfDaysStr = os.Getenv("NO_OF_DAYS")
	noOfDays, err := strconv.Atoi(noOfDaysStr)
	if err != nil {
		noOfDays = 1
	}

	// Parse the command-line flags
	flag.IntVar(&daysFlag, "days", noOfDays, "Number of days to forecast")
	flag.StringVar(&cityFlag, "city", os.Getenv("DEFAULT_LOCATION"), "Enter the name of the city")
	flag.Parse()

	// Retrieve the access token from the environment
	token := os.Getenv("WEATHER_ACCESS_TOKEN")

	// Construct the URL
	url := fmt.Sprintf("https://api.weatherapi.com/v1/forecast.json?q=%s&days=%d&key=%s",
		cityFlag, daysFlag, token)

	// Call the API
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer func() {
		if cerr := res.Body.Close(); cerr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error closing response body: %v\n", cerr)
		}
	}()

	// If Status is not OK 200, panic
	if res.StatusCode != 200 {
		panic("Weather Api not available")
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// Unmarshal the json into the provide
	// struct reference
	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	// Extract data from the struct
	location, current, forecastDay := weather.Location, weather.Current, weather.Forecast.ForecastDay

	rows = append(rows, table.Row{"------", "-------",
		fmt.Sprintf("%s, %s", location.Name, location.Country),
		"-----", "----", "-----"})

	// Build and display current conditions
	rows = append(rows, table.Row{
		"Now",
		fmt.Sprintf("%2.0f", current.Temperature),
		fmt.Sprintf("%-25s", current.Condition.Text),
		fmt.Sprintf("%3d%%", current.ChanceOfRain),
		fmt.Sprintf("%3.0f", current.WindSpeed),
		fmt.Sprintf("%3.0f", current.Gusts),
	})

	for _, fday := range forecastDay {
		hours := fday.Hours

		fdate := time.Unix(fday.DateEpoch, 0).Format("2006-01-02")
		rows = append(rows, table.Row{"-----", "-------", fdate,
			"-----", "----", "-----"})

		// Get the hourly forecasts and
		// construct an output message
		for _, hour := range hours {
			date := time.Unix(hour.DateEpoch, 0)

			// If the hourly forecast is in the past
			// ignore it and continue along
			if date.Before(time.Now()) {
				continue
			}

			rows = append(rows, table.Row{
				date.Format("15:04"),
				fmt.Sprintf("%2.0f", hour.Temperature),
				fmt.Sprintf("%-25s", hour.Condition.Text),
				fmt.Sprintf("%3d%%", hour.ChanceOfRain),
				fmt.Sprintf("%3.0f", hour.WindSpeed),
				fmt.Sprintf("%3.0f", hour.Gusts),
			})

		}
	}

	t := setupTable(columns, rows)

	m := model{t}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func setupColumns() []table.Column {
	columns := []table.Column{
		{Title: "Time", Width: 6},
		{Title: "Temp °C", Width: 7},
		{Title: "Conditions", Width: 25},
		{Title: "Rain", Width: 5},
		{Title: "Wind", Width: 4},
		{Title: "Gusts", Width: 5},
	}
	return columns
}

func setupEnvironment() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func setupTable(columns []table.Column, rows []table.Row) table.Model {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}
