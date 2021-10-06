/**
 *
 * Copyright 2021 Victor Shinya
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var (
	wg sync.WaitGroup
)

type functionBodyRequest struct {
	Username string `json:"username"`
	Apikey string `json:"apikey"`
	On bool `json:"power"`
	VsiNames []string `json:"name"`
}

// VirtualServer model to handle only the useful data
type VirtualServer struct {
	ID         int    `json:"id"`
	DomainName string `json:"fullyQualifiedDomainName"`
}

// Main function to run the Action source code
func Main(params map[string]interface{}) map[string]interface{} {
	m, _ := json.Marshal(params)
	var bodyRequest functionBodyRequest
	_ = json.Unmarshal(m, &bodyRequest)

	var arr []string
	err := json.Unmarshal([]byte(params["name"].(string)), &arr)
	if err != nil {
		log.Fatalf("Error while parsing Virtual Server names into a string array: %v", err)
	}
	bodyRequest.VsiNames = arr

	body, err := httpRequest("GET", "https://api.softlayer.com/rest/v3.1/SoftLayer_Account/getVirtualGuests", nil, bodyRequest.Username, bodyRequest.Apikey)
	if err != nil {
		log.Fatalf("Error to retrieve the list of all Virtual Servers from Softlayer API: %v", err)
		return map[string]interface{}{
			"success": false,
			"statusCode": 500,
			"message": "Internal Server Error",
		}
	}
	var vsis []VirtualServer
	err = json.Unmarshal(body, &vsis)
	if err != nil {
		log.Fatalf("Error to parse Softlayer API response into Virtual Server struct: %v", err)
		return map[string]interface{}{
			"success": false,
			"statusCode": 500,
			"message": "Internal Server Error",
		}
	}
	for _, vsi := range vsis {
		for _, name := range bodyRequest.VsiNames {
			if vsi.DomainName == name {
				fmt.Printf("ID = %d; VSI Name = %s\n", vsi.ID, name)
				wg.Add(1)
				go power(bodyRequest.On, vsi, bodyRequest.Username, bodyRequest.Apikey)
			}
		}
	}
	wg.Wait()
	return map[string]interface{}{
		"success": true,
		"statusCode": 200,
		"message": "Done",
	}
}

func power(on bool, vsi VirtualServer, slUsername string, slApiKey string) {
	url := "https://api.softlayer.com/rest/v3.1/SoftLayer_Virtual_Guest/" + strconv.Itoa(vsi.ID) + "/powerOff"
	if on {
		url = "https://api.softlayer.com/rest/v3.1/SoftLayer_Virtual_Guest/" + strconv.Itoa(vsi.ID) + "/powerOn"
	}
	body, err := httpRequest("GET", url, nil, slUsername, slApiKey)
	if err != nil {
		log.Fatalf("Error to execute the Softlayer Virtual Guest API call to power on/off: %v", err)
	}

	fmt.Printf("VSI = %s; BODY = %s\n", vsi.DomainName, string(body))
	defer wg.Done()
}

func httpRequest(method string, url string, body io.Reader, username string, password string) (bodyResponse []byte, errorMessage error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatalf("Error to create a new HTTP Request: %v", err)
		errorMessage = err
		return
	}

	req.SetBasicAuth(username, password)
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error to execute the HTTP Request: %v", err)
		errorMessage = err
		return
	}

	byteBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error to read Body object from HTTP Request: %v", err)
		errorMessage = err
		return
	}

	bodyResponse = byteBody
	return
}
