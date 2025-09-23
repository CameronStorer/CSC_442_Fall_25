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
	// if EVEN
	if len(stdin)%2 == 0 {
		bitSize = 8
	} else { // if ODD
		bitSize = 7
	}
	// separate the binary string into chucks
	var output string
	for i := bitSize; i < len(stdin); i += bitSize {
		// convert the binary into numbers composed of bits
		number, _ := strconv.ParseInt(stdin[i-bitSize:i], 2, bitSize)
		fmt.Println(number)
		output += string(rune(number))
	}

	fmt.Println(output)

	// handle error
	/*if err != nil {
		// print error message
		fmt.Fprintf(os.Stderr, "Error reading the redirected standard input.\n")
		// kill program
		os.Exit(0)
	}*/
	// fmt.Printf("The binary string '%s' is equal to the number %d\n", stdin, number)
}

// main function
func main() {
	// collect the redirected stdin if there is any...
	stdin := storeStdin()
	// decode the binary
	binaryToASCII(stdin)

}

/*
1001000
1100101
1101100
1101100
1101111
0100000
1010111
1101111
1110010
1101100
1100100
0100001
*/
