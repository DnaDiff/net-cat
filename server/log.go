package server

import (
	"fmt"
	"io"
	"os"
)

func enableLogging() (os.File, error) {
	// Create a pipe to capture standard output
	r, w, err := os.Pipe()
	if err != nil {
		return *w, err
	}

	// Set standard output to the write end of the pipe
	os.Stdout = w

	// Asynchronously write the standard output to a file
	go func() {
		f, err := os.Create("output.log")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		_, err = io.Copy(f, r)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	return *w, nil
}
