package main

import (
	"btrzaws"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	versionNumber = 1
)

func main() {
	operation := flag.String("op", "", "operation to perform: update,fix-mail,fix-reg,smtp")
	repo := flag.String("repo", "", "repository to post in the registry")
	env := flag.String("env", "", "repository environment")
	flag.Parse()
	if *operation == "update" {
		runDefault()
	} else {
		if *operation == "fix-email" {
			fixEmail(*env)
		} else if *operation == "fix-reg" {
			fixRegistration(*repo, *env)
		}
		if *operation == "smtp" {
			translateEmailKey()
		}
	}
}

func translateEmailKey() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("please enter aws key value")
	key, _ := reader.ReadString('\n')
	key = strings.Trim(key, " \n")
	smtpCode, _ := btrzaws.GenerateSMTPPasswordFromSecret(key)
	fmt.Printf("for key '%s':\r\nsmtp code: '%s'\r\nDone.\r\n", key, smtpCode)
}

func fixRegistration(repo, env string) {
	if repo == "" {
		fmt.Println("No repo name to fix. exiting.")
		os.Exit(1)
	}
	if env == "" {
		fmt.Println("No repo env to fix. exiting.")
		os.Exit(1)
	}
	EnsureRepoRegistered(repo, env)
}

func fixEmail(environment string) {
	if environment == "" {
		log.Println("You must provide running environment")
		os.Exit(1)
	}
	updateMissingEmailInfo(environment)
}

func runDefault() {
	log.Println("No flags, delay for 5...")
	time.Sleep(time.Second * 5)
	log.Println("starting.")
	updateGroupsAndUsers()
}
