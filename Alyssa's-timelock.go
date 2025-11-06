// Name: Alaskan Bullworm 
// Date: 10-27-2025
// Assignment: Timelock
package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"
	"crypto/md5"
)

func main() {
	// getting integers from standard in
	var y, mo, d, hh, mm, ss int
	// scanning for the standard in with error
	_, err := fmt.Fscan(os.Stdin, &y, &mo, &d, &hh, &mm, &ss)
	// checking errors
	if err != nil {
		fmt.Fprintln(os.Stderr, "Provide epoch on stdin: YYYY, MM, DD, HH, mm, SS")
		os.Exit(1)
	}

	// getting the epock time through the systems local time
	loc := time.Local
	epoch := time.Date(y, time.Month(mo), d, hh, mm, ss, 0, loc)

	// the current local system time
	now := time.Now().In(loc)

	// calculating the seconds through signed 64
	elapsed := now.Unix() - epoch.Unix()


		// aligning the beginning interval
		// finding the correct minute/second
		interval := (elapsed / 60 ) * 60

		// if there are decimals??
		intervalStr := fmt.Sprintf("%d", interval)

		// MD5(MD5(interval))
		inner := md5.Sum([]byte(intervalStr))
		innerHex := hex.EncodeToString(inner[:])
		outer := md5.Sum([]byte(innerHex))
		outerHex := hex.EncodeToString(outer[:])

		//Extracting the first two letters
		letters := "" // starts off blank
		for i := 0; i < len(outerHex) && len(letters) < 2; i++ {
			c := outerHex[i]
			if (c >= 'a' && c <= 'f'){
				letters += string(c) // put the first two letters into string
			}
			

		}
	

	// same thing but with last two digits
	digits := "" 
	for i := len(outerHex) - 1; i >= 0 && len(digits) < 2; i-- {
		c := outerHex[i]
		if c >= '0' && c <= '9' {
			digits += string(c)
		}
	}

	// checking if there is not enough numbers or letters
	
	for len(letters) < 2 {
		letters += "?" // using ? to see where something is missing
	}

	// same thing for numbers
	for len(digits) < 2 {
		digits += "?"
	}

	// combining both letters and numbers to print
	code := letters + digits

	// printing it out onto a new line
	io.WriteString(os.Stdout, code+"\n")

}
