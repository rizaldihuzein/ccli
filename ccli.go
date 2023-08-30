package ccli

import (
	"flag"
	"fmt"
	"strings"

	"github.com/rizaldihuzein/ccli/src"
)

const (
	errorMSG      = "sorry we encountered error \n"
	errorPanicMSG = "sorry we encountered panic \n"
)

func panicWrapper(f func()) {
	if f == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(errorPanicMSG, r)
		}
	}()

	f()
}

func ProcessCommand() {
	panicWrapper(processCommand)
}

func processCommand1() {
	src.Build()
	data, err := src.GetFromSource()
	if err != nil {
		fmt.Println(errorMSG, err)
		return
	}

	var csvPath string
	if csvPath == "" {
		csvPath = "data.csv"
	}

	err = src.SetAndReplaceToCSV(data, csvPath)
	if err != nil {
		fmt.Println(errorMSG, err)
		return
	}
}

func processCommand2(tags []string) {
	src.Build()
	data, err := src.SearchFromCSV(tags, "data.csv")
	if err == src.ErrMissingFile {
		fmt.Println("CSV file not found, generating new one...")
		processCommand1()
		data, err = src.SearchFromCSV(tags, "data.csv")
	}
	if err != nil {
		fmt.Println(errorMSG, err)
		return
	}
	for _, v := range data {
		fmt.Printf("ID: %s, Balance: %s\n", v.ID, v.Balance)
	}
}

func processCommand() {
	var tagStr = flag.String("tag", "-1", "tags to search separated by comma")
	flag.Parse()
	if tagStr == nil || *tagStr == "-1" {
		processCommand1()
		fmt.Println("No tags found.\nGenerating CSV instead...\nTo search data, please use -tag flag\n e.g. -tag=sed,quis")
	}
	if tagStr != nil {
		tags := strings.Split(*tagStr, ",")
		if *tagStr == "" {
			tags = []string{}
		}
		processCommand2(tags)
	}
}
