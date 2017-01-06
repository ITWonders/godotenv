package godotenv

import (
	"bufio"
	"os"
)

// ReadEnv same as Read, except also getting value from Environment
// if the variable exists
func ReadEnv(filenames ...string) (envMap map[string]string, err error) {
	filenames = filenamesOrDefault(filenames)
	envMap = make(map[string]string)

	for _, filename := range filenames {
		individualEnvMap, individualErr := readFileEnv(filename)

		if individualErr != nil {
			err = individualErr
			return // return early on a spazout
		}

		for key, value := range individualEnvMap {
			envMap[key] = value
		}
	}

	return
}

func readFileEnv(filename string) (envMap map[string]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	envMap = make(map[string]string)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for _, fullLine := range lines {
		if !isIgnoredLine(fullLine) {
			key, value, err := parseLine(fullLine)

			if os.Getenv(key) != "" {
				envMap[key] = os.Getenv(key)
			} else if err == nil {
				envMap[key] = value
			}
		}
	}
	return
}
