package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	versionNumber = 1
)

func main() {
	operation := flag.String("op", "", "operation to perform")
	repo := flag.String("repo", "", "repository to post in the registry")
	env := flag.String("env", "", "repository environment")
	flag.Parse()
	if *operation == "" {
		log.Println("No flags, delay for 5...")
		time.Sleep(time.Second * 5)
		log.Println("starting.")
		updateGroupsAndUsers()
	} else {
		if *operation == "fix-email" {
			if *env == "" {
				log.Println("You must provide running environment")
				os.Exit(1)
			}
			updateMissingEmailInfo(*env)
		} else if *operation == "fix-reg" {
			if *repo == "" {
				fmt.Println("No repo name to fix. exiting.")
				os.Exit(1)
			}
			if *env == "" {
				fmt.Println("No repo env to fix. exiting.")
				os.Exit(1)
			}
			EnsureRepoRegistered(*repo, *env)
		}
	}
}
