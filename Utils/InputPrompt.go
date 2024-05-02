package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func InputPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func InputPromptInt(label string) int {
	var i int
	for {
		fmt.Print(label + " ")
		_, err := fmt.Scanf("%d", &i)
		if err == nil {
			break
		}
	}
	return i
}
