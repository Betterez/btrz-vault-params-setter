package main

import (
	"btrzaws"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	versionNumber = 1
)

func main() {
	operation := flag.String("op", "", "operation to perform: update,fix-mail,fix-reg,smtp,vault")
	repo := flag.String("repo", "", "repository to post in the registry")
	env := flag.String("env", "", "repository environment")
	filename := flag.String("vault_file", "", "vault file name (full path) to process")
	flag.Parse()
	if *operation == "update" {
		updateGroupsAndUsers()
	} else if *operation == "vault" {
		if *repo != "" && *env != "" {
			if err := showVaultDataForRepo(*repo, *env); err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}
		} else if *filename != "" {
			if err := updateVault(*filename); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			fmt.Println("bad params")
		}
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
