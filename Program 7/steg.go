/*
Name: AlaskanBullWorm
Date: 10/27/25
Desc: simple Steg program
*/

package main

// IMPORT LIBRARIES
import (
	"flag" // for handling command-line arguments
	"fmt"  // printing
	"os"   // for making system calls and manual argument parsing
	"strconv" // convert strings to int
	"strings" // string manipulation
)

// DEFINE CONSTANTS
const (
	READDIRECTION = "L2R" // defines the bit processing order
	SENTINELSIZE  = 6
)

// DECLARE FLAG VALUE VARIABLES
var (
	// sentinel
	SENTINEL = []byte{0x00, 0xff, 0x00, 0x00, 0xff, 0x00}

	// store or retrieve
	store    bool
	retrieve bool

	// bit or byte
	useBits  bool
	useBytes bool

	// the remaining parameters
	offset              int
	interval            int
	wrapperFileLocation string
	hiddenFileLocation  string
)

// INITIALIZE FLAG VARIABLES TO THEIR DEFAULT VALUES BEFORE MAIN FUNC
func init() {
	// IMPORTANT: We define all flags here so they appear in the Usage message.
	// The value-based flags will be parsed manually in parseManualArgs().

	// store or retrieve
	flag.BoolVar(&store, "s", false, "store hidden data")
	flag.BoolVar(&retrieve, "r", false, "read hidden data")

	// bit or byte
	flag.BoolVar(&useBits, "b", false, "bit method")
	flag.BoolVar(&useBytes, "B", false, "byte method")

	// the remaining parameters (defined for Usage display only)
	flag.IntVar(&offset, "o", 0, "offset value default 0")
	flag.IntVar(&interval, "i", 1, "interval value default 1")
	flag.StringVar(&wrapperFileLocation, "w", "", "wrapper file location")
	flag.StringVar(&hiddenFileLocation, "h", "", "hidden file location")
}

// Manually parses combined flags (e.g., -o1024)
func parseManualArgs() {
	// set defaults manually
	offset = 0
	interval = 1

	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-o") {
			// handle -o flag
			valStr := strings.TrimPrefix(arg, "-o")
			if valStr == "" {
				fmt.Fprintln(os.Stderr, "Error: -o flag requires a value.")
				os.Exit(1)
			}
			val, err := strconv.Atoi(valStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Invalid integer for -o: %s\n", valStr)
				os.Exit(1)
			}
			offset = val

		} else if strings.HasPrefix(arg, "-i") {
			// handle -i flag
			valStr := strings.TrimPrefix(arg, "-i")
			if valStr == "" {
				fmt.Fprintln(os.Stderr, "Error: -i flag requires a value.")
				os.Exit(1)
			}
			val, err := strconv.Atoi(valStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Invalid integer for -i: %s\n", valStr)
				os.Exit(1)
			}
			interval = val

		} else if strings.HasPrefix(arg, "-w") {
			// handle -w flag
			wrapperFileLocation = strings.TrimPrefix(arg, "-w")
			if wrapperFileLocation == "" {
				fmt.Fprintln(os.Stderr, "Error: -w flag requires a file path.")
				os.Exit(1)
			}

		} else if strings.HasPrefix(arg, "-h") {
			// handle -h flag
			hiddenFileLocation = strings.TrimPrefix(arg, "-h")
			if hiddenFileLocation == "" {
				fmt.Fprintln(os.Stderr, "Error: -h flag requires a file path.")
				os.Exit(1)
			}
		}
	}
}

// MAIN FUNC
func main() {
	// 1. Manually parse combined arguments first
	parseManualArgs()

	// 2. Filter os.Args to contain only the program name and the simple boolean flags (-s, -r, -b, -B).
	// This prevents flag.Parse() from seeing and erroring on combined arguments like "-o1024".
	var filteredArgs []string
	for _, arg := range os.Args {
		// Keep the program name (os.Args[0])
		if arg == os.Args[0] {
			filteredArgs = append(filteredArgs, arg)
			continue
		}
		// Only keep the simple, non-value flags for flag.Parse()
		if arg == "-s" || arg == "-r" || arg == "-b" || arg == "-B" {
			filteredArgs = append(filteredArgs, arg)
		}
	}

	// 3. Temporarily set os.Args to the filtered list for flag.Parse()
	// We only need to do this if arguments were actually passed.
	if len(os.Args) > 1 {
		os.Args = filteredArgs
	}

	// 4. Collect boolean command-line arguments
	flag.Parse()

	// input sanitization
	if store == retrieve {
		fmt.Println("-s or -r required.")
		os.Exit(1)
	}
	if useBits == useBytes {
		fmt.Println("-b or -B required.")
		os.Exit(1)
	}
	// Validation now checks manually parsed values
	if wrapperFileLocation == "" {
		fmt.Println("wrapper file required (-w).")
		os.Exit(1)
	}
	if store && hiddenFileLocation == "" {
		fmt.Println("hidden file (-h) required when storing (-s).")
		os.Exit(1)
	}

	// read the wrapper file
	wrapperContents := readWrapper()

	// take action based on the command arguments
	var result []byte
	if retrieve { // retrieve
		result = retrieveData(wrapperContents)
	} else { // store
		result = storeData(wrapperContents)
	}

	// handle the resulting data
	_, e := os.Stdout.Write(result)
	if e != nil {
		fmt.Println("Failure to write output: ", e)
		os.Exit(1)
	}

}

// HELPER FUNCS

