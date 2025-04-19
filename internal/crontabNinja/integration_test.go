package crontabninja

import (
	"os"
	"testing"
)

func TestIntegrationParseSchedule(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	apiURL := os.Getenv("CRONTAB_NINJA_URL")
	if apiURL == "" {
		apiURL = "https://cronly.app/api/ai/generate" // default URL
	}

	tests := []struct {
		name     string
		schedule string
		want     string
		wantErr  bool
	}{
		{
			name:     "parse daily schedule",
			schedule: "every day at 9am",
			want:     "0 9 * * *",
			wantErr:  false,
		},
		// {
		// 	name:     "parse weekly schedule",
		// 	schedule: "every monday at noon",
		// 	want:     "0 12 * * 1",
		// 	wantErr:  false,
		// },
		{
			name:     "parse invalid schedule",
			schedule: "invalid schedule format xyz",
			want:     "Invalid cron expression.",
			wantErr:  false,
		},
	}

	client := NewClient(apiURL)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.ParseSchedule(tt.schedule)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSchedule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseSchedule() = %v, want %v", got, tt.want)
			}
		})
	}
}
