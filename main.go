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
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	params := map[string]interface{}{
		"username": os.Getenv("SOFTLAYER_USERNAME"),
		"apikey": os.Getenv("SOFTLAYER_APIKEY"),
		"name": os.Getenv("VSIS_NAME"),
		"power": os.Getenv("POWER"),
	}
	res := Main(params)
	b, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("Error to parse Function response to []byte: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.Write(b)
}

func main() {
	godotenv.Load()

	http.HandleFunc("/", Handler)
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
