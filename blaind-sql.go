package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Define command-line flags and parse them.
	urlFlag := flag.String("u", "", "Single target URL")
	urlFileFlag := flag.String("f", "", "File containing a list of URLs")
	verbose := flag.Bool("v", false, "Enable verbose output")
	outputFile := flag.String("o", "output.txt", "Output file to write results")
	flag.Parse()

	var urls []string

	if *urlFlag != "" {
		urls = append(urls, *urlFlag)
	} else if *urlFileFlag != "" {
		urlsFromFile, err := readLines(*urlFileFlag)
		if err != nil {
			fmt.Printf("Error reading URLs from %s: %v\n", *urlFileFlag, err)
			return
		}
		urls = append(urls, urlsFromFile...)
	}

	output, err := os.Create(*outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %s\n", err)
		return
	}
	defer output.Close()

	// Define ANSI escape codes for colors
	blueColor := "\033[94m"
	redColor := "\033[91m"
	resetColor := "\033[0m"

	for _, url := range urls {
		requestURL := fmt.Sprintf("Request URL: %s", url)
		fmt.Printf("%s%s%s\n", blueColor, requestURL, resetColor)

		payload := "%27%22%60" // Default payload

		modifiedURL := url + payload
		originalLength, err := getContentLength(url)
		if err != nil {
			fmt.Printf("Error fetching content length for %s: %s\n", url, err)
			continue
		}
		modifiedLength, err := getContentLength(modifiedURL)
		if err != nil {
			fmt.Printf("Error fetching content length for %s: %s\n", modifiedURL, err)
			continue
		}
		if originalLength != modifiedLength {
			result := fmt.Sprintf("Not Vulnerable: %s (Content Length Unchanged)\n", modifiedURL)
			if *verbose {
				fmt.Println(result)
			}
		} else {
			vulnerableResult := fmt.Sprintf("Vulnerable: %s (Content Length Changed)\n", modifiedURL)
			fmt.Printf("%s%s%s\n", redColor, vulnerableResult, resetColor)
			output.WriteString(vulnerableResult) // Write only for vulnerable URLs
		}
	}
}

func getContentLength(url string) (int, error) {
	response, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	return int(response.ContentLength), nil
}

func readLines(filename string) ([]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := []string{}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}
