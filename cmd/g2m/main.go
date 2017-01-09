package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	medium "github.com/kalambet/medium-sdk-go"
)

func main() {
	_, token := parseArgs()

	m := medium.NewClientWithAccessToken(*token)

	u, err := m.GetUser()
	if err != nil {
		exitWithError(err.Error())
	}

	publications, err := m.GetPublications(u.ID)
	if err != nil {
		exitWithError(err.Error())
	}

	if len(*publications) == 0 {
		exitWithError("There are no publications where you can import your archive.")
	}

	printInColor("List of all Publications you can contribute to:\n", Green)
	for idx, pub := range *publications {
		contributors, err := m.GetContributors(pub.ID)
		if err != nil {
			exitWithError(err.Error())
		}

		for _, c := range *contributors {
			if c.UserID == u.ID {
				fmt.Printf("\t%d: %s (URL: %s)\n", idx, pub.Name, pub.URL)
			}
		}
	}

	fmt.Print("\nPlease select Publication you want import your archive to: ")
	reader := bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')

	pubIdx, err := strconv.Atoi(strings.Trim(userInput, "\n"))
	if err != nil || pubIdx > len(*publications) {
		exitWithError(err.Error())
	}

	fmt.Printf("Tags: %#v\n", (*publications)[pubIdx])

	//_, err = g2m.DecodeJSONArchive(*path)

	//if err != nil {
	//	exitWithError(err.Error())
	//}

}

// Text colors
const (
	Black int = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

func parseArgs() (archive *string, token *string) {
	archive = flag.String("archive", "", "Path to ghost archive JSON")
	token = flag.String("token", "", "Self issues Medium access token")

	flag.Parse()

	if archive == nil || *archive == "" {
		exitWithError("In order to migrate yor blog 📑  to Medium you must provide path to JSON file with archive 🗄")
	} else if _, err := os.Stat(*archive); os.IsNotExist(err) {
		exitWithError("Path to JSON file with archive 🗄  you provided is incorrect of file can't be found")
	} else if token == nil || *token == "" {
		exitWithError("In order to migrate yor blog 📑  to Medium you must provide Self Issued Access Token 🔓")
	}

	return
}

func logInColor(message string, color int) {
	log.Printf("\033[%dm%s\033[0m\n\n", color, message)
}

func printInColor(message string, color int) {
	fmt.Printf("\033[%dm%s\033[0m\n", color, message)
}

func exitWithError(message string) {
	logInColor(message, Red)
	flag.Usage()
	os.Exit(1)
}
