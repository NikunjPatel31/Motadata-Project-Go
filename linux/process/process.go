package process

import (
	"MotadataPlugin/linux/util"
	"errors"
	"golang.org/x/crypto/ssh"
	"strings"
)

func GetStat(connection *ssh.Client) (statistics map[string]interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(e.(string))
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

	processSplit := strings.Split(strings.TrimSpace(strings.ReplaceAll(string(processStat), "\n", " ")), " ")

	statistics = make(map[string]interface{})

	var processes []map[string]interface{}

	for index := 0; index < len(processSplit); index++ {

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
	}

	statistics[util.SystemProcess] = processes

	return
}
