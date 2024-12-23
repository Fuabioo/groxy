package cmd

import (
	"strings"

	catpuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var (
	theme     = catpuccin.Mocha
	highlight = lipgloss.NewStyle().
			Bold(true).
			Underline(true).
			Foreground(lipgloss.Color(theme.Peach().Hex)).
			Render
)

func init() {
	log.SetStyles(&log.Styles{
		Timestamp: lipgloss.NewStyle().Faint(true),
		Caller:    lipgloss.NewStyle().Faint(true),
		Prefix:    lipgloss.NewStyle().Bold(true).Faint(true),
		Message:   lipgloss.NewStyle(),
		Key: lipgloss.NewStyle().Faint(true).
			Foreground(lipgloss.Color(theme.Pink().Hex)),
		Value: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Peach().Hex)),
		Separator: lipgloss.NewStyle().Faint(true),
		Levels: map[log.Level]lipgloss.Style{
			log.DebugLevel: lipgloss.NewStyle().
				SetString(strings.ToUpper(log.DebugLevel.String())).
				Bold(true).
				MaxWidth(4).
				Foreground(lipgloss.Color(theme.Lavender().Hex)),
			log.InfoLevel: lipgloss.NewStyle().
				SetString(strings.ToUpper(log.InfoLevel.String())).
				Bold(true).
				MaxWidth(4).
				Foreground(lipgloss.Color(theme.Blue().Hex)),
			log.WarnLevel: lipgloss.NewStyle().
				SetString(strings.ToUpper(log.WarnLevel.String())).
				Bold(true).
				MaxWidth(4).
				Foreground(lipgloss.Color(theme.Yellow().Hex)),
			log.ErrorLevel: lipgloss.NewStyle().
				SetString(strings.ToUpper(log.ErrorLevel.String())).
				Bold(true).
				MaxWidth(4).
				Foreground(lipgloss.Color(theme.Maroon().Hex)),
			log.FatalLevel: lipgloss.NewStyle().
				SetString(strings.ToUpper(log.FatalLevel.String())).
				Bold(true).
				MaxWidth(4).
				Foreground(lipgloss.Color(theme.Red().Hex)),
		},
		Keys: map[string]lipgloss.Style{},
		Values: map[string]lipgloss.Style{
			"domain": lipgloss.NewStyle().
				Italic(true).
				Foreground(lipgloss.Color(theme.Mauve().Hex)),
			"host": lipgloss.NewStyle().
				Italic(true).
				Foreground(lipgloss.Color(theme.Mauve().Hex)),
			"scheme": lipgloss.NewStyle().
				Italic(true).
				Bold(true).
				Foreground(lipgloss.Color(theme.Teal().Hex)),
		},
	})
}

func colorizeSignature(signature string) string {
	// for each word, colorize it using a color array
	// for example, the first word is colored with the first color in the array
	// the second word is colored with the second color in the array
	// and so on
	colours := []string{
		theme.Pink().Hex,
		theme.Mauve().Hex,
		theme.Red().Hex,
		theme.Peach().Hex,
		theme.Yellow().Hex,
		theme.Green().Hex,
		theme.Sapphire().Hex,
		theme.Blue().Hex,
		theme.Lavender().Hex,
	}
	lines := strings.Split(signature, "\n")
	baseStyle := lipgloss.NewStyle().Bold(true)
	var colorized strings.Builder
	for i, line := range lines {
		words := strings.Fields(line)
		for j, word := range words {
			// colour gradient
			jindex := i + j
			colour := colours[jindex%len(colours)]
			colorized.WriteString(
				baseStyle.
					Foreground(lipgloss.Color(colour)).
					Render(word),
			)
			colorized.WriteString(" ")
		}
		if i < len(lines)-1 {
			colorized.WriteString("\n")
		}
	}
	return colorized.String()
}
