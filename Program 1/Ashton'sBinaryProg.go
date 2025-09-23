package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// filterBits takes a string and keeps only the 0s and 1s
func filterBits(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '0' || s[i] == '1' {
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

// decodeChunks splits bits into groups (7 or 8), turns them into ascii
// it returns the decoded text, how many bad bytes we saw, and how many bits were leftover
func decodeChunks(bits string, chunk int) (string, int, int) {
	invalid := 0
	leftover := len(bits) % chunk

	// go through each chunk
	runes := make([]rune, 0, len(bits)/chunk)
	for i := 0; i+chunk <= len(bits); i += chunk {
		part := bits[i : i+chunk]
		val := 0
		for j := 0; j < chunk; j++ {
			val = (val << 1) | int(part[j]-'0')
		}

		// if ascii value is 8 (backspace), remove last char
		if val == 8 {
			if len(runes) > 0 {
				runes = runes[:len(runes)-1]
			}
			continue
		}

		// keep only normal printable stuff or whitespace
		if (val >= 32 && val <= 126) || val == 9 || val == 10 || val == 13 {
			runes = append(runes, rune(val))
		} else {
			invalid++
		}
	}

	return string(runes), invalid, leftover
}

// chooseBest tries both 7-bit and 8-bit, picks the one with fewer errors
// tie rules: less leftover bits, then longer text, else choose 8-bit
func chooseBest(bits string) (string, int) {
	out7, inv7, rem7 := decodeChunks(bits, 7)
	out8, inv8, rem8 := decodeChunks(bits, 8)

	if inv7 < inv8 {
		return out7, 7
	} else if inv8 < inv7 {
		return out8, 8
	}

	if rem7 < rem8 {
		return out7, 7
	} else if rem8 < rem7 {
		return out8, 8
	}

	if len(out7) > len(out8) {
		return out7, 7
	} else if len(out8) > len(out7) {
		return out8, 8
	}

	return out8, 8
}

func main() {
	// read all text from stdin
	reader := bufio.NewReader(os.Stdin)
	var sb strings.Builder
	for {
		chunk, err := reader.ReadString('\n')
		sb.WriteString(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "read error:", err)
			os.Exit(1)
		}
	}

	// clean input to only bits
	clean := filterBits(sb.String())
	if len(clean) == 0 {
		return
	}

	// decode and print result (no extra newline)
	decoded, _ := chooseBest(clean)
	fmt.Print(decoded)
}
