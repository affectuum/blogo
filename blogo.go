package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	// "strings"
	// "bufio"
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
	// fmt.Println(mdSplit)
	mdParse := blockParsing(mdSplit)
	fmt.Println("--------------------------------------------------------")
	fmt.Println("Final result: ", mdParse)
}

func filePrepper(file []byte) string {
	// fileLines := strings.Split(string(file), "\r\n")
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
	data, err2 := ioutil.ReadAll(jsonFile)
	if err2 != nil {
		fmt.Println("error:", err2)
	}
	err3 := json.Unmarshal(data, &parsingMap)
	if err3 != nil {
		fmt.Println("error:", err3)
	}
	endResult := lines
	fmt.Println(endResult)
	for i := 0; i < len(parsingMap.Elem); i++ {
		matchPattern := `(?m)` + parsingMap.Elem[i].MatchPattern
		re, err := regexp.Compile(matchPattern)
		if err != nil {
			fmt.Print(err)
		}
		subst := parsingMap.Elem[i].SubstitutionPattern
		// endResult = re.ReplaceAllString(lines, subst)
		fmt.Printf("Processing: %s\n", parsingMap.Elem[i].ElementName)
		fmt.Printf("Pattern: %s\n", matchPattern)
		fmt.Printf("Substitution: %s\n", subst)
		matches := re.FindAllString(endResult, -1)
		if len(matches) > 0 {
			fmt.Printf("Found %d matches\n", len(matches))
			for _, match := range matches {
				fmt.Printf("Match: %s\n", match)
			}
		}
		endResult = re.ReplaceAllStringFunc(endResult, func(match string) string {
			return re.ReplaceAllString(match, subst)
		})
		fmt.Printf("Result after " + parsingMap.Elem[i].ElementName + ": \n" + endResult)
		fmt.Printf("\n###########################################################################################\n")
	}
	return endResult
}
