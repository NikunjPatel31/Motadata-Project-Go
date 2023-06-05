package cpu

import (
	. "MotadataPlugin/linux/util"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
)

const cmd = "nproc --all && mpstat -P ALL | awk 'NR>3 {print $4 \" \" $7 \" \" $5 \" \" $14}'"

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

	cpuStat, err := session.Output(cmd)

	if err != nil {
		return
	}

	cpuSplit := strings.Split(strings.TrimSpace(string(cpuStat)), NewLineSeparator)

	response[SystemCPUCore] = cpuSplit[0]

	var CPUs []map[string]interface{}

	for index := 1; index < len(cpuSplit); index++ {

		cpu := make(map[string]interface{})

		CPUCore := strings.Split(cpuSplit[index], SpaceSeparator)

		if CPUCore[0] == "all" {

			if cpuPer, err := strconv.ParseFloat(CPUCore[1], 64); err == nil {

				response[SystemCPUPercentage] = cpuPer

			}

			if percentage, err := strconv.ParseFloat(CPUCore[2], 64); err == nil {

				response[SystemCPUUserPercentage] = percentage

			}

			if cpuPer, err := strconv.ParseFloat(CPUCore[3], 64); err == nil {

				response[SystemCPUIdlePercentage] = cpuPer

			}

		} else {

			if cpuNo, err := strconv.Atoi(CPUCore[0]); err == nil {

				cpu[SystemCPUCore] = cpuNo

			}

			if percentage, err := strconv.ParseFloat(CPUCore[1], 64); err == nil {

				cpu[SystemCPUPercentage] = percentage

			}

			if percentage, err := strconv.ParseFloat(CPUCore[2], 64); err == nil {

				cpu[SystemCPUUserPercentage] = percentage

			}

			if percentage, err := strconv.ParseFloat(CPUCore[3], 64); err == nil {

				cpu[SystemCPUIdlePercentage] = percentage

			}

			CPUs = append(CPUs, cpu)
		}
	}

	response[SystemCPU] = CPUs

	return
}
