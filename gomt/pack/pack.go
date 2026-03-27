package pack

// func PackSetupPacket(input string, output string) error {
// 	if ok := utils.ExistsDir(input); !ok {
// 		return errors.Errorf("invalid input path")
// 	}

// 	zipFiles, err := utils.GetFileList(input, ".zip")
// 	if err != nil {
// 		return errors.Wrap(err, "git zip files list")
// 	}
// 	txtFiles, err := utils.GetFileList(input, ".txt")
// 	if err != nil {
// 		return errors.Wrap(err, "git txt files list")
// 	}

// 	return nil
// }

// func sha256Hash(files ...string) (string, error) {
// 	// Create a new SHA256 hash.
// 	h := sha256.New()

// 	// Iterate over the files and hash their contents.
// 	for _, file := range files {
// 		f, err := os.Open(file)
// 		if err != nil {
// 			return "", err
// 		}
// 		defer f.Close()

// 		if _, err := io.Copy(h, f); err != nil {
// 			return "", err
// 		}
// 	}

// 	// Get the hash sum.
// 	sum := h.Sum(nil)

// 	// Convert the hash sum to a hexadecimal string.
// 	hex := ""
// 	for _, b := range sum {
// 		hex += string('0' + (b >> 4))
// 		hex += string('0' + (b & 0x0f))
// 	}

// 	return hex, nil
// }
