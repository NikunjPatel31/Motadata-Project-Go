package main

import (
	"MotadataPlugin/linux"
	"encoding/json"
	"fmt"
	"os"
)

func main() {

	var response []byte

	_ = response

	input := make(map[string]any)

	defer func() {
		if err := recover(); err != nil {
			input["error"] = fmt.Sprintf("%v", err)
		}

		response, _ = json.Marshal(input)

		fmt.Println(string(response))
	}()

	fmt.Println(os.Args[1])

	err := json.Unmarshal([]byte(os.Args[1]), &input)

	if err != nil {

		input["error"] = fmt.Sprintf("%v", err)

		return
	}

	switch fmt.Sprint(input["Operation"]) {

	case "Discovery":
		discoveryResponse, err := linux.Discovery(input["credentialProfile"].(map[string]any), input["discoveryProfile"].(map[string]any))

		if err != nil {

			input["error"] = fmt.Sprintf("%v", err)

			return
		}

		input["result"] = discoveryResponse

	case "Polling":
		//response :=
		collectResponse, err := linux.Collect(input["credentialProfile"].(map[string]any), input["discoveryProfile"].(map[string]any), input["matrices"].([]any))

		if err != nil {
			input["error"] = fmt.Sprintf("%v", err)

			return
		}

		input["result"] = collectResponse
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
