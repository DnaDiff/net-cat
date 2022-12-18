package server

import (
	"fmt"
	"io"
	"os"
	"time"
)

var mw io.Writer

func enableLogging(isEnabled bool) (os.File, error) {
	if !isEnabled {
		mw = io.MultiWriter(os.Stdout)
		return *os.Stdout, nil
	}

	// Create a pipe to capture standard output
	r, w, err := os.Pipe()
	if err != nil {
		return *w, err
	}
	mw = io.MultiWriter(w, os.Stdout)

	// Asynchronously write the standard output to a file
	go func() {
		err := os.Mkdir("./logs", 0755)
		if err != nil && !os.IsExist(err) {
			fmt.Fprintln(mw, err)
			return
		}
		f, err := os.Create("./logs/output_" + time.Now().Format("2006.01.02_15.04") + ".log")
		if err != nil {
			fmt.Fprintln(mw, err)
			return
		}
		defer f.Close()

		_, err = io.Copy(f, r)
		if err != nil {
			fmt.Fprintln(mw, err)
			f.Close()
			return
		}
	}()

	return *w, nil // REMEMBER TO CALL w.Close() WHEN CALLING THIS FUNCTION
}
