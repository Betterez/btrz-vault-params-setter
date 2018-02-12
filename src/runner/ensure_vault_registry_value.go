package main

import (
	"btrzutils"
	"log"
)

func EnsureRepoRegistered(repositoryName, environment string) (int, error) {
	const fileName = "secrets/secrets.json"
	params, err := btrzutils.LoadVaultInfoFromJSONFile(fileName, environment)
	if err != nil {
		log.Println(err, "loading valut info for ", environment)
		return 0, err
	}
	log.Println("connecting to vault server...")
	connection, err := btrzutils.CreateVaultConnection(params)
	if err != nil {
		log.Println(err, "connecting to vault server!")
		return 0, err
	}
	log.Println("ensuring repo", repositoryName)
	code, err := connection.EnsureRepositoryExists(repositoryName)
	if err != nil {
		log.Println(err, "ensuring repo")
	}
	log.Println("repo", repositoryName, code, err)
	return code, err
}
