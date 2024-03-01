package ui

import (
	"bufio"
	"fmt"
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

// Define ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
)

// messageColors maps MessageType values to their corresponding color codes.
var messageColors = map[MessageType]string{
	MessageTypeInfo:    colorReset,
	MessageTypeWarning: colorYellow,
	MessageTypeError:   colorRed,
	MessageTypeNotice:  colorCyan,
	MessageTypeSuccess: colorGreen,
}

func AskYesNo(prompt string, defaultYes bool) bool {
	// default answer should add the correct uppercase to the (Y/n) prompt
	if defaultYes {
		prompt += " (Y/n)"
	} else {
		prompt += " (y/N)"
	}
	fmt.Println(prompt)
	proceed := ""
	_, _ = fmt.Scanln(&proceed) // ignore errors as we only care about the user input
	yesStrings := []string{"Y", "YES", "YEAH", "YEP", "YEA", "YEAH", "YUP"}
	if defaultYes {
		yesStrings = append(yesStrings, "")
	}
	isYes := slices.Contains(yesStrings, strings.TrimSpace(strings.ToUpper(proceed)))
	return isYes
}

// ReadUserInput reads a line of input from the user.
func ReadUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')
	return strings.TrimSpace(userInput)
}

// PrintMessage prints a message to the user.
func PrintMessage(message string, messageType MessageType) {
	color, ok := messageColors[messageType]
	if !ok {
		// If the messageType is not found in the map, use a default color or no color.
		fmt.Print(message)
		return
	}

	// Print the message with color.
	fmt.Printf("%s%s%s", color, message, colorReset)
}
