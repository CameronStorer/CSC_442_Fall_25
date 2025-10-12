/*	Author: AlaskanBullWorm
	Date: 10/10/25
	Desc: simple go program to decode covert timing channel messages
*/

package main

// import required libraries
import (
	"errors"  // for handling errors
	"fmt"     // for printing
	"io"      // input output
	"net"     // for connecting to chat server
	"os"      // for making system calls
	"strconv" // for to and from string conversions
	"time"    // for reading time
)

// define some constants
const (
	address  = "138.47.99.228" // chat server address
	port     = "31337"         // chat server port
	msForBit = 65000000        // ms value that will determine when to call a time 1 or 0
)

// helper method to convert binary strings into ASCII
// (this method was recycled from a previous assignment)
func bin2ASCII(binaryStringOut string) string {
	// now we need to convert the binary to ASCII
	output := ""

	// ensure that the binary string is % 8 by padding so that
	// there are no slicing errors/
	if len(binaryStringOut)%8 != 0 {
		for {
			binaryStringOut += "0"
			if len(binaryStringOut)%8 == 0 {
				break
			}
		}
	}

	// for each 8 bits of the string:
	for i := 0; i < len(binaryStringOut); i += 8 {

		// store the bits
		tempEight := binaryStringOut[i : i+8]

		// convert these eight bits into an int
		decValue, err := strconv.ParseUint(tempEight, 2, 8)
		// error checking
		if err != nil {
			fmt.Println("Binary to ASCII conversion error: ", err)
		}
		output += string(byte(decValue))
	}
	return output
}

// main runtime
func main() {
	// connect to the specified chat server using tcp protocol
	connection, err := net.Dial("tcp", address+":"+port)
	fmt.Println("[connect to the chat server]")
	// handle errors
	if err != nil {
		fmt.Println("Connection error: ", err)
		os.Exit(1)
	}

	// now start recieving messages
	// fmt.Println("Recieving message: ")
	// define a start time to compare to time of first letter
	prevChar := time.Now()
	// declare a variable to hold bits
	var currentBit int
	// declare a string to record the converted bit values
	var bitValues string
	// read response from the server
	for {
		// set the read size to be 1 byte
		buffer := make([]byte, 1)
		// attempt to read a byte
		n, err := connection.Read(buffer)
		// handle errors
		if err != nil {
			// if the chat server indicates "end of file"
			if errors.Is(err, io.EOF) {
				// fmt.Println("\nMessage fully recieved.")
				break
			}
			fmt.Println("Reading error: ", err)
			return
		}
		// print recieved message as it arrives
		if n > 0 {
			// convert the elapsed time since previous character to a bit
			// value
			if time.Since(prevChar).Nanoseconds() <= msForBit {
				currentBit = 0
			} else {
				currentBit = 1
			}
			bitValues += strconv.Itoa(currentBit)
			fmt.Print(currentBit)
			// fmt.Println("\t", time.Since(prevChar).Nanoseconds())
			// fmt.Printf(string(buffer[:n]))
			prevChar = time.Now()
		}
	}
	// ensure the connection closes when no longer needed
	connection.Close()
	fmt.Println("\n[disconnect from the chat server]")
	// slice off first useless bit
	bitValues = bitValues[1:]
	// call method to convert the binary string into ASCII
	fmt.Println("Covert message: ", bin2ASCII(bitValues))

}
