package lib

import "os"

func CheckOrCreateDir(dirPath string) error {
	// Check if the directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// The directory does not exist (or maybe some other error, like a permission issue)
		err := os.MkdirAll(dirPath, 0755) // Create the directory and any necessary parents
		if err != nil {
			return err
		}
	}

	return nil
}
