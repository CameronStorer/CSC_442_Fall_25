package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

//import "bufio"

// helper method to store the standard input
// thank you cameron !!!
func storeStdin() string {
	// first, ensure that standard input has been redirected
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Println("Expected stdin !! Use ./binary < something.txt")
		os.Exit(0)
	}
	// attempt to read the stdin by storing it and error checking
	input, err := io.ReadAll(os.Stdin)
	// if there is an error while reading stdin
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading the redirected standard input.\n")
		// kill program
		os.Exit(0)
	}

	return string(input)
}

func main() {

	raw := storeStdin()

	// detect 7-bit or 8-bit
	var kind int = 7
	if (len(raw)-1)%8 == 0 {
		kind = 8
	}

	// split into groups of 7/8 int64s; stores decimal values in letters
	letters := []int64{}
	for i := 0; i < len(raw)-1; i += kind {
		end := i + kind
		if end > len(raw)-1 {
			end = len(raw) - 1
		}
		tempe, err := strconv.ParseInt(raw[i:end], 2, 0)
		if err != nil {
			fmt.Printf("error during parseint, %v", err)
		}
		letters = append(letters, tempe)

	}
	// debug stuff
	//fmt.Println("")
	//fmt.Println((len(raw) - 1) % 8)
	//fmt.Println(kind)
	//fmt.Println(letters)
	//out_int := []int{}

	// output string slices
	out := []string{}

	// convert every int64 type into an int, then into a rune, then that into a string
	for i := range letters {
		out = append(out, string(rune(int(letters[i]))))
	}

	// finally print out every individual string
	for i := range out {
		fmt.Printf("%v", out[i])
	}
}
