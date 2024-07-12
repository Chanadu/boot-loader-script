package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
)

func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("[isRoot] Unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
}

func getEntryType(fileName string) (entryType string, isFallback bool, err error) {
	isFallback = false
	entryType = ""
	lowerFileName := strings.ToLower(fileName)
	if strings.Contains(lowerFileName, "arch") {
		entryType = "arch"
	}
	if strings.Contains(lowerFileName, "windows") {
		entryType = "windows"
	}
	if strings.Contains(lowerFileName, "lts") {
		entryType = "lts"
	}
	if strings.Contains(lowerFileName, "fallback") {
		isFallback = true
	}
	if entryType == "" {
		err = fmt.Errorf("Unknown entry type")
		return "", false, err
	}
	return entryType, isFallback, nil
}

func main() {
	if !isRoot() {
		log.Fatalf("You must be root to run this script")
		os.Exit(1)
	}

	fmt.Println("Script Starting!")
	fmt.Println()

	files, err := os.ReadDir("/efi/loader/entries")
	if err != nil {
		log.Fatalf("Unable to read directory: %s", err)
		os.Exit(1)
	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
		if file.IsDir() {
			continue
		}

		fileName := file.Name()

		fileData, err := os.ReadFile("/efi/loader/entries/" + fileName)
		if err != nil {
			log.Fatalf("Unable to read file: %s", err)
			os.Exit(1)
		}
		// fmt.Println()
		// fmt.Println(string(fileData))
		// fmt.Println()

		fileLines := strings.Split(string(fileData), "\n")
		entryType, isFallback, err := getEntryType(fileName)
		if err != nil {
			log.Fatalf("Unable to determine entry type: %s", err)
			os.Exit(1)
		}

		entryTypeToTitle := map[string]string{
			"arch":    "EndeavourOS-Arch",
			"windows": "Windows",
			"lts":     "EndeavourOS-LTS",
		}

		entryTypeToSortKeyAddition := map[string]int{
			"arch":    1,
			"windows": 2,
			"lts":     3,
		}
		fallbackSortKeyAddition := 3

		for _, line := range fileLines {
			if strings.Contains(line, "title") {
				fallbackText := ""
				if isFallback {
					fallbackText = "-Fallback"
				}

				fmt.Println(line)
				fmt.Println("Replace with :", entryTypeToTitle[entryType]+fallbackText)
			} else if strings.Contains(line, "sort-key") {
				fmt.Println(line)
				sortkeyNumber := entryTypeToSortKeyAddition[entryType]
				if isFallback {
					sortkeyNumber += fallbackSortKeyAddition
				}
				splitLine := strings.Split(line, " ")
				newSplitLine := []string{}
				for _, item := range splitLine {
					if !(item == "" || item == " ") {
						newSplitLine = append(newSplitLine, item)
					}
				}
				splitLine = newSplitLine

				if len(splitLine) != 2 {
					log.Fatalf("Invalid sort-key line: %s", line)
					os.Exit(1)
				}
				fmt.Printf("Replace with: %s %d-%s",
					splitLine[0],
					sortkeyNumber,
					splitLine[1],
				)
			}
		}
	}

}
