package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"postmantidy/core"
	"strings"
)

/*
*
setup rules
*/
var sectionPathForFolder = 1
var groupedRequests = map[string][]core.ItemStructure{}
var hostReplaceUrl = "https://minerva.petrosea.com"

func main() {

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Printf("first argument is postman file")
		return
	}

	// mapping argument values
	postmanFilepath := args[0]

	if _, err := os.Stat(postmanFilepath); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		fmt.Printf("file %v not found", postmanFilepath)
		return
	}

	jsonFile, err := os.Open(postmanFilepath)
	if err != nil {
		fmt.Printf("json file cannot be open: %v", err.Error())
		return
	}

	fmt.Println("json file successfully opened")

	defer jsonFile.Close()

	byteJsonFile, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Printf("json file cannot be read: %v", err.Error())
		return
	}

	postData := core.PostmanStructure{}

	err = json.Unmarshal(byteJsonFile, &postData)
	if err != nil {
		fmt.Printf("json file cannot be unmarshal: %v", err.Error())
		return
	}

	newPostmanData := core.PostmanStructure{}

	newPostmanData.Info = postData.Info

	ProcessItem(&postData.Item, &newPostmanData.Item)

	BuildPostmanNewStructure(&newPostmanData)

	fJson, err := json.MarshalIndent(newPostmanData, "", " ")
	if err != nil {
		fmt.Printf("json file cannot be marshall indent: %v", err.Error())
		return
	}

	err = ioutil.WriteFile("test.json", fJson, 0644)
	if err != nil {
		fmt.Printf("json file cannot be write file: %v", err.Error())
		return
	}

}

func ProcessItem(sourcePostman *[]core.ItemStructure, destPostman *[]core.ItemStructure) {
	for _, item := range *sourcePostman {
		if !core.IsUniqueUrl(item) {
			continue
		}

		newName := item.Request.Url.Path[len(item.Request.Url.Path)-1]
		item.Name = newName

		// check authrization in header
		isNeedAuth := false
		for index, keyHeader := range item.Request.Header {
			if keyHeader.Key == "Authorization" {
				item.Request.Header = append(item.Request.Header[:index], item.Request.Header[index+1:]...)
				isNeedAuth = true
			}

			if keyHeader.Key == "X-Minerva-RBAC-Domain" {
				item.Request.Header[index].Value = "{{rbac_domain}}"
			}
		}

		if isNeedAuth {
			item.Request.Auth = &core.ItemRequestAuth{
				Type: "bearer",
				Bearer: []core.ItemRequestAuthBearer{
					{
						Key:   "token",
						Value: "{{token}}",
						Type:  "string",
					},
				},
			}
		}

		// setup host
		item.Request.Url.Host = []string{"{{host}}"}
		item.Request.Url.Protocol = nil
		item.Request.Url.Raw = strings.Replace(item.Request.Url.Raw, hostReplaceUrl, "{{host}}", 1)

		if _, ok := groupedRequests[item.Request.Url.Path[sectionPathForFolder]]; !ok {
			//do something here
			groupedRequests[item.Request.Url.Path[sectionPathForFolder]] = []core.ItemStructure{
				item,
			}
		} else {
			groupedRequests[item.Request.Url.Path[sectionPathForFolder]] = append(groupedRequests[item.Request.Url.Path[sectionPathForFolder]], item)

		}
	}
}

func BuildPostmanNewStructure(newStructure *core.PostmanStructure) {

	for key, _ := range groupedRequests {
		items := groupedRequests[key]
		newStructure.Item = append(newStructure.Item, core.ItemStructure{
			Name:    key,
			Item:    &items,
			Request: nil,
		})
	}

}
