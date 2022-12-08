package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/google/shlex"
)

func readCrontabfile(path string) ([]cronjob, error) {
	readFile, err := os.Open("cronfile.txt")
	if err != nil {
		return nil, fmt.Errorf("Error opening file: %v", err)
	}
	defer readFile.Close()
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string
	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}
	result := make([]cronjob, 0)
	for _, line := range fileLines {
		val, err := shlex.Split(line)
		if err != nil {
			return nil, fmt.Errorf("Error parsing line: %v", err)
		}
		cj := cronjob{
			Crontab: strings.Join(val[0:6], " "),
			Command: val[6],
			Args:    val[7:],
		}
		result = append(result, cj)
	}

	return result, nil
}
