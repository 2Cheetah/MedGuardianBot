package texttodate

import (
	"fmt"
	"time"

	"github.com/2Cheetah/MedGuardianBot/internal/groq"
)

type Parser struct {
	llm *groq.Client
}

func NewParser(apiKey string) *Parser {
	llm := groq.NewClient(apiKey)
	return &Parser{llm: llm}
}

func (p *Parser) ParseText(text string) (time.Time, error) {
	timestampString, err := p.llm.ParseTextToISO8601(text)
	if err != nil {
		return time.Time{}, fmt.Errorf("couldn't parse text to ISO 8601 string, text: %s, error: %w", text, err)
	}
	dateISO8601, err := time.Parse(time.DateOnly, timestampString)
	if err != nil {
		return time.Time{}, fmt.Errorf("couldn't parse ISO 8601 string to time.Time representation, string: %s, error: %w", timestampString, err)
	}
	return dateISO8601, nil
}
