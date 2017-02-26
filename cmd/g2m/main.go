package main

import (
	"errors"
	"flag"
	"os"

	g2m "github.com/kalambet/ghost2medium/ghost2medium"
	utils "github.com/kalambet/go-utils"
)

func main() {
	archive, token, migarte, err := parseArgs()
	if err != nil {
		utils.ExitWithError(err.Error())
	}

	posts, err := g2m.DecodeJSONArchive(*archive)
	if err != nil {
		utils.ExitWithError(err.Error())
	}

	err = g2m.Import(*token, *migarte, posts)
	if err != nil {
		utils.ExitWithError(err.Error())
	}

	utils.PrintInColorln("Finished ✅ \n", utils.Green)
}

func parseArgs() (archive *string, token *string, migrate *bool, err error) {
	archive = flag.String("archive", "", "Path to ghost archive JSON")
	token = flag.String("token", "", "Self issues Medium access token")
	migrate = flag.Bool("migrate", false, "Without this options we'll Do all steps except actual posting")

	flag.Parse()

	if *migrate {
		utils.PrintInColorln("\n\tThis run will import all the entries", utils.Yellow)
	} else {
		utils.PrintInColorln("\n\tThis is a test run", utils.Yellow)
	}

	if archive == nil || *archive == "" {
		return nil, nil, nil, errors.New("In order to migrate yor blog 📑  to Medium you must provide path to JSON file with archive 🗄")
	} else if _, err := os.Stat(*archive); os.IsNotExist(err) {
		return nil, nil, nil, errors.New("Path to JSON file with archive 🗄  you provided is incorrect of file can't be found")
	} else if token == nil || *token == "" {
		return nil, nil, nil, errors.New("In order to migrate yor blog 📑  to Medium you must provide Self Issued Access Token 🔓")
	}

	return
}
