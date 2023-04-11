package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
)

// UniqueStringSlice returns unique items in a slice
func UniqueStringSlice(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}

func BytesDiff(b1, b2 []byte, path string) (data []byte, err error) {
	f1, err := os.CreateTemp("", "")
	if err != nil {
		return
	}
	defer os.Remove(f1.Name())
	defer f1.Close()

	f2, err := os.CreateTemp("", "")
	if err != nil {
		return
	}
	defer os.Remove(f2.Name())
	defer f2.Close()

	_, _ = f1.Write(b1)
	_, _ = f2.Write(b2)

	data, err = exec.Command("diff", "--label=old/"+path, "--label=new/"+path, "-u", f1.Name(), f2.Name()).CombinedOutput()
	if len(data) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	}
	return
}

func WritetoFile(filepath, content string) {
	f, err := os.Create(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func ReadFile(f string) ([]byte, error) {
	var data []byte
	var err error
	// if filename is an empty string it is a stdin
	if f == "" {
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			return data, err
		}
		log.Debug("Reading from stdin")
	} else {
		data, err = os.ReadFile(f)
		if err != nil {
			return data, err
		}
		log.Debugf("Reading from file %s", f)
	}
	return data, nil
}
