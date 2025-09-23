// Description: Binary Decoder
// Cyberstorm Group Name: AlaskanBullWorm

package main

// import appropriate libraries
import (
	"fmt"     // library used for printing
	"io"      // libarary used for taking stdin
	"os"      // libary used for making system calls
	"strconv" // library used for converting strings
)

// helper method to store the standard input
func storeStdin() string {
	// first, ensure that standard input has been redirected
	stat, _ := os.Stdin.Stat()
	// if not ask the user to do so next time
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Println("Please redirect standard input.")
		fmt.Println("EX: './binary < example.txt'")
		// kill program to allow user to retry
		os.Exit(0)
	}
	// attempt to read the stdin by storing it and error checking
	input, err := io.ReadAll(os.Stdin)
	// if there is an error while reading stdin
	if err != nil {
		// print error message
		fmt.Fprintf(os.Stderr, "Error reading the redirected standard input.\n")
		// kill program
		os.Exit(0)
	}
	// convert the input into string and return
	return string(input)
}

// helper method to decode the binary
func binaryToASCII(stdin string) {
	// determine if the number of bits is even or odd
	bitSize := 0
	// if divisible by 8
	if len(stdin[:len(stdin)-1])%8 == 0 {
		bitSize = 8
	} else { // otherwise
		bitSize = 7
	}
	// separate the binary string into chucks
	var output string
	for i := bitSize; i < len(stdin); i += bitSize {
		// convert the binary into numbers composed of bits
		number, _ := strconv.ParseInt(stdin[i-bitSize:i], 2, bitSize+1)
		output += string(rune(number))
		// fmt.Println(stdin[i-bitSize:i])
		// fmt.Println(number)
	}
	// fmt.Println(strconv.ParseInt("1001000", 2, bitSize))
	// print out the output
	fmt.Println(output)
}

// main function
func main() {
	// collect the redirected stdin if there is any...
	stdin := storeStdin()
	// decode the binary
	binaryToASCII(stdin)

}

