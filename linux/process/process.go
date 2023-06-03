package process

import (
	"MotadataPlugin/linux/util"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
)

func GetStat(connection *ssh.Client) (response map[string]interface{}, err error) {

	response = make(map[string]interface{})

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()

	session, err := connection.NewSession()

	if err != nil {
		// send response from here
		return
	}

	processStat, err := session.Output("ps aux | awk 'NR>1 {print $2 \" \" $3 \"% \" $4 \"% \" $1\" \"$11}'")

	if err != nil {
		return
	}

	var processes []map[string]interface{}

	processSplit := strings.Split(strings.TrimSpace(string(processStat)), util.NewLineSeparator)

	for index := 0; index < len(processSplit); index++ {

		process := strings.Split(processSplit[index], util.SpaceSeparator)

		processInfo := make(map[string]interface{})

		processInfo[util.SystemProcessPid] = process[0]
		index++
		processInfo[util.SystemProcessCPU] = process[1]
		index++
		processInfo[util.SystemProcessMemory] = process[2]
		index++
		processInfo[util.SystemProcessUser] = process[3]
		index++
		processInfo[util.SystemProcessCommand] = process[4]

		processes = append(processes, processInfo)
	}

	//processSplit := strings.Split(strings.TrimSpace(strings.ReplaceAll(string(processStat), "\n", " ")), " ")

	/*for index := 0; index < len(processSplit); index++ {

		process := make(map[string]interface{})

		process[util.SystemProcessPid] = processSplit[index]
		index++
		process[util.SystemProcessCPU] = processSplit[index]
		index++
		process[util.SystemProcessMemory] = processSplit[index]
		index++
		process[util.SystemProcessUser] = processSplit[index]
		index++
		process[util.SystemProcessCommand] = processSplit[index]

		processes = append(processes, process)
	}*/

	response[util.SystemProcess] = processes

	return
}
