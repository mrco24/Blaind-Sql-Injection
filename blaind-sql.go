package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings" // Add this line to import the "strings" package
)


func main() {
    // Define command-line flags and parse them.
    urlFile := flag.String("u", "url.txt", "File containing a list of URLs")
    payloadFile := flag.String("p", "payload.txt", "File containing a list of payloads")
    verbose := flag.Bool("v", false, "Enable verbose output")
    outputFile := flag.String("o", "output.txt", "Output file to write results")
    flag.Parse()

    urls, err := readLines(*urlFile)
    if err != nil {
        fmt.Printf("Error reading URL file: %s\n", err)
        return
    }

    payloads, err := readLines(*payloadFile)
    if err != nil {
        fmt.Printf("Error reading payload file: %s\n", err)
        return
    }

    output, err := os.Create(*outputFile)
    if err != nil {
        fmt.Printf("Error creating output file: %s\n", err)
        return
    }
    defer output.Close()

    // Define ANSI escape codes for red color
    redColor := "\033[91m"
    resetColor := "\033[0m"

    for _, url := range urls {
        for _, payload := range payloads {
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
                result := fmt.Sprintf("Not Vulnerable: %s (Content Length Unchanged)\n", modifiedURL) // Mark as "Not Vulnerable"
                if *verbose {
                    fmt.Println(result)
                }
                output.WriteString(result)
            } else {
                result := fmt.Sprintf("%sVulnerable: %s (Content Length Changed)%s\n", redColor, modifiedURL, resetColor) // Mark as "Vulnerable" in red
                fmt.Println(result)
                output.WriteString(result)
            }
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
