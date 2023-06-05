package system

import (
	"MotadataPlugin/linux/util"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
)

func GetStat(connection *ssh.Client) (statistics map[string]interface{}, err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()

	statistics = make(map[string]interface{})

	session, err := connection.NewSession()

	if err != nil {
		return
	}

	systemStat, err := session.Output("hostname |tr '\\n' \" \" && uname |tr '\\n' \" \" && ps -eo nlwp | awk '{ num_threads += $1 } END { print num_threads }' | tr '\\n' \" \" && vmstat | tail -n 1 | awk '{print $12}' | tr '\\n' \" \" && ps axo state | grep \"R\" | wc -l | tr '\\n' \" \" && ps axo stat | grep \"D\" | wc -l && uptime -p | awk 'gsub(\"up \",\"\")' && hostnamectl | grep \"Operating System\"")

	if err != nil {
		return
	}

	systemSplit := strings.Split(string(systemStat), util.NewLineSeparator)

	row1 := strings.Split(systemSplit[0], util.SpaceSeparator)

	statistics[util.SystemName] = row1[0]

	statistics[util.SystemOSName] = row1[1]

	if threads, err := strconv.Atoi(row1[2]); err == nil {

		statistics[util.SystemThreads] = threads

	}

	if contextSwitch, err := strconv.Atoi(row1[3]); err == nil {

		statistics[util.SystemContextSwtiches] = contextSwitch

	}

	if runningProcess, err := strconv.Atoi(row1[4]); err == nil {

		statistics[util.SystemRunningProcesses] = runningProcess

	}

	if blockedProcess, err := strconv.Atoi(row1[5]); err == nil {

		statistics[util.SystemBlockedProcesses] = blockedProcess

	}

	statistics[util.SystemUptime] = systemSplit[1]

	statistics[util.SystemOSVersion] = strings.TrimSpace(strings.Split(systemSplit[2], ":")[1])

	return
}