// read the wrapper file
func readWrapper() []byte {
	// read and store
	wrapperContents, e := os.ReadFile(wrapperFileLocation)
	if e != nil {
		fmt.Println("Failure to read wrapper file: ", e)
		os.Exit(1)
	}
	return wrapperContents
}

// METHODS FOR THE RETRIEVAL AND STORING OF THE HIDDEN DATA

// retrieve hidden data
func retrieveData(wrapperContents []byte) []byte {

	// read the hidden data based on the chosen method
	if useBits {
		// required variables
		var hiddenBytes []byte
		currentOffset := offset
		wrapperLength := len(wrapperContents)
		sentinelMatch := 0

		// runs until end of wrapper, or sentinel found
		for currentOffset < wrapperLength {
			hiddenByte := byte(0)

			// collect a byte from the 8 least significant bits in the wrapper
			for j := 0; j < 8; j++ {
				// if wrapper ends before full byte obtained
				if currentOffset >= wrapperLength {
					fmt.Println("Failure to obtain byte. Reached the end of wrapper.")
					os.Exit(1)
				}

				// collect lsb from current byte of wrapper
				bit := wrapperContents[currentOffset] & 0x01

				if READDIRECTION == "L2R" { // MSB-first
					// shift existing bits left
					hiddenByte <<= 1
					// store the new bit
					hiddenByte |= bit

				} else if READDIRECTION == "R2L" { // LSB-first
					// shift new bit left
					bitToPlace := bit << j
					// combine bit with the existing hiddenByte
					hiddenByte |= bitToPlace
				}

				// increment offset by interval
				currentOffset += interval
			}

			// once the sentinel is found, the matching continues
			if sentinelMatch < SENTINELSIZE && hiddenByte == SENTINEL[sentinelMatch] {
				sentinelMatch++
				if sentinelMatch == SENTINELSIZE {
					return hiddenBytes // all hidden bytes found
				}
			} else {
				// if partially matching, the bytes are part of the hidden message and must be added
				if sentinelMatch > 0 {
					hiddenBytes = append(hiddenBytes, SENTINEL[:sentinelMatch]...)
				}
				// add the current byte to the hidden data.
				hiddenBytes = append(hiddenBytes, hiddenByte)
				// set the match count back to 0
				sentinelMatch = 0
			}
		}

		// if sentinel is never found
		fmt.Println("Failure to find sentinel in wrapper file.")
		os.Exit(1)
		return nil

	} else { //useBytes
		// required variables
		var hiddenBytes []byte
		currentOffset := offset
		wrapperLength := len(wrapperContents)
		sentinelMatch := 0

		// iterate through the wrapper
		for currentOffset < wrapperLength {
			b := wrapperContents[currentOffset]

			// check for match w sentinel
			if sentinelMatch < SENTINELSIZE && b == SENTINEL[sentinelMatch] {
				// potentially a sentinel match
				sentinelMatch++
				if sentinelMatch == SENTINELSIZE {
					// found all hidden bytes
					return hiddenBytes
				}
			} else {
				// if partially matching, the bytes are part of the hidden message
				if sentinelMatch > 0 {
					hiddenBytes = append(hiddenBytes, SENTINEL[:sentinelMatch]...)
				}
				// add the current byte to the hidden data.
				hiddenBytes = append(hiddenBytes, b)
				// set the match count back to 0
				sentinelMatch = 0
			}

			// increment
			currentOffset += interval
		}

		// if sentinel is never found
		fmt.Println("Failure to find sentinel in wrapper file.")
		os.Exit(1)
		return nil
	}
}

// store the hidden data
func storeData(wrapperContents []byte) []byte {
	// read the hidden data
	hiddenData, e := os.ReadFile(hiddenFileLocation)
	if e != nil {
		fmt.Println("Failure to read hidden file: ", e)
		os.Exit(1)
	}

	// add the sentinel
	hiddenData = append(hiddenData, SENTINEL...) // put it at the end

	// store the data in proper format
	if useBits {
		// required variables
		currentOffset := offset
		wrapperLength := len(wrapperContents)

		// iterate through the data to hide
		for i := 0; i < len(hiddenData); i++ {
			hiddenByte := hiddenData[i]

			// 8 bits of byte
			for j := 0; j < 8; j++ {
				// if value of interval/offset is too high
				if currentOffset >= wrapperLength {
					fmt.Println("Failure to store data. interval or offset too large for wrapper.")
					os.Exit(1)
				}

				var bitToHide byte

				if READDIRECTION == "L2R" { // MSB-first
					bitToHide = (hiddenByte >> 7) & 0x01
					hiddenByte <<= 1

				} else if READDIRECTION == "R2L" { // LSB-first
					bitToHide = hiddenByte & 0x01
					hiddenByte >>= 1
				}

				// clear wrapper's lsb
				wrapperContents[currentOffset] &= 0xFE

				// store bit
				wrapperContents[currentOffset] |= bitToHide

				currentOffset += interval
			}
		}
		return wrapperContents
	} else { // useBytes
		// required variables
		currentOffset := offset
		wrapperLength := len(wrapperContents)
		hiddenLen := len(hiddenData)

		for i := 0; i < hiddenLen; i++ {
			// if value of interval/offset is too high
			if currentOffset >= wrapperLength {
				fmt.Fprintln(os.Stderr, "Failure to sore data. interval or offset too large for wrapper.")
				os.Exit(1)
			}
			// place hidden byte as wrapper byte
			wrapperContents[currentOffset] = hiddenData[i]
			currentOffset += interval
		}
		return wrapperContents
	}
}