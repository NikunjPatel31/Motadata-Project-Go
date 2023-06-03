package system

import (
	"MotadataPlugin/linux/util"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
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

	statistics[util.SystemThreads] = row1[2]

	statistics[util.SystemContextSwtiches] = row1[3]

	statistics[util.SystemRunningProcesses] = row1[4]

	statistics[util.SystemBlockedProcesses] = row1[5]

	statistics[util.SystemUptime] = systemSplit[1]

	statistics[util.SystemOSVersion] = strings.TrimSpace(strings.Split(systemSplit[2], ":")[1])

	return
}
