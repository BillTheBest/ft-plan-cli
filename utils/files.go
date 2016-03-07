/**
* This is used everywhere. So it needs to be in a shared space.
**/

package utils

import "os"

// CloseAndDelete closes the file and then deletes it.
func CloseAndDelete(file *os.File, filePath string) {
	file.Close()
	os.Remove(filePath)
}
