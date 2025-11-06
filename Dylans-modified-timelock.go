// Name: Alaskan Bullworm
// Date: 10-27-2025
// Assignment: Timelock
package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

func main() {
	var year, month, day, hour, min, sec int
	if _, err := fmt.Fscan(os.Stdin, &year, &month, &day, &hour, &min, &sec); err != nil {
		fmt.Fprintln(os.Stderr, "Usage: provide epoch on stdin as 6 integers: YYYY MM DD HH mm SS")
		os.Exit(1)
	}

	loc := time.Local
	epoch := time.Date(year, time.Month(month), day, hour, min, sec, 0, loc)

	// Continuous loop: recompute, and print only when the minute window changes
	var lastWindow int64 = -1
	for {
		now := time.Now().In(loc)
		elapsed := now.Unix() - epoch.Unix()
		if elapsed < 0 {
			// If the epoch is in the future, wait and try again shortly.
			sleepUntilNextSecond()
			continue
		}

		window := elapsed / 60 // which minute since epoch
		if window != lastWindow {
			code := generateTimelockCode(elapsed)
			fmt.Printf("%s  %s\n", now.Format("2006-01-02 15:04:05"), code)
			lastWindow = window
		}

		// Sleep until the next second to avoid busy-waiting
		sleepUntilNextSecond()
	}
}

// generateTimelockCode returns the 4-character code produced by
// MD5(MD5(intervalSecondsString)), where intervalSeconds is aligned
// to the nearest minute (seconds == 0). The code is two hex letters
// (a-f) taken from the left of the outer hex and two digits taken
// from the right of the outer hex. Missing characters are replaced
// with '?'.
func generateTimelockCode(elapsedSeconds int64) string {
	// align to minute: floor to whole minutes then convert to seconds
	window := elapsedSeconds / 60
	intervalSeconds := window * 60
	intervalStr := fmt.Sprint(intervalSeconds)

	// MD5(MD5(interval))
	inner := md5.Sum([]byte(intervalStr))
	innerHex := hex.EncodeToString(inner[:])
	outer := md5.Sum([]byte(innerHex))
	outerHex := hex.EncodeToString(outer[:]) // lowercase hex string

	// collect first two letters a-f from left to right
	letters := make([]byte, 0, 2)
	for i := 0; i < len(outerHex) && len(letters) < 2; i++ {
		c := outerHex[i]
		if c >= 'a' && c <= 'f' {
			letters = append(letters, c)
		}
	}
	for len(letters) < 2 {
		letters = append(letters, '?')
	}

	// collect last two digits from right to left, then reverse to keep original order
	digitsRev := make([]byte, 0, 2)
	for i := len(outerHex) - 1; i >= 0 && len(digitsRev) < 2; i-- {
		c := outerHex[i]
		if c >= '0' && c <= '9' {
			digitsRev = append(digitsRev, c)
		}
	}
	digits := make([]byte, 0, 2)
	for i := len(digitsRev) - 1; i >= 0; i-- {
		digits = append(digits, digitsRev[i])
	}
	for len(digits) < 2 {
		digits = append(digits, '?')
	}

	return string(letters) + string(digits)
}

// sleepUntilNextSecond pauses until the start of the next whole second.
func sleepUntilNextSecond() {
	now := time.Now()
	next := now.Truncate(time.Second).Add(time.Second)
	time.Sleep(time.Until(next))
}
