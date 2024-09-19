package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCalculateNextRunTime(t *testing.T) {
	tests := []struct {
		name            string
		currentTime     time.Time
		runTime         time.Time
		expectedNextRun time.Time
	}{
		{
			name:            "Next run is today",
			currentTime:     time.Date(2024, 9, 18, 8, 0, 0, 0, time.UTC),
			runTime:         time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
			expectedNextRun: time.Date(2024, 9, 18, 9, 0, 0, 0, time.UTC),
		},
		{
			name:            "Next run is tomorrow",
			currentTime:     time.Date(2024, 9, 18, 10, 0, 0, 0, time.UTC),
			runTime:         time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
			expectedNextRun: time.Date(2024, 9, 19, 9, 0, 0, 0, time.UTC),
		},
		{
			name:            "Run at midnight",
			currentTime:     time.Date(2024, 9, 18, 23, 59, 0, 0, time.UTC),
			runTime:         time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedNextRun: time.Date(2024, 9, 19, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worker := &APODWorker{
				RunTime: tt.runTime,
			}

			nextRun := worker.calculateNextRunTime(tt.currentTime)
			assert.Equal(t, tt.expectedNextRun, nextRun)
		})
	}
}
