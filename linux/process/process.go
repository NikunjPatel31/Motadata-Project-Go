package process

import (
	. "MotadataPlugin/linux/util"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strconv"
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

	processSplit := strings.Split(strings.TrimSpace(string(processStat)), NewLineSeparator)

	for index := 0; index < len(processSplit); index++ {

		process := strings.Split(processSplit[index], SpaceSeparator)

		processInfo := make(map[string]interface{})

		if pid, err := strconv.Atoi(process[0]); err == nil {

			processInfo[SystemProcessPid] = pid

		}

		index++

		if processCPU, err := strconv.ParseFloat(strings.ReplaceAll(process[1], "%", ""), 64); err == nil {

			processInfo[SystemProcessCPU] = processCPU

		}

		index++

		if processMem, err := strconv.ParseFloat(strings.ReplaceAll(process[2], "%", ""), 64); err == nil {

			processInfo[SystemProcessMemory] = processMem

		}

		index++

		processInfo[SystemProcessUser] = process[3]

		index++

		processInfo[SystemProcessCommand] = process[4]

		processes = append(processes, processInfo)
	}

	response[SystemProcess] = processes

	return
}
