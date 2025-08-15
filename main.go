package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jvanrhyn/sun/cmd/sun"
	"github.com/jvanrhyn/sun/internal/ui"
)

type status int

const (
	loading status = iota
	success
	errorState
)

type fetchMsg struct{}

type loadedMsg struct{ w sun.Weather }

type errMsg struct{ err error }

type keyMap struct {
	Quit key.Binding
	ToggleFocus key.Binding
	Refresh key.Binding
	Help key.Binding
}

func (k keyMap) ShortHelp() []key.Binding { return []key.Binding{k.Refresh, k.ToggleFocus, k.Help, k.Quit} }
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Refresh, k.ToggleFocus, k.Help, k.Quit}}
}

type model struct {
	table   table.Model
	spin    spinner.Model
	help    help.Model
	keys    keyMap
	status  status
	err     error
	height  int
	width   int
	city    string
	days    int
	client  *sun.Client
	lastUpd time.Time
	th      ui.Theme
}

func initialModel(city string, days int, client *sun.Client) model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	t := table.New(table.WithColumns(defaultColumns()), table.WithRows(nil), table.WithFocused(true))
	h := help.New()
	return model{
		table:  t,
		spin:   sp,
		help:   h,
		th:     ui.DefaultTheme(),
		keys: keyMap{
			Quit:        key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
			ToggleFocus: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "toggle focus")),
			Refresh:     key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
			Help:        key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		},
		status: loading,
		city:   city,
		days:   days,
		client: client,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("Weather Forecast"),
		m.spin.Tick,
		func() tea.Msg { return fetchMsg{} },
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		setAdaptiveDimensions(&m)
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	case fetchMsg:
		m.status = loading
		return m, tea.Batch(m.spin.Tick, fetchWeatherCmd(m.city, m.days, m.client))
	case loadedMsg:
		m.status = success
		m.err = nil
		m.lastUpd = time.Now()
		m.table.SetRows(buildRows(msg.w))
		return m, nil
	case errMsg:
		m.status = errorState
		m.err = msg.err
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.ToggleFocus):
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
			return m, nil
		case key.Matches(msg, m.keys.Refresh):
			return m, func() tea.Msg { return fetchMsg{} }
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	default:
		var cmds []tea.Cmd
		var scmd tea.Cmd
		m.spin, scmd = m.spin.Update(msg)
		cmds = append(cmds, scmd)
		var tcmd tea.Cmd
		m.table, tcmd = m.table.Update(msg)
		cmds = append(cmds, tcmd)
		return m, tea.Batch(cmds...)
	}
}

func (m model) View() string {
	switch m.status {
	case loading:
		return m.th.Base.Render(fmt.Sprintf(" %s Fetching weather for %s (%dd)…", m.spin.View(), m.city, m.days)) + "\n"
	case errorState:
		msg := m.th.Error.Render(fmt.Sprintf("Error: %v", m.err)) + "\nPress r to retry, q to quit."
		return m.th.Base.Render(msg) + "\n"
	case success:
		header := m.th.Header.Render(fmt.Sprintf(" Weather for %s  •  %dd days", m.city, m.days))
		footer := m.th.Footer.Render(fmt.Sprintf("%s  •  Updated %s", m.help.View(m.keys), m.lastUpd.Format(time.Kitchen)))
		return m.th.Base.Render(header+"\n"+m.table.View()) + "\n" + footer + "\n"
	default:
		return m.th.Base.Render("")
	}
}

func defaultColumns() []table.Column {
	return []table.Column{
		{Title: "Time", Width: 6},
		{Title: "Temp °C", Width: 7},
		{Title: "Conditions", Width: 25},
		{Title: "Rain", Width: 5},
		{Title: "Wind", Width: 4},
		{Title: "Gusts", Width: 5},
	}
}

func setAdaptiveDimensions(m *model) {
	min := 6 + 7 + 5 + 4 + 5 + 5 // base for non-Conditions plus spacing
	cond := m.width - min
	if cond < 20 {
		cond = 20
	}
	cols := []table.Column{
		{Title: "Time", Width: 6},
		{Title: "Temp °C", Width: 7},
		{Title: "Conditions", Width: cond},
		{Title: "Rain", Width: 5},
		{Title: "Wind", Width: 4},
		{Title: "Gusts", Width: 5},
	}
	m.table.SetHeight(max(7, m.height-5))
	m.table.SetColumns(cols)
}

func max(a, b int) int { if a > b { return a }; return b }

func buildRows(w sun.Weather) []table.Row {
	var rows []table.Row
	loc, cur, days := w.Location, w.Current, w.Forecast.ForecastDay
	rows = append(rows, table.Row{"------", "-------", fmt.Sprintf("%s, %s", loc.Name, loc.Country), "-----", "----", "-----"})
	rows = append(rows, table.Row{
		"Now",
		fmt.Sprintf("%2.0f", cur.Temperature),
		fmt.Sprintf("%-25s", cur.Condition.Text),
		fmt.Sprintf("%3d%%", cur.ChanceOfRain),
		fmt.Sprintf("%3.0f", cur.WindSpeed),
		fmt.Sprintf("%3.0f", cur.Gusts),
	})
	for _, f := range days {
		hours := f.Hours
		fdate := time.Unix(f.DateEpoch, 0)
		rows = append(rows, table.Row{"-----", "-------", fmt.Sprintf("%s (%s)", fdate.Format("2006-01-02"), fdate.Weekday().String()), "-----", "----", "-----"})
		for _, h := range hours {
			d := time.Unix(h.DateEpoch, 0)
			if d.Before(time.Now()) { continue }
			rows = append(rows, table.Row{
				d.Format("15:04"),
				fmt.Sprintf("%2.0f", h.Temperature),
				fmt.Sprintf("%-25s", h.Condition.Text),
				fmt.Sprintf("%3d%%", h.ChanceOfRain),
				fmt.Sprintf("%3.0f", h.WindSpeed),
				fmt.Sprintf("%3.0f", h.Gusts),
			})
		}
	}
	return rows
}

func fetchWeatherCmd(city string, days int, client *sun.Client) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		w, err := client.Forecast(ctx, city, days)
		if err != nil {
			return errMsg{err}
		}
		return loadedMsg{w: w}
	}
}

func main() {
	_ = godotenv.Load()
	var (
		city string
		days int
	)
	defCity := os.Getenv("DEFAULT_LOCATION")
	if defCity == "" { defCity = "London" }
	defDays := 1
	if v := os.Getenv("NO_OF_DAYS"); v != "" { if n, err := strconvAtoiSafe(v); err == nil { defDays = n } }
	flag.StringVar(&city, "city", defCity, "Enter the name of the city")
	flag.IntVar(&days, "days", defDays, "Number of days to forecast")
	flag.Parse()
	key := os.Getenv("WEATHER_ACCESS_TOKEN")
	client := &sun.Client{APIKey: key, HTTP: &http.Client{Timeout: 10 * time.Second}}
	m := initialModel(city, days, client)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func strconvAtoiSafe(s string) (int, error) { return strconv.Atoi(s) }
