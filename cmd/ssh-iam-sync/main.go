package main

import (
	"internal/awssync"
	"internal/config"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	c := config.GetAppConfig()
	keys := awssync.GetSSHKeys(c)

	log.Println("Writing SSH keys to file: " + c.AuthorizedKeys)
	var file *os.File

	// Generate absolute path to load
	p, err := filepath.Abs(c.AuthorizedKeys)
	if err != nil {
		log.Fatal(err)
	}
	c.AuthorizedKeys = p

	if c.Overwrite {
		file, _ = os.Create(c.AuthorizedKeys)
	} else {
		file, _ = os.OpenFile(c.AuthorizedKeys, os.O_APPEND|os.O_WRONLY, 0600)
	}

	file.Write([]byte(strings.Join(keys, "\n")))

	defer file.Close()
}
