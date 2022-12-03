package diff

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func writeTempFile(dir, prefix string, data []byte) (string, error) {
	file, err := os.CreateTemp(dir, prefix)
	if err != nil {
		return "", err
	}

	_, err = file.Write(data)
	if err1 := file.Close(); err == nil {
		err = err1
	}
	if err != nil {
		os.Remove(file.Name())
		return "", err
	}
	return file.Name(), nil
}

// Diff returns a human-readable description of the differences between s1 and s2.
func Diff(b1, b2 []byte) (string, error) {
	if bytes.Equal(b1, b2) {
		return "", nil
	}

	cmd := "diff"
	if _, err := exec.LookPath(cmd); err != nil {
		return "", fmt.Errorf("diff command unavailable\nold: %q\nnew: %q", b1, b2)
	}

	f1, err := writeTempFile("", "gen_test", b1)
	if err != nil {
		return "", err
	}
	defer os.Remove(f1)

	f2, err := writeTempFile("", "gen_test", b2)
	if err != nil {
		return "", err
	}
	defer os.Remove(f2)

	data, err := exec.Command(cmd, "-u", f1, f2).CombinedOutput()
	if len(data) > 0 { // diff exits with a non-zero status when the files don't match
		err = nil
	}
	if err != nil {
		return "", err
	}
	return string(data), nil
}
