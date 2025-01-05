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
			return "", fmt.Errorf("error reading input: %w", err)
		}

		// End input on newline or carriage return
		if char == '\n' || char == '\r' {
			break
		}

		// Append the character to the password
		password += string(char)

		// Optionally mask characters with '*' if the mask flag is set
		fmt.Print(mask)
	}

	// Clear the line after input
	fmt.Println()
	fmt.Print("\033[2K\033[0G\r")

	return password, nil
}
