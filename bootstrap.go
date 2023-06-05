package main

import (
	"MotadataPlugin/linux"
	"MotadataPlugin/linux/util"
	"encoding/json"
	"fmt"
	"os"
)

func main() {

	var response []byte

	input := make(map[string]interface{})

	defer func() {
		if err := recover(); err != nil {
			input[util.Message] = fmt.Sprintf("%v", err)
		}

		response, _ = json.Marshal(input)

		fmt.Println(string(response))
	}()

	err := json.Unmarshal([]byte(os.Args[1]), &input)

	if err != nil {

		input[util.Input] = os.Args[1]

		input[util.Message] = fmt.Sprintf("%v", err)

		return
	}

	switch fmt.Sprint(input["Operation"]) {

	case "Discovery":
		discoveryResponse, err := linux.Discovery(input["credentialProfile"].(map[string]interface{}), input["discoveryProfile"].(map[string]interface{}))

		if err != nil {

			input[util.Status] = util.Fail

			input[util.Message] = fmt.Sprintf("%v", err)

			return
		}

		input[util.Result] = discoveryResponse

	case "Polling":
		//response :=
		collectResponse, err := linux.Collect(input["credentialProfile"].(map[string]interface{}), input["discoveryProfile"].(map[string]interface{}), input["matrices"].([]interface{}))

		if err != nil {
			input[util.Message] = fmt.Sprintf("%v", err)

			return
		}

		input[util.Result] = collectResponse
	}
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
