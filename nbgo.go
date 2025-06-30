package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	listStyle   = lipgloss.NewStyle().Margin(1, 2)
	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true).Padding(0, 1)
)

type item struct {
	path, title string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.path }
func (i item) FilterValue() string { return i.title }

type model struct {
	list           list.Model
	input          textinput.Model
	currentDir     string
	addingNote     bool
	addingBookmark bool
	err            error
}

func InitialModel() model {
	homeDir, _ := os.UserHomeDir()
	defaultDir := filepath.Join(homeDir, ".nbgo", "default")
	items := loadNotes(defaultDir)

	listDelegate := list.NewDefaultDelegate()
	listDelegate.Styles.SelectedTitle = listDelegate.Styles.SelectedTitle.Foreground(lipgloss.Color("39")).Bold(true)
	listDelegate.Styles.SelectedDesc = listDelegate.Styles.SelectedDesc.Foreground(lipgloss.Color("242"))

	l := list.New(items, listDelegate, 80, 20)
	l.Title = "Notes in " + defaultDir
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	ti := textinput.New()
	ti.Placeholder = "Enter note title..."
	ti.CharLimit = 80
	ti.Width = 50

	return model{
		list:           l,
		input:          ti,
		currentDir:     defaultDir,
		addingNote:     false,
		addingBookmark: false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func loadNotes(dir string) []list.Item {
	var items []list.Item
	files, err := os.ReadDir(dir)
	if err != nil {
		return items
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") || strings.HasSuffix(file.Name(), ".bookmark.md") {
			name := strings.TrimSuffix(file.Name(), ".md")
			name = strings.TrimSuffix(name, ".bookmark")
			items = append(items, item{path: filepath.Join(dir, file.Name()), title: name})
		}
	}
	return items
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.addingNote || m.addingBookmark {
			switch msg.String() {
			case "esc":
				m.addingNote = false
				m.addingBookmark = false
				m.input.Reset()
				return m, nil
			case "enter":
				if m.addingNote {
					title := m.input.Value()
					if title != "" {
						filename := fmt.Sprintf("%s.md", time.Now().Format("20060102150405"))
						filePath := filepath.Join(m.currentDir, filename)
						if err := os.MkdirAll(m.currentDir, 0755); err != nil {
							m.err = err
							return m, nil
						}
						if err := os.WriteFile(filePath, []byte("# "+title+"\n"), 0644); err != nil {
							m.err = err
							return m, nil
						}
						cmd := exec.Command("84", filePath)
						cmd.Stdin = os.Stdin
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						if err := cmd.Run(); err != nil {
							m.err = err
							return m, nil
						}
						m.list.SetItems(loadNotes(m.currentDir))
					}
				} else if m.addingBookmark {
					url := m.input.Value()
					if url != "" {
						filename := fmt.Sprintf("%s.bookmark.md", time.Now().Format("20060102150405"))
						filePath := filepath.Join(m.currentDir, filename)
						if err := os.MkdirAll(m.currentDir, 0755); err != nil {
							m.err = err
							return m, nil
						}
						if err := os.WriteFile(filePath, []byte(url+"\n"), 0644); err != nil {
							m.err = err
							return m, nil
						}
						cmd := exec.Command("84", filePath)
						cmd.Stdin = os.Stdin
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						if err := cmd.Run(); err != nil {
							m.err = err
							return m, nil
						}
						m.list.SetItems(loadNotes(m.currentDir))
					}
				}
				m.addingNote = false
				m.addingBookmark = false
				m.input.Reset()
				return m, nil
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "a":
			m.addingNote = true
			m.addingBookmark = false
			m.input.Placeholder = "Enter note title..."
			m.input.Focus()
			return m, textinput.Blink
		case "b":
			m.addingBookmark = true
			m.addingNote = false
			m.input.Placeholder = "Enter bookmark URL..."
			m.input.Focus()
			return m, textinput.Blink
		case "e":
			if selected, ok := m.list.SelectedItem().(item); ok {
				cmd := exec.Command("84", selected.path)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					m.err = err
					return m, nil
				}
				m.list.SetItems(loadNotes(m.currentDir))
			}
			return m, nil
		case "v":
			if selected, ok := m.list.SelectedItem().(item); ok {
				cmd := exec.Command("glow", selected.path)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					m.err = err
					return m, nil
				}
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("%s\n\n%s", headerStyle.Render("Error"), m.err.Error())
	}

	if m.addingNote || m.addingBookmark {
		return fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			headerStyle.Render("Add Note"),
			m.input.View(),
			"Press Enter to save, Esc to cancel",
		)
	}

	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		headerStyle.Render("Notes in "+m.currentDir),
		m.list.View(),
		"Press a to add note, b for bookmark, e to edit, v to view, q to quit",
	)
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "use" {
		homeDir, _ := os.UserHomeDir()
		newDir := filepath.Join(homeDir, ".nbgo", args[1])
		if err := os.MkdirAll(newDir, 0755); err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			os.Exit(1)
		}
		if err := os.WriteFile(filepath.Join(homeDir, ".nbgo", ".current"), []byte(newDir), 0644); err != nil {
			fmt.Printf("Error setting current notebook: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Switched to notebook: %s\n", args[1])
		os.Exit(0)
	}

	homeDir, _ := os.UserHomeDir()
	currentDir := filepath.Join(homeDir, ".nbgo", "default")
	if current, err := os.ReadFile(filepath.Join(homeDir, ".nbgo", ".current")); err == nil {
		currentDir = string(current)
	}

	m := InitialModel()
	m.currentDir = currentDir
	m.list.SetItems(loadNotes(currentDir))
	m.list.Title = "Notes in " + currentDir

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
