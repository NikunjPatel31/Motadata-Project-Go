package linux

import (
	"MotadataPlugin/linux/cpu"
	"MotadataPlugin/linux/disk"
	"MotadataPlugin/linux/memory"
	"MotadataPlugin/linux/process"
	"MotadataPlugin/linux/system"
	. "MotadataPlugin/linux/util"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func Discovery(credentialProfile map[string]interface{}, discoveryProfile map[string]interface{}) (response map[string]interface{}, err error) {
	response = make(map[string]interface{})

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))

		}
	}()

	connection, err := Connect(credentialProfile, discoveryProfile)

	if err != nil {

		return
	}

	defer func() {
		e := connection.Close()

		if e != nil {

			err = errors.New(fmt.Sprintf("%v", e))

			errorResponse := make(map[string]string)

			errorResponse[Status] = Fail

			errorResponse[Message] = fmt.Sprintf("%v", e)

			response[Result] = errorResponse
		}
	}()

	session, err := connection.NewSession()

	if err != nil {

		errorResponse := make(map[string]string)

		errorResponse[Status] = Fail

		errorResponse[Message] = fmt.Sprintf("%v", err)

		response[Result] = errorResponse

		return
	}

	command := "hostname"

	discoveryOutput, err := session.Output(command)

	if err != nil {

		errorResponse := make(map[string]string)

		errorResponse[Status] = Fail

		errorResponse[Message] = fmt.Sprintf("%v", err)

		response[Result] = errorResponse

		return
	}

	response[Status] = Success

	response[Message] = strings.TrimSpace(string(discoveryOutput))

	return
}

func Collect(credentialProfile map[string]interface{}, discoveryProfile map[string]interface{}, metrics string) (response map[string]interface{}, err error) {

	response = make(map[string]interface{})

	defer func() {
		if e := recover(); e != nil {
			// send response from here
			//err = e.(error)
			err = errors.New(fmt.Sprintf("%v", e))

			response[Message] = fmt.Sprintf("%v", e)
		}
	}()

	switch metrics {

	case CPU:

		connection, e := Connect(credentialProfile, discoveryProfile)

		if e != nil {

			err = e

			return
		}

		defer func() {
			err = connection.Close()
		}()

		cpuStat, e := cpu.GetStat(connection)

		if e != nil {

			err = e

			return
		}

		response[SystemCPU] = cpuStat[SystemCPU]

		response[SystemCPUCore] = cpuStat[SystemCPUCore]

		response[SystemCPUIdlePercentage] = cpuStat[SystemCPUIdlePercentage]

		response[SystemCPUUserPercentage] = cpuStat[SystemCPUUserPercentage]

		response[SystemCPUPercentage] = cpuStat[SystemCPUPercentage]

		return
	case Memory:

		connection, e := Connect(credentialProfile, discoveryProfile)

		if e != nil {

			err = e

			return
		}

		defer func() {
			err = connection.Close()
		}()

		memoryStat, e := memory.GetStat(connection) // shadow

		if e != nil {

			err = e

			return
		}

		for key, value := range memoryStat {

			response[key] = value

		}

		return
	case Process:

		connection, e := Connect(credentialProfile, discoveryProfile)

		if e != nil {

			err = e

			return
		}

		defer func() {
			e = connection.Close()

			if e != nil {
				err = e
			}
		}()

		processStat, e := process.GetStat(connection)

		if e != nil {

			err = e

			return
		}

		response[SystemProcess] = processStat[SystemProcess]

		return
	case System:

		connection, e := Connect(credentialProfile, discoveryProfile)

		if e != nil {

			err = errors.New(fmt.Sprintf("%v", e))

			return
		}

		defer func() {
			err = connection.Close()
		}()

		systemStat, e := system.GetStat(connection)

		if e != nil {
			err = errors.New(fmt.Sprintf("%v", e))

			return
		}

		for key, value := range systemStat {
			response[key] = value
		}

		return
	case Disk:

		connection, e := Connect(credentialProfile, discoveryProfile)

		if e != nil {

			err = e

			return
		}

		defer func() {
			err = connection.Close()
		}()

		diskStat, e := disk.GetStat(connection) // shadow

		if e != nil {

			err = e

			return
		}

		response[SystemDisk] = diskStat[SystemDisk]

		return

	}

	return
}

func Connect(credentialProfile map[string]interface{}, discoveryProfile map[string]interface{}) (connection *ssh.Client, err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()

	config := &ssh.ClientConfig{

		User: fmt.Sprint(credentialProfile["username"]),

		Auth: []ssh.AuthMethod{ssh.Password(fmt.Sprint(credentialProfile["password"]))},

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		Timeout: time.Second * 20,
	}

	ip := fmt.Sprint(discoveryProfile["ip"], ":", discoveryProfile["port"])

	connection, err = ssh.Dial("tcp", ip, config)

	return
}
