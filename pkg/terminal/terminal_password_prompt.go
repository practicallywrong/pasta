package terminal

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

// PromptForPassword prompts the user for a password, with password feedback
func PromptForPassword(promptMessage string, mask rune) (string, error) {
	fmt.Print(promptMessage)

	fd := int(os.Stdin.Fd())
	termState, err := term.GetState(fd)
	if err != nil {
		return "", err
	}
	defer term.Restore(fd, termState)

	// Set terminal to raw mode to disable echoing of typed characters
	if _, err := term.MakeRaw(fd); err != nil {
		return "", err
	}

	var password string
	reader := bufio.NewReader(os.Stdin)

	// Read the password one rune at a time
	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			fmt.Println()
			return "", fmt.Errorf("error reading input: %w", err)
		}

		// End input on newline or carriage return (Enter key)
		if char == '\n' || char == '\r' {
			break
		}

		// Handle Ctrl+C (ASCII 3), terminate the program
		if char == 3 { // Ctrl+C
			fmt.Println()
			fmt.Print("\033[2K\033[0G\r")
			return "", fmt.Errorf("operation interrupted by user (Ctrl+C)")
		}

		// Handle backspace (remove last character)
		if char == '\b' || char == '\x7f' {
			if len(password) > 0 {
				password = password[:len(password)-1]
				fmt.Print("\b \b") // Erase the last mask character printed
			}
			continue
		}

		password += string(char)

		fmt.Print(string(mask))
	}

	// Clear the line after input and reset the cursor position
	fmt.Println()
	fmt.Print("\033[2K\033[0G\r")

	return password, nil
}
