package utils

import (
	"fmt"
	"os"
	"strings"
)

func ReadENV(key string) (value string) {
	var result string

	path, err := os.Getwd()

	if err != nil {
		fmt.Println("Error occurred while getting path:", err.Error())
	}

	_, ok := strings.CutSuffix(path, "/")

	var filepath string
	if ok {
		filepath = path + ".env"
	} else {
		filepath = path + "/" + ".env"
	}
	content, err := os.ReadFile(filepath)

	if err != nil {
		fmt.Println("Error occurred while reading file ", filepath, " error: ", err.Error())
	}

	lines := strings.Split(string(content[:]), "\n")

	for _, line := range lines {
		if strings.Contains(line, key) {
			result = strings.Split(line, fmt.Sprintf("%s=", key))[1]
		}
	}

	if strings.Contains(result, `"`) || strings.Contains(result, "'") {
		result = strings.ReplaceAll(result, `"`, "")
		result = strings.ReplaceAll(result, `'`, "")
	}

	return result
}
