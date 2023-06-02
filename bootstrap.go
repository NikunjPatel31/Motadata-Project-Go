package main

import (
	"MotadataPlugin/linux"
	"encoding/json"
	"fmt"
	"os"
)

func main() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic: ", err)
		}
	}()

	input := make(map[string]any)

	fmt.Println(os.Args[1])

	err := json.Unmarshal([]byte(os.Args[1]), &input)

	if err != nil {
		return
	}

	switch fmt.Sprint(input["Operation"]) {
	case "Discovery":
		discoveryResponse, err := linux.Discovery(input["credentialProfile"].(map[string]any), input["discoveryProfile"].(map[string]any))

		if err != nil {
			fmt.Println(discoveryResponse["message"])
		}

		response, err := json.Marshal(discoveryResponse)

		if err != nil {
			return
		}

		fmt.Println(string(response))
	case "Collect":
		//response :=
		collectResponse, err := linux.Collect(input["credentialProfile"].(map[string]any), input["discoveryProfile"].(map[string]any), input["matrices"].([]any))

		response, err := json.Marshal(collectResponse)

		if err != nil {
			fmt.Println("Error in json.Marshal")
		}

		fmt.Println(string(response))
	}

	//fmt.Println(fmt.Sprint(response["message"]))

}

/*

'{"Operation":"Discovery","credentialProfile":{"username":"nikunj",
"password":"Rajv@123"},
"discoveryProfile":{"ip":"10.20.40.197",
"port":"1267"}}'
*/

/*

'{
	"Operation":"Collect",
	"credentialProfile":
		{
			"username":"nikunj",
			"password":"Rajv@123"
		},
	"discoveryProfile":
		{
			"ip":"10.20.40.197",
			"port":"1267"
		}
	"matrices":["memory", "cpu"]
}'

*/
