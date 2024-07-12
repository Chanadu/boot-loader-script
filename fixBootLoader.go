package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
)

func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("[isRoot] Unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
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
		// if file.IsDir() {
		// 	continue
		// }

		// fileName := file.Name()
		// if fileName == "arch.conf" {
		// 	continue
		// }

		// err := os.Remove("/efi/loader/entries/" + fileName)
		// if err != nil {
		// 	os.Exit(1)
		// }
	}
}
