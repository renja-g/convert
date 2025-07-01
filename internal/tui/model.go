package tui

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/renja-g/convert/internal/alias"
	"github.com/renja-g/convert/internal/converter"
	"github.com/renja-g/convert/internal/detect"
)

// fileInfo holds meta data of a processed file.
// It is used as Bubble Tea message payload as well.
// (We declare the message alias further down.)
type fileInfo struct {
	path     string
	mimeType string
	ext      string // extension (with leading dot)
	err      error
}

type fileInfoMsg fileInfo

type convertDoneMsg struct {
	outputPath string
	err        error
}

// resetMsg is sent after showing a success or error banner to return to the start screen.
type resetMsg struct{}

// model represents the Bubble Tea application state.

type model struct {
	styles      styles
	spinner     spinner.Model
	processing  bool // while detecting mime type
	converting  bool // while actual conversion happens
	file        fileInfo
	choices     []string // destination extensions (with dot)
	cursor      int
	choice      string // chosen destination extension
	width       int
	height      int
	input       textinput.Model
	suggestions []string
	selIdx      int
	searching   bool

	// Post-conversion feedback
	showSuccess bool
	successPath string
}

func initialModel() model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()

	return model{
		styles:      defaultStyles(),
		spinner:     sp,
		input:       ti,
		suggestions: searchFiles(""),
		searching:   true,
	}
}

// checkFileCmd inspects the dropped file and returns a fileInfoMsg.
func checkFileCmd(path string) tea.Cmd {
	return func() tea.Msg {
		file, err := os.Open(path)
		if err != nil {
			return fileInfoMsg{path: path, err: fmt.Errorf("could not open file: %w", err)}
		}
		defer file.Close()

		buffer := make([]byte, 512)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return fileInfoMsg{path: path, err: fmt.Errorf("could not read file: %w", err)}
		}

		mimeType := http.DetectContentType(buffer[:n])

		ext, ok := detect.ExtensionFromMimeType(mimeType)
		if !ok {
			ext = strings.ToLower(filepath.Ext(path))
		}
		ext = alias.Resolve(ext)

		return fileInfoMsg{path: path, mimeType: mimeType, ext: ext, err: nil}
	}
}

// convertCmd executes the conversion and returns a convertDoneMsg.
func convertCmd(srcPath, fromExt, toExt string) tea.Cmd {
	return func() tea.Msg {
		conv, ok := converter.GetConverter(fromExt, toExt)
		if !ok {
			return convertDoneMsg{err: fmt.Errorf("no converter from %s to %s", fromExt, toExt)}
		}

		base := strings.TrimSuffix(srcPath, filepath.Ext(srcPath))
		outputPath := base + toExt

		if err := conv.Convert(srcPath, outputPath, nil); err != nil {
			return convertDoneMsg{err: err}
		}
		return convertDoneMsg{outputPath: outputPath}
	}
}

// Init implements tea.Model.
func (m model) Init() tea.Cmd {
	return tea.Batch(tea.EnableBracketedPaste, m.spinner.Tick)
}

// Update implements tea.Model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyRunes {
			maybePath := sanitizeDroppedPath(string(msg.Runes))
			if filepath.IsAbs(maybePath) {
				m.processing = true
				m.searching = false
				m.file = fileInfo{path: maybePath}
				m.choice = ""
				m.choices = nil
				return m, checkFileCmd(maybePath)
			}
		}

		// If we are in searching mode (typing filename)
		if m.searching && !m.processing && !m.converting {
			prev := m.input.Value()
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			// Update suggestions when input changes
			if m.input.Value() != prev {
				m.suggestions = searchFiles(m.input.Value())
				m.selIdx = 0
			}

			// Additional navigation keys for suggestion list
			switch msg.String() {
			case "up", "k":
				if len(m.suggestions) > 0 && m.selIdx > 0 {
					m.selIdx--
				}
			case "down", "j":
				if len(m.suggestions) > 0 && m.selIdx < len(m.suggestions)-1 {
					m.selIdx++
				}
			case "enter", "tab":
				var chosen string
				if len(m.suggestions) > 0 {
					chosen = m.suggestions[m.selIdx]
				} else {
					chosen = m.input.Value()
				}
				if chosen != "" {
					m.processing = true
					m.searching = false
					m.file = fileInfo{path: chosen}
					return m, tea.Batch(cmd, checkFileCmd(chosen))
				}
			case "esc":
				if m.input.Value() == "" {
					return m, tea.Quit
				}
				m.input.Reset()
				m.suggestions = nil
				m.selIdx = 0
				return m, cmd
			}
			return m, cmd
		}

		// While in choice selector
		if len(m.choices) > 0 && !m.converting && !m.processing {
			return m.handleChoiceKeys(msg)
		}

		// Global keybindings
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}

		return m, nil

	case fileInfoMsg:
		m.processing = false
		m.file = fileInfo(msg)
		if m.file.err == nil {
			// Build choices list dynamically
			convMap := converter.GetConvertersFor(m.file.ext)
			for dest := range convMap {
				m.choices = append(m.choices, dest)
			}
			if len(m.choices) == 0 {
				m.choices = []string{} // no converters, keep empty
			} else {
				m.cursor = 0
			}
		}
		return m, nil

	case convertDoneMsg:
		m.converting = false
		if msg.err != nil {
			m.file.err = msg.err
		} else {
			m.showSuccess = true
			m.successPath = msg.outputPath
			m.file.err = nil
		}

		// Clear chooser state
		m.choices = nil
		m.cursor = 0

		// Schedule reset back to start screen after 2 seconds
		return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg { return resetMsg{} })

	case resetMsg:
		// Reset state to initial search screen
		m.showSuccess = false
		m.file.err = nil
		m.searching = true
		m.input.Reset()
		m.suggestions = searchFiles("")
		m.selIdx = 0
		m.file = fileInfo{}
		m.choice = ""
		m.choices = nil
		return m, nil

	default:
		// Always update spinner
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m *model) handleChoiceKeys(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}
	case "enter":
		m.choice = m.choices[m.cursor]
		m.converting = true
		return *m, convertCmd(m.file.path, m.file.ext, m.choice)
	case "ctrl+c", "esc", "q":
		return *m, tea.Quit
	}
	return *m, nil
}

