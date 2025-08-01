package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

type elementToParse struct {
	ElementName         string `json:"elementName"`
	MatchPattern        string `json:"matchPattern"`
	SubstitutionPattern string `json:"substitutionPattern"`
}

type elementsToParse struct {
	Elem []elementToParse `json:"elementsToParse"`
}

func main() {
	filePath := os.Args[1]
	markdownFile, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	mdSplit := filePrepper(markdownFile)
	mdParse := blockParsing(mdSplit)
	fmt.Println("<html><head><title>Result</title></head><body>", mdParse, "</body></html>")
}

func filePrepper(file []byte) string {
	fileLines := string(file)
	return fileLines
}

func blockParsing(lines string) string {
	var parsingMap elementsToParse
	jsonFile, err := os.Open("./elements.json")
	if err != nil {
		fmt.Print(err)
	}
	defer jsonFile.Close()
	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("error:", err)
	}
	err = json.Unmarshal(data, &parsingMap)
	if err != nil {
		fmt.Println("error:", err)
	}
	endResult := lines
	for i := 0; i < len(parsingMap.Elem); i++ {
		matchPattern := `(?m)` + parsingMap.Elem[i].MatchPattern
		re, err := regexp.Compile(matchPattern)
		if err != nil {
			fmt.Print(err)
		}
		subst := parsingMap.Elem[i].SubstitutionPattern
		endResult = re.ReplaceAllStringFunc(endResult, func(match string) string {
			return re.ReplaceAllString(match, subst)
		})
	}
	return endResult
}
