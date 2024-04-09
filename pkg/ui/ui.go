package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

type MessageType int

const (
	MessageTypeInfo MessageType = iota
	MessageTypeWarning
	MessageTypeError
	MessageTypeNotice
	MessageTypeSuccess
)

// Define ANSI color codes.
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
)

type UI struct {
	Reader io.Reader
	Writer io.Writer
}

type opt func(*UI)

func WithReader(reader io.Reader) opt {
	return func(ui *UI) {
		ui.Reader = reader
	}
}

func WithWriter(writer io.Writer) opt {
	return func(ui *UI) {
		ui.Writer = writer
	}
}

func NewUI(opts ...opt) UI {
	ui := UI{ //nolint:varnamelen
		Reader: os.Stdin,
		Writer: os.Stdout,
	}

	for _, opt := range opts {
		opt(&ui)
	}

	return ui
}

func (u UI) AskYesNo(prompt string, defaultYes bool) bool {
	// default answer should add the correct uppercase to the (Y/n) prompt
	if defaultYes {
		prompt += " (Y/n)"
	} else {
		prompt += " (y/N)"
	}

	_, err := u.Writer.Write([]byte(prompt + "\n"))
	if err != nil {
		return false
	}

	yesStrings := []string{"Y", "YES", "YEAH", "YEP", "YEA", "YEAH", "YUP"}

	scanner := bufio.NewScanner(u.Reader)
	for scanner.Scan() {
		answer := strings.ToUpper(scanner.Text())
		if answer == "" {
			return defaultYes
		}

		if slices.Contains(yesStrings, answer) {
			return true
		}
	}

	return false
}

// ReadUserInput reads a line of input from the user.
func (u UI) ReadUserInput() string {
	scanner := bufio.NewScanner(u.Reader)
	scanner.Scan()
	userInput := scanner.Text()

	return strings.TrimSpace(userInput)
}

// PrintMessage prints a message to the user.
func (u UI) PrintMessage(message string, messageType MessageType) {
	if messageType == MessageTypeInfo {
		fmt.Fprint(u.Writer, message)
		return
	}

	messageColors := map[MessageType]string{
		MessageTypeWarning: colorYellow,
		MessageTypeError:   colorRed,
		MessageTypeNotice:  colorCyan,
		MessageTypeSuccess: colorGreen,
	}

	color, ok := messageColors[messageType]
	if !ok {
		// If the messageType is not found in the map, use a default color or no color.
		fmt.Fprint(u.Writer, message)
		return
	}

	// Print the message with color.
	fmt.Fprint(u.Writer, color+message+colorReset)
}
