package tracktime

import (
	"fmt"
	"time"
)

// TrackTime 計算並打印兩個時間點之間的持續時間
func TrackTime(pre time.Time) time.Duration {
	expend := time.Since(pre)
	fmt.Println("expend : ", expend)
	return expend
}
