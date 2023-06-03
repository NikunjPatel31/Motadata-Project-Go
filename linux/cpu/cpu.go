package cpu

import (
	"MotadataPlugin/linux/util"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
)

func GetStat(connection *ssh.Client) (response map[string]interface{}, err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()

	response = make(map[string]interface{})

	session, err := connection.NewSession()

	if err != nil {
		return
	}

	cpuStat, err := session.Output("nproc --all && mpstat -P ALL | awk 'NR>3 {print $4 \" \" $7 \" \" $5 \" \" $14}'")

	if err != nil {
		return
	}

	cpuSplit := strings.Split(strings.TrimSpace(string(cpuStat)), util.NewLineSeparator)

	response[util.SystemCPUCore] = cpuSplit[0]

	var CPUs []map[string]string

	for index := 1; index < len(cpuSplit); index++ {

		cpu := make(map[string]string)

		CPUCore := strings.Split(cpuSplit[index], util.SpaceSeparator)

		if CPUCore[0] == "all" {

			response[util.SystemCPUPercentage] = CPUCore[1]

			response[util.SystemCPUUserPercentage] = CPUCore[2]

			response[util.SystemCPUIdlePercentage] = CPUCore[3]

		} else {

			cpu[util.SystemCPUCore] = CPUCore[0]

			cpu[util.SystemCPUPercentage] = CPUCore[1]

			cpu[util.SystemCPUUserPercentage] = CPUCore[2]

			cpu[util.SystemCPUIdlePercentage] = CPUCore[3]

			CPUs = append(CPUs, cpu)
		}
	}

	response[util.SystemCPU] = CPUs

	return
}
