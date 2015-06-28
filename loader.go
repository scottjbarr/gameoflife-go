// See http://psoup.math.wisc.edu/mcell/ca_files_formats.html
//
// This loader is a non-strict Life 1.05 parser. See docs/glider.L for an
// example of a loadable file.
package main

import (
	"bufio"
	"os"
)

// Read a Life data file and return the string data.
func ReadLifeData(path string) ([]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()
		if string(text[0]) != "#" {
			lines = append(lines, text)
		}
	}

	return lines, scanner.Err()
}
