package texttocron

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseSchedule(t *testing.T) {
	tests := []struct {
		name           string
		schedule       string
		serverResponse string
		statusCode     int
		want           string
		wantErr        bool
	}{
		{
			name:           "successful parsing",
			schedule:       "every day at 3pm",
			serverResponse: `{"crontab": "0 15 * * *"}`,
			statusCode:     http.StatusOK,
			want:           "0 15 * * *",
			wantErr:        false,
		},
		{
			name:           "server error",
			schedule:       "invalid schedule",
			serverResponse: `{"error": "Invalid schedule"}`,
			statusCode:     http.StatusBadRequest,
			want:           "",
			wantErr:        true,
		},
		{
			name:           "invalid json response",
			schedule:       "every minute",
			serverResponse: `invalid json`,
			statusCode:     http.StatusOK,
			want:           "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			client := NewClient(server.URL)
			got, err := client.ParseSchedule(tt.schedule)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSchedule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseSchedule() = %v, want %v", got, tt.want)
			}
		})
	}
}
