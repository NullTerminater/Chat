package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	serverPort    = ":8080"
	discoveryPort = ":9999"
	headerHeight  = 3
	footerHeight  = 3
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

type udpMsg string
type errMsg error

type model struct {
	viewport   viewport.Model
	textInput  textinput.Model
	messages   []string
	conn       *net.UDPConn
	serverAddr *net.UDPAddr
	err        error
	ready      bool
}

func initialModel(conn *net.UDPConn, serverAddr *net.UDPAddr) model {
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput:  ti,
		messages:   []string{},
		conn:       conn,
		serverAddr: serverAddr,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		listenForMessages(m.conn),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			userInput := m.textInput.Value()
			if userInput == "" {
				break
			}

			if m.conn != nil && m.serverAddr != nil {
				_, err := m.conn.Write([]byte(userInput))
				if err != nil {
					m.err = err
				}
			}

			m.messages = append(m.messages, fmt.Sprintf("[%s][You]: %s", time.Now().Format("15:04:05"), userInput))
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textInput.Reset()
			m.viewport.GotoBottom()
		}

	case udpMsg:
		m.messages = append(m.messages, string(msg))
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()
		return m, listenForMessages(m.conn)

	case errMsg:
		m.err = msg
		return m, tea.Quit

	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-headerHeight-footerHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - headerHeight - footerHeight
		}
	}

	m.textInput, tiCmd = m.textInput.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	header := titleStyle.Render("UDP Chat Client")
	footer := m.textInput.View()

	return fmt.Sprintf("%s\n%s\n\n%s", header, m.viewport.View(), footer)
}

func listenForMessages(conn *net.UDPConn) tea.Cmd {
	return func() tea.Msg {
		buf := make([]byte, 2048)
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			return errMsg(err)
		}
		text := string(buf[:n])
		return udpMsg(fmt.Sprintf("[Server]: %s", text))
	}
}

func main() {
	serverIP := findServer()

	serverAddr, err := net.ResolveUDPAddr("udp", serverIP+serverPort)
	if err != nil {
		fmt.Println("Error resolving server address:", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	p := tea.NewProgram(initialModel(conn, serverAddr), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func findServer() string {
	fmt.Println("Listening for server beacon on UDP :9999...")
	addr, err := net.ResolveUDPAddr("udp", discoveryPort)
	if err != nil {
		fmt.Println("Address resolve error:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Listen error (is another client running?):", err)
		os.Exit(1)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	_, remoteAddr, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("Read error:", err)
		os.Exit(1)
	}

	return remoteAddr.IP.String()
}
