package tracktime

import (
	"testing"
	"time"
)

func TestTrackTime(t *testing.T) {
	defer TrackTime(time.Now())

	time.Sleep(1 * time.Second)
}
