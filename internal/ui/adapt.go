package ui

import "github.com/charmbracelet/bubbles/table"

func AdaptiveColumns(width int) []table.Column {
	min := 6 + 7 + 5 + 4 + 5 + 5
	cond := width - min
	if cond < 20 {
		cond = 20
	}
	return []table.Column{
		{Title: "Time", Width: 6},
		{Title: "Temp °C", Width: 7},
		{Title: "Conditions", Width: cond},
		{Title: "Rain", Width: 5},
		{Title: "Wind", Width: 4},
		{Title: "Gusts", Width: 5},
	}
}
