package validation

import (
	"fmt"
	"regexp"
	"strings"
)

type TextValidation struct {
	emptyCnt uint32
	totalCnt uint32
	viaCnt uint32
}

func(validation *TextValidation) Validate(text string) (error string) {
	if len(text) > 300 {
		return fmt.Sprintf("🚫 Text is too long, 300 characters max (you have %d)", len(text))
	}

	if len(text) < 20 {
		return fmt.Sprintf("🚫 Text is too short, 20 characters min (you have %d)", len(text))
	}

	if !strings.HasPrefix(text, "🚗") && !strings.HasPrefix(text, "👋") && !strings.HasPrefix(text, "🚘") {
		return "🚫 Text must start with 🚗 or 👋"
	}

	if strings.Contains(text, "@") {
		return "🚫 No @ symbol please"
	}

	lowerText := strings.ToLower(text)

	if strings.Contains(lowerText, "fuck") {
		return "🚫 No English bad words please 🤦"
	}

	if strings.Contains(lowerText, "хуй") || strings.Contains(lowerText, "пизда") || strings.Contains(lowerText, "блядь") {
		return "🚫 No Russian bad words please 🤦"
	}

	if strings.Contains(lowerText, "foobar") {
		return "🚫 Whoa, don't just copy and paste text, change \"Pick up\" and \"Drop off\" at least"
	}

	if strings.Contains(lowerText, "http") || strings.Contains(lowerText, "ftp://") || strings.Contains(lowerText,"www") {
		return "🚫 No links please"
	}

	var re = regexp.MustCompile(`\w\.\w{2,5}`)

	if len(re.FindAllString(lowerText, -1)) > 0 {
		return "🚫 No domains please"
	}

	lines := strings.Split(text, "\n")

	for i := range lines {
		line := strings.TrimSpace(lines[i])

		if len(line) == 0 {
			validation.emptyCnt++
		} else {
			if strings.HasPrefix(line, "via ") {
				validation.viaCnt++
			}
			validation.totalCnt++
		}
	}

	if validation.emptyCnt > 2 {
		return fmt.Sprintf("🚫 Only 2 empty lines max please (you have %d)", validation.emptyCnt)
	}

	if validation.totalCnt > 9 {
		return fmt.Sprintf("🚫 Only 9 max lines please (you have %d)", validation.totalCnt)
	}

	if validation.totalCnt < 5 {
		return fmt.Sprintf("🚫 At least 5 lines please (you have %d)", validation.totalCnt)
	}

	if validation.viaCnt > 0 {
		return "🚫 Do not include \"via\", it will be added automatically"
	}

	return ""
}

func NewTextValidation() (textValidation *TextValidation) {
	validation := &TextValidation{}
	return validation
}