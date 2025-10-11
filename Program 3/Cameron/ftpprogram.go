/*
author: 	AlaskanBullWorm
date:		9/28/25
description:	a Go program to decode storage covert channel messages through ftp
*/

package main

// import the required libraries
import (
	"fmt"     // for printing
	"log"     // for debugginig errors
	"strconv" // library for converting binary strings to ints

	"github.com/secsy/goftp" // third party library for handling ftp connections
)

// create some constant variables
const (
	METHOD   = "10"              // variable holding either 10 or 7 (number of right-most 'bits' to consider)
	address  = "138.47.99.161" // ip address of ftp server
	port     = "21"            // port of ftp server
	username = "anonymous"     // ftp server username
	password = ""              // ftp server password
	path     = "/"             // path on ftp server to connect to
)

// beginning of main runtime
func main() {

	// convert from string to int
	METHODINT, _ := strconv.Atoi(METHOD)

	// connect to a ftp server
	// give necessary info to goftp library
	config := goftp.Config{
		User:            username,
		Password:        password,
		ActiveTransfers: true,
	}
	// connect
	client, err := goftp.DialConfig(config, address+":"+port)
	// error checking
	if err != nil {
		log.Fatal(err)
	}
	// fetch the file listing and respective permissions
	entries, err := client.ReadDir(path)
	// error checking
	if err != nil {
		log.Fatal(err)
	}
	// leave the ftp server
	client.Close()

	// convert each file info into readable text
	// variable to store each file's permissions
	perms := []string{}
	// append to the above list each file's string of permissions
	for _, entry := range entries {
		perms = append(perms, entry.Mode().String())
	}

	fmt.Println(perms)

	// determine if the first 7 or 10 right-most bits of the permission
	// strings should be considered
	if METHOD == "7" {
		// create new slice to hold values
		tempSlice := []string{}
		// slice down to only the 7 right-most bits
		for _, perm := range perms {
			tempSlice = append(tempSlice, perm[3:])
		}
		perms = tempSlice
	}

	// if there all 10 bits are used, do nothing here
	// now concatonate all characters of all strings together into a long string
	tempString := ""
	for _, perm := range perms {
		tempString += perm
	}

	// now we need to determine if each character of the long string is a '-' or anything else
	binaryStringOut := ""
	for _, char := range tempString {
		if char == '-' {
			binaryStringOut += "0"
		} else {
			binaryStringOut += "1"
		}
	}

	fmt.Println(binaryStringOut)

	// ensure that the binary string is % 7 by padding so that
	// there are no slicing errors/
	if METHOD == "10" {
		if len(binaryStringOut)%7 != 0 {
			for {
				binaryStringOut += "0"
				if len(binaryStringOut)%7 == 0 {
					break
				}
			}
		}
	}

	// now we need to convert the binary to ASCII
	output := ""
	// for each 8 bits of the string:
	for i := 0; i < len(binaryStringOut); i += 7 {

		// store the bits
		tempEight := binaryStringOut[i : i+7]

		fmt.Println()
		fmt.Println(tempEight)

		// convert these eight bits into an int
		decValue, err := strconv.ParseUint(tempEight, 2, METHODINT)
		// error checking
		if err != nil {
			log.Fatal(err)
		}
		output += string(byte(decValue))
	}
	fmt.Println(output)

}
