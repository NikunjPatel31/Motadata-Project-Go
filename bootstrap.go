package main

import (
	"MotadataPlugin/linux"
	"encoding/json"
	"fmt"
	"os"
)

func main() {

	input := make(map[string]any)

	fmt.Println(os.Args[1])

	err := json.Unmarshal([]byte(os.Args[1]), &input)

	if err != nil {
		return
	}

	switch fmt.Sprint(input["Operation"]) {
	case "Discovery":
		fmt.Println("We are inside operation")
		discoveryResponse, err := linux.Discovery(input["credentialProfile"].(map[string]any), input["discoveryProfile"].(map[string]any))

		if err != nil {
			fmt.Println(discoveryResponse["message"])
		}

		response, err := json.Marshal(discoveryResponse)

		if err != nil {
			return
		}

		fmt.Println("Response: ", string(response))
	case "Collect":
		collectResponse, err := linux.Collect(input["credentialProfile"].(map[string]any), input["discoveryProfile"].(map[string]any), input["matrices"].([]any))

		if err != nil {
			fmt.Println("Error in collect")
		}

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
