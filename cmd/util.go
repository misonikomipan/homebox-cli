package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func confirm(prompt string) bool {
	fmt.Printf("%s [y/N]: ", prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return answer == "y" || answer == "yes"
}

func unmarshalJSON(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
