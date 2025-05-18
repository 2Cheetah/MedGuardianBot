package texttotime

import (
	"testing"
	"time"
)

type LLMClientStub struct{}

func (llmC *LLMClientStub) ParseTextToISO8601(text string) (string, error) {
	switch text {
	case "until 10th Sep 2026":
		return "2026-09-10", nil
	case "tomorrow":
		return time.Now().Add(24 * time.Hour).Format(time.DateOnly), nil
	case "foo":
		return "bar", nil
	}
	return "", nil
}

func TestParseText(t *testing.T) {
	p := Parser{&LLMClientStub{}}
	tests := []struct {
		name           string
		input          string
		expectedOutput time.Time
		hasError       bool
	}{
		{
			name:           "simple date",
			input:          "until 10th Sep 2026",
			expectedOutput: time.Date(2026, time.September, 10, 0, 0, 0, 0, time.UTC),
			hasError:       false,
		},
		{
			name:           "relative date",
			input:          "tomorrow",
			expectedOutput: time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour),
			hasError:       false,
		},
		{
			name:           "wrong date",
			input:          "foo",
			expectedOutput: time.Time{},
			hasError:       true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dateISO8601, err := p.ParseText(test.input)
			if err != nil {
				if test.hasError == false {
					t.Errorf("error %v not expected", err)
				}
			}
			if !dateISO8601.Equal(test.expectedOutput) {
				t.Errorf("expected %s, got %s", test.expectedOutput, dateISO8601)
			}
		})
	}
}
