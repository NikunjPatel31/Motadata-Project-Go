package cpu

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

	statistics = make(map[string]interface{})

	defer func() {

		if e := recover(); e != nil {
			err = errors.New(e.(string))
		}

	}()

	session, err := connection.NewSession()

	if err != nil {
		statistics[util.Status] = "error"

		statistics[util.Message] = err.Error()

		return
	}

	cpuStat, err := session.Output("nproc --all && mpstat -P ALL | awk 'NR>3 {print $4 \" \" $7 \" \" $5 \" \" $14}'")

	if err != nil {

		statistics[util.Status] = "error"

		statistics[util.Message] = err.Error()

		return
	}

	cpuSplit := strings.Split(strings.TrimSpace(string(cpuStat)), "\n")

	statistics[util.SystemCPUCore] = cpuSplit[0]

	var CPUs []map[string]string

	for index := 1; index < len(cpuSplit); index++ {

		cpu := make(map[string]string)

		CPUCore := strings.Split(cpuSplit[index], " ")

		if CPUCore[0] == "all" {

			statistics[util.SystemCPUPercentage] = CPUCore[1]

			statistics[util.SystemCPUUserPercentage] = CPUCore[2]

			statistics[util.SystemCPUIdlePercentage] = CPUCore[3]

		} else {

			cpu[util.SystemCPUCore] = CPUCore[0]

			cpu[util.SystemCPUPercentage] = CPUCore[1]

			cpu[util.SystemCPUUserPercentage] = CPUCore[2]

			cpu[util.SystemCPUIdlePercentage] = CPUCore[3]

			CPUs = append(CPUs, cpu)
		}
	}

	statistics[util.SystemCPU] = CPUs

	return
}
