package utils

import "os"

func WriteInline(data string) error {
	_, err := os.Stdout.Write([]byte(data + "\r"))
	if err != nil {
		return err
	}
	return os.Stdout.Sync()
}
