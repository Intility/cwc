package ui_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/intility/cwc/pkg/ui"
)

func TestAskYesNo(t *testing.T) {
	testCases := []struct {
		name       string
		prompt     string
		userInput  string
		want       bool
		defaultYes bool
	}{
		{
			name:       "Yes, default",
			prompt:     "Do you want to proceed?",
			userInput:  "y\n",
			defaultYes: true,
			want:       true,
		},
		{
			name:       "No, default",
			prompt:     "Do you want to proceed?",
			userInput:  "n\n",
			defaultYes: false,
			want:       false,
		},
		{
			name:       "Yes, default, uppercase",
			prompt:     "Do you want to proceed?",
			userInput:  "Y\n",
			defaultYes: true,
			want:       true,
		},
		{
			name:       "No, default, uppercase",
			prompt:     "Do you want to proceed?",
			userInput:  "N\n",
			defaultYes: false,
			want:       false,
		},
		{
			name:       "Yes, not default",
			prompt:     "Do you want to proceed?",
			userInput:  "y\n",
			defaultYes: false,
			want:       true,
		},
		{
			name:       "No, not default",
			prompt:     "Do you want to proceed?",
			userInput:  "n\n",
			defaultYes: false,
			want:       false,
		},
		{
			name:       "Random input",
			prompt:     "Do you want to proceed?",
			userInput:  "kochicomputeren\n",
			defaultYes: false,
			want:       false,
		},
		{
			name:       "Empty input",
			prompt:     "Do you want to proceed?",
			userInput:  "\n",
			defaultYes: false,
			want:       false,
		},
		{
			name:       "Empty input, default yes",
			prompt:     "Do you want to proceed?",
			userInput:  "\n",
			defaultYes: true,
			want:       true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := bytes.NewBufferString(tc.userInput)
			buf := &bytes.Buffer{}
			ui := ui.NewUI(ui.WithReader(reader), ui.WithWriter(buf))

			t.Run("Return value is correct", func(t *testing.T) {
				got := ui.AskYesNo(tc.prompt, tc.defaultYes)
				if got != tc.want {
					t.Errorf(cmp.Diff(got, tc.want))
				}
			})

			t.Run("Prompt is correct", func(t *testing.T) {
				if tc.defaultYes {
					tc.prompt += " (Y/n)"
				} else {
					tc.prompt += " (y/N)"
				}
				wantPrompt := fmt.Sprintf("%s\n", tc.prompt)
				if buf.String() != wantPrompt {
					t.Errorf(cmp.Diff(wantPrompt, buf.String()))
				}
			})
		})
	}
}

func TestReadUserInput(t *testing.T) {
	testCases := []struct {
		name      string
		userInput string
		want      string
	}{
		{
			name:      "Simple input",
			userInput: "Hello, world!\n",
			want:      "Hello, world!",
		},
		{
			name:      "Empty input",
			userInput: "\n",
			want:      "",
		},
		{
			name:      "Input with newline",
			userInput: "Hello, world!\n",
			want:      "Hello, world!",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := bytes.NewBufferString(tc.userInput)
			buf := &bytes.Buffer{}
			ui := ui.NewUI(ui.WithReader(reader), ui.WithWriter(buf))
			t.Run("Return value is correct", func(t *testing.T) {
				got := ui.ReadUserInput()
				if got != tc.want {
					t.Errorf(cmp.Diff(got, tc.want))
				}
			})
		})
	}
}

func TestPrintMessage(t *testing.T) {
	testCases := []struct {
		name        string
		message     string
		messageType ui.MessageType
		want        string
	}{
		{
			name:        "Info message",
			message:     "Just do it!",
			messageType: ui.MessageTypeInfo,
			want:        "Just do it!",
		},
		{
			name:        "Warning message",
			message:     "Just do it!",
			messageType: ui.MessageTypeWarning,
			want:        "\x1b[33mJust do it!\x1b[0m",
		},
		{
			name:        "Error message",
			message:     "Just do it!",
			messageType: ui.MessageTypeError,
			want:        "\x1b[31mJust do it!\x1b[0m",
		},
		{
			name:        "Success message",
			message:     "Just do it!",
			messageType: ui.MessageTypeSuccess,
			want:        "\x1b[32mJust do it!\x1b[0m",
		},
		{
			name:        "Notice message",
			message:     "Just do it!",
			messageType: ui.MessageTypeNotice,
			want:        "\x1b[36mJust do it!\x1b[0m",
		},
		{
			name:        "Unknown message type",
			message:     "Just do it!",
			messageType: ui.MessageType(42),
			want:        "Just do it!",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := bytes.NewBufferString("")
			buf := &bytes.Buffer{}
			ui := ui.NewUI(ui.WithReader(reader), ui.WithWriter(buf))

			ui.PrintMessage(tc.message, tc.messageType)
			if buf.String() != tc.want {
				t.Errorf(cmp.Diff(buf.String(), tc.want))
			}
		})
	}
}
