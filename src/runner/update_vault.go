package main

import (
	"btrzutils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type VaultServerData struct {
	Keys    []string `json:"keys"`
	Address string   `json:"address"`
	Port    int      `json:"port"`
	Token   string   `json:"token"`
}

type VaultEntry struct {
	ServerData VaultServerData `json:"vault"`
}
type VaultRegistry struct {
	Servers map[string]VaultEntry `json:"servers"`
}

type vaultUpdater struct {
	Server btrzutils.VaultConnectionParameters `json:"server"`
	Values map[string]map[string]string        `json:"values"`
}

func loadVaultInfoFromFile(filename string) (*vaultUpdater, error) {
	var err error
	var data []byte
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		return nil, err
	}
	if data, err = ioutil.ReadFile(filename); err != nil {
		return nil, err
	}
	info := &vaultUpdater{}
	json.Unmarshal(data, info)
	return info, err
}

func updateVault(filename string) error {
	updater, err := loadVaultInfoFromFile(filename)
	if err != nil {
		return err
	}
	server, err := btrzutils.CreateVaultConnection(&updater.Server)
	if err != nil {
		fmt.Printf("%v\nexiting", err)
		os.Exit(1)
	}
	if server.GetVaultStatus() != btrzutils.VaultOnline {
		return fmt.Errorf("Server not online:%s", server.GetVaultStatus())
	}
	fmt.Println("server is online")
	for key := range updater.Values {
		jsonData, err := json.Marshal(updater.Values[key])
		if err != nil {
			return fmt.Errorf("error %v marshaling", err)
		}
		jsonString := string(jsonData)
		if _, err = server.SetValuesForRepository(key, jsonString, false); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		json, err := server.GetRepositoryValues(key)
		if err != nil {
			return err
		}
		fmt.Println(json)
	}

	return nil
}

func loadVaultRegistryFromFile(filename string) (*VaultRegistry, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, err
	}
	var bytes []byte
	var err error
	if bytes, err = ioutil.ReadFile(filename); err != nil {
		return nil, err
	}
	result := &VaultRegistry{}
	if err = json.Unmarshal(bytes, result); err != nil {
		return nil, err
	}
	return result, nil
}

func loadVaultRegistry() (*VaultRegistry, error) {
	const (
		secretsFile = "secrets/secrets.json"
	)
	return loadVaultRegistryFromFile(secretsFile)
}

func showVaultDataForRepo(repo, env string) error {
	servers, err := loadVaultRegistry()
	if err != nil {
		return err
	}
	serverInfo, ok := servers.Servers[env]
	if !ok {
		return fmt.Errorf("servers for %s were not found", env)
	}
	fmt.Println("got server info for env")
	vaultServer, err := btrzutils.CreateVaultConnectionFromParameters(serverInfo.ServerData.Address, serverInfo.ServerData.Token, serverInfo.ServerData.Port)
	if err != nil {
		return err
	}
	fmt.Printf("values for %s", env)
	fmt.Println(vaultServer.GetRepositoryValues(repo))
	return nil
}
