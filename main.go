/**
 *
 * Copyright 2020 Victor Shinya
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
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	username string
	apikey   string
)

// VirtualServer model to handle only the useful data
type VirtualServer struct {
	ID         int    `json:"id"`
	DomainName string `json:"fullyQualifiedDomainName"`
}

// Main function to run the Action source code
func Main(params map[string]interface{}) map[string]interface{} {
	username = params["username"].(string)
	apikey = params["apikey"].(string)
	on := params["power"].(bool)
	name := params["name"].(string)

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.softlayer.com/rest/v3.1/SoftLayer_Account/getVirtualGuests", nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(username, apikey)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var vsis []VirtualServer
	err = json.Unmarshal(body, &vsis)
	r := make(map[string]interface{})
	for _, vsi := range vsis {
		// fmt.Printf("ID = %d; VSI = %s\n", vsi.Id, vsi.DomainName)
		if vsi.DomainName == name {
			success := power(on, &vsi)
			r["success"] = success
			r["id"] = vsi.ID
			r["name"] = vsi.DomainName
		}
	}
	return r
}

func power(on bool, vsi *VirtualServer) bool {
	url := "https://api.softlayer.com/rest/v3.1/SoftLayer_Virtual_Guest/" + strconv.Itoa(vsi.ID) + "/powerOff"
	if on {
		url = "https://api.softlayer.com/rest/v3.1/SoftLayer_Virtual_Guest/" + strconv.Itoa(vsi.ID) + "/powerOn"
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(username, apikey)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	res := true
	if string(body) == "false" {
		res = false
	}
	return res
}
