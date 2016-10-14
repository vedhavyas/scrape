package utils

import "os"

//WriteInLine will write the data in single line by erasing the previously written data
func WriteInline(data string) error {
	_, err := os.Stdout.Write([]byte(data + "\r"))
	if err != nil {
		return err
	}
	return os.Stdout.Sync()
}
