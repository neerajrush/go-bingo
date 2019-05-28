package main

import (
	"os"
	"fmt"
	"github.com/ChimeraCoder/gojson"
)

// Convert beam-forming json document to go struct
func json2GoStruct(apiName string) (goStruct string) {
	jsonFile := apiName + ".json"
	api, err := os.Open(jsonFile)
	if err != nil {
		fmt.Println("error opening", jsonFile, err)
		return
	}
	structName := ""
	for _,c := range apiName {
		if c == '-' || c == '_' {
			continue
		}
		structName += string(c)
	}

	actual, err := gojson.Generate(api, gojson.ParseJson, structName, "gojson", []string{"json"}, false, true)

	if err != nil {
		fmt.Println("error parsing", jsonFile, err)
		return
	}

	goStruct = string(actual)
	return
}

func main() {
	api_list := []string{"Beam-Forming", "Dynamic-Frequency-Change"}
	for _,v := range api_list {
		go_struct := json2GoStruct(v)
		fmt.Println(go_struct)
	}
	fmt.Println(api_list)
}
