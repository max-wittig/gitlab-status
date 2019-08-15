package main

import (
	"testing"
	"time"
)

func TestFindCurrentShouldStatus(t *testing.T) {
	crons := make(map[string]Status)
	crons["0 7 * * *"] = Status{Emoji: "7", Message: "7"}
	crons["0 10 * * *"] = Status{Emoji: "10", Message: "10"}
	crons["0 12 * * *"] = Status{Emoji: "12", Message: "12"}
	crons["0 13 * * *"] = Status{Emoji: "13", Message: "13"}
	crons["0 15 * * *"] = Status{Emoji: "15", Message: "15"}
	crons["0 18 * * *"] = Status{Emoji: "18", Message: "18"}
	cronsWithOnlyOne := make(map[string]Status)
	cronsWithOnlyOne["0 12 * * *"] = Status{Emoji: "12", Message: "12"}
	tests := []struct {
		now        time.Time
		wantStatus Status
		name       string
		crons      map[string]Status
	}{
		{
			now:        time.Date(2019, 8, 13, 10, 30, 2, 0, time.UTC),
			wantStatus: crons["0 10 * * *"],
			crons:      crons,
			name:       "test1",
		},
		{
			now:        time.Date(2019, 8, 13, 10, 00, 00, 0, time.UTC),
			wantStatus: crons["0 10 * * *"],
			crons:      crons,
			name:       "test2",
		},
		{
			now:        time.Date(2019, 8, 13, 12, 30, 0, 5, time.UTC),
			wantStatus: crons["0 12 * * *"],
			crons:      crons,
			name:       "test3",
		},
		{
			now:        time.Date(2019, 8, 13, 14, 30, 6, 2, time.UTC),
			wantStatus: crons["0 13 * * *"],
			crons:      crons,
			name:       "test4",
		},
		{
			now:        time.Date(2019, 8, 13, 20, 00, 5, 1, time.UTC),
			wantStatus: crons["0 18 * * *"],
			crons:      crons,
			name:       "test5",
		},
		{
			now:        time.Date(2019, 8, 13, 6, 30, 25, 0, time.UTC),
			wantStatus: crons["0 18 * * *"],
			crons:      crons,
			name:       "test6",
		},
		{
			now:        time.Date(2019, 8, 13, 8, 20, 2, 0, time.UTC),
			wantStatus: crons["0 7 * * *"],
			crons:      crons,
			name:       "test7",
		},
		{
			now:        time.Date(2019, 8, 13, 8, 20, 2, 0, time.UTC),
			wantStatus: cronsWithOnlyOne["0 12 * * *"],
			crons:      cronsWithOnlyOne,
			name:       "test8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{Timezone: "Europe/Zurich", Crons: tt.crons}
			result, err := findCurrentShouldStatus(&config, tt.now)
			if err != nil {
				t.Errorf("Error, while processing %s", tt.name)
			}
			if *result != tt.wantStatus {
				t.Errorf("Invalid status for %s found. Should be %s. Is %s", tt.now, tt.wantStatus, *result)
			}
		})
	}
}