// View renders the UI.
func (m model) View() string {
	var content string

	if m.processing {
		content = fmt.Sprintf("%s Processing %s…", m.spinner.View(), m.file.path)
	} else if m.converting {
		content = fmt.Sprintf("%s Converting %s → %s…", m.spinner.View(), m.file.ext, m.choice)
	} else if m.showSuccess {
		content = m.styles.InfoBox.Render(m.styles.Success.Render(fmt.Sprintf("✓ Converted to %s", m.successPath)))
	} else if m.file.err != nil && !m.searching {
		content = m.styles.ErrorBox.Render(m.styles.Error.Render(fmt.Sprintf("Error: %v", m.file.err)))
	} else if len(m.choices) > 0 {
		var sb strings.Builder
		sb.WriteString("What format would you like to convert to?\n\n")
		for i, c := range m.choices {
			cursor := " "
			label := strings.TrimPrefix(c, ".")
			if m.cursor == i {
				cursor = ">"
				label = m.styles.Choice.Render(label)
			}
			sb.WriteString(fmt.Sprintf("%s %s\n", cursor, label))
		}
		sb.WriteString("\n" + m.styles.Help.Render("(Enter to select, Esc to quit)"))
		content = m.styles.InfoBox.Render(sb.String())
	} else if m.searching {
		var sb strings.Builder
		sb.WriteString(m.input.View() + "\n")

		for i, s := range m.suggestions {
			cursor := " "
			if m.selIdx == i {
				cursor = ">"
				s = m.styles.Choice.Render(s)
			}
			sb.WriteString(fmt.Sprintf("%s %s\n", cursor, s))
		}

		if len(m.suggestions) == 0 {
			sb.WriteString("No matches found\n")
		}
		// Always show help line
		sb.WriteString("\n" + m.styles.Help.Render("Drop a file here or start typing to search."))

		content = m.styles.InfoBox.Render(sb.String())
	} else if m.file.path == "" {
		content = "Drop a file onto this terminal window or start typing to search."
	} else {
		// Final result
		if m.file.err != nil {
			content = m.styles.ErrorBox.Render(m.styles.Error.Render(fmt.Sprintf("Error: %v", m.file.err)))
		} else {
			info := lipgloss.JoinVertical(lipgloss.Left,
				fmt.Sprintf("File: %s", m.file.path),
				fmt.Sprintf("MIME: %s", m.styles.Success.Render(m.file.mimeType)),
			)
			content = m.styles.InfoBox.Render(info)
		}
	}

	title := m.styles.Title.Render("Convert – Interactive Mode")
	help := m.styles.Help.Render("Press ESC or Ctrl+C to quit.")

	ui := lipgloss.JoinVertical(lipgloss.Center, title, "\n\n", content, "\n\n", help)

	return m.styles.App.Render(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, ui))
}

// Run launches the interactive TUI. Exposed to CLI package.
func Run() error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func searchFiles(query string) []string {
	var matches []string
	lowerQuery := strings.ToLower(query)
	_ = filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}

		// Only filter by query when query is non-empty
		if lowerQuery != "" && !strings.Contains(strings.ToLower(path), lowerQuery) {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		ext = alias.Resolve(ext)
		if convs := converter.GetConvertersFor(ext); len(convs) == 0 {
			return nil // unsupported mime/extension
		}

		matches = append(matches, path)
		if len(matches) >= 10 {
			return filepath.SkipDir
		}
		return nil
	})
	return matches
}

func sanitizeDroppedPath(raw string) string {
	raw = strings.TrimSpace(raw)

	raw = strings.TrimPrefix(raw, "file://")

	// Strip surrounding single or double quotes
	if len(raw) >= 2 {
		if (raw[0] == '\'' && raw[len(raw)-1] == '\'') || (raw[0] == '"' && raw[len(raw)-1] == '"') {
			raw = raw[1 : len(raw)-1]
		}
	}

	// Unescape space characters ("\\ " -> " ")
	raw = strings.ReplaceAll(raw, "\\ ", " ")

	return raw
}
