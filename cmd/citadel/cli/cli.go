package cli

import (
	"fmt"

	"github.com/sveltinio/prompti/toggle"
)

func AskYesOrNo(question string) bool {
	result, _ := toggle.Run(&toggle.Config{Question: question})

	return result
}

func Ask(question string) string {
	fmt.Println(question)

	var answer string
	fmt.Scanln(&answer)

	return answer
}

func isLowercaseAlphaNumericString(s string) bool {
	for _, c := range s {
		if !isLowercaseAlphaNumericRune(c) {
			return false
		}
	}

	return true
}

func isLowercaseAlphaNumericRune(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
}
