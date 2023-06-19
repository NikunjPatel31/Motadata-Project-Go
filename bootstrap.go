package main

import (
	"MotadataPlugin/linux"
	. "MotadataPlugin/linux/util"
	"encoding/json"
	"fmt"
	"os"
)

func main() {

	var response []byte

	input := make(map[string]interface{})

	defer func() {
		if err := recover(); err != nil {
			input[Message] = fmt.Sprintf("%v", err)
		}

		response, _ = json.Marshal(input)

		fmt.Println(string(response))
	}()

	err := json.Unmarshal([]byte(os.Args[1]), &input)

	if err != nil {

		input[Message] = fmt.Sprintf("%v", err)

		return
	}

	switch fmt.Sprint(input["Operation"]) {

	case "Discovery":
		discoveryResponse, err := linux.Discovery(input["credentialProfile"].(map[string]interface{}), input["discoveryProfile"].(map[string]interface{}))

		if err != nil {

			errorResponse := make(map[string]interface{})

			errorResponse[Status] = Fail

			errorResponse[Message] = fmt.Sprintf("%v", err)

			input[Result] = errorResponse

			return
		}

		input[Result] = discoveryResponse

	case "Polling":

		collectResponse, err := linux.Collect(input["credentialProfile"].(map[string]interface{}), input["discoveryProfile"].(map[string]interface{}), fmt.Sprintf("%v", input["metrics"]))

		if err != nil {

			var errorResult = make(map[string]string)

			errorResult[Message] = fmt.Sprintf("%v", err)

			input[Error] = errorResult

			return
		}

		input[Result] = collectResponse

	default:
		var errorResult = make(map[string]string)

		errorResult[Message] = "invalid operation"

		input[Error] = errorResult

	}
}
