package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	headerHeight = 3
	footerHeight = 3
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A0A0A0")).
			Italic(true)
)

type UdpMsg string
type errMsg error

type SendFunc func(msg string) error

type ChatModel struct {
	Viewport     viewport.Model
	TextInput    textinput.Model
	Messages     []string
	SendCallback SendFunc
	Title        string
	Err          error
	Ready        bool
}

func InitialModel(title string, sendCallback SendFunc) ChatModel {
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return ChatModel{
		TextInput:    ti,
		Messages:     []string{},
		SendCallback: sendCallback,
		Title:        title,
	}
}

func (m ChatModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			userInput := m.TextInput.Value()
			if userInput == "" {
				break
			}

			if m.SendCallback != nil {
				err := m.SendCallback(userInput)
				if err != nil {
					m.Err = err
				}
			}

			m.Messages = append(m.Messages, fmt.Sprintf("[%s][You]: %s", time.Now().Format("15:04:05"), userInput))
			m.Viewport.SetContent(strings.Join(m.Messages, "\n"))
			m.TextInput.Reset()
			m.Viewport.GotoBottom()
		}

	case UdpMsg:
		m.Messages = append(m.Messages, string(msg))
		m.Viewport.SetContent(strings.Join(m.Messages, "\n"))
		m.Viewport.GotoBottom()

	case errMsg:
		m.Err = msg
		return m, tea.Quit

	case tea.WindowSizeMsg:
		if !m.Ready {
			m.Viewport = viewport.New(msg.Width, msg.Height-headerHeight-footerHeight)
			m.Viewport.YPosition = headerHeight
			m.Viewport.SetContent(strings.Join(m.Messages, "\n"))
			m.Ready = true
		} else {
			m.Viewport.Width = msg.Width
			m.Viewport.Height = msg.Height - headerHeight - footerHeight
		}
	}

	m.TextInput, tiCmd = m.TextInput.Update(msg)
	m.Viewport, vpCmd = m.Viewport.Update(msg)

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m ChatModel) View() string {
	if !m.Ready {
		return "\n  Initializing..."
	}

	header := titleStyle.Render(m.Title)
	footer := m.TextInput.View()

	return fmt.Sprintf("%s\n%s\n\n%s", header, m.Viewport.View(), footer)
}
