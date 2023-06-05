package linux

import (
	"MotadataPlugin/linux/cpu"
	"MotadataPlugin/linux/disk"
	"MotadataPlugin/linux/memory"
	"MotadataPlugin/linux/process"
	"MotadataPlugin/linux/system"
	"MotadataPlugin/linux/util"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"sync"
	"time"
)

func Discovery(credentialProfile map[string]interface{}, discoveryProfile map[string]interface{}) (response map[string]interface{}, err error) {
	response = make(map[string]interface{})

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))

			response[util.Status] = util.Fail

			response[util.Message] = fmt.Sprintf("%v", e)
		}
	}()

	connection, err := connect(credentialProfile, discoveryProfile)

	if err != nil {
		response[util.Status] = util.Fail

		response[util.Message] = fmt.Sprintf("%v", err)

		return
	}

	defer func() {
		e := connection.Close()

		if e != nil {

			err = errors.New(fmt.Sprintf("%v", e))

			response[util.Status] = util.Fail

			response[util.Message] = fmt.Sprintf("%v", err)
		}
	}()

	session, err := connection.NewSession()

	if err != nil {

		response[util.Status] = util.Fail

		response[util.Message] = fmt.Sprintf("%v", err)

		return
	}

	command := "hostname"

	_, err = session.Output(command)

	if err != nil {

		response[util.Status] = util.Fail

		response[util.Message] = fmt.Sprintf("%v", err)

		return
	}

	response[util.Status] = util.Success

	response[util.Message] = "Connection established"

	return
}

func Collect(credentialProfile map[string]interface{}, discoveryProfile map[string]interface{}, matrices []interface{}) (response map[string]interface{}, err error) {

	response = make(map[string]interface{})

	defer func() {
		if e := recover(); e != nil {
			// send response from here
			//err = e.(error)
			err = errors.New(fmt.Sprintf("%v", e))

			response[util.Status] = util.Fail

			response[util.Message] = fmt.Sprintf("%v", e)
		}
	}()

	/*session, err := connection.NewSession()

	if err != nil {
		panic(err)
	}

	memoryStat, err := session.Output("free -b | awk 'NR>1 {print $2\" \" $3\" \" ((($2 - $7) * 100) / $2) \" \" $4 \" \" (($4 * 100) / $2) \" \" $7}'| head -n 1|tr '\\n' \" \" && free -b | awk 'NR>2 {print $2}'")

	if err != nil {
		panic(err)
	}*/

	//var channel = make(chan interface{}, 5)

	var wg sync.WaitGroup

	wg.Add(5)

	// memory statistics
	go func() {

		defer func() {
			if e := recover(); e != nil {
				err = errors.New(fmt.Sprintf("%v", e))
			}
			wg.Done()
		}()

		connection, e := connect(credentialProfile, discoveryProfile)

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

	}()

	// process statistics
	go func() {

		defer func() {
			if e := recover(); e != nil {

				err = errors.New(fmt.Sprintf("%v", e))

			}
			wg.Done()
		}()

		connection, e := connect(credentialProfile, discoveryProfile)

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

		response[util.SystemProcess] = processStat[util.SystemProcess]

	}()

	// cpu statistics
	go func() {

		defer func() {
			if e := recover(); e != nil {

				err = errors.New(fmt.Sprintf("%v", e))

			}
			wg.Done()
		}()

		connection, e := connect(credentialProfile, discoveryProfile)

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

		response[util.SystemCPU] = cpuStat[util.SystemCPU]

		response[util.SystemCPUCore] = cpuStat[util.SystemCPUCore]

		response[util.SystemCPUIdlePercentage] = cpuStat[util.SystemCPUIdlePercentage]

		response[util.SystemCPUUserPercentage] = cpuStat[util.SystemCPUUserPercentage]

		response[util.SystemCPUPercentage] = cpuStat[util.SystemCPUPercentage]

	}()

	// system statistics
	go func() {

		defer func() {
			if e := recover(); e != nil {
				err = errors.New(fmt.Sprintf("%v", e))
			}
			wg.Done()
		}()

		connection, e := connect(credentialProfile, discoveryProfile)

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

	}()

	// disk statistics
	go func() {

		defer func() {
			if e := recover(); e != nil {
				err = errors.New(fmt.Sprintf("%v", e))
			}
			wg.Done()
		}()

		connection, e := connect(credentialProfile, discoveryProfile)

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

		response[util.SystemDisk] = diskStat[util.SystemDisk]

		/*for key, value := range memoryStat {

			response[key] = value

		}*/

	}()

	wg.Wait()

	return
}

func connect(credentialProfile map[string]interface{}, discoveryProfile map[string]interface{}) (connection *ssh.Client, err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()

	config := &ssh.ClientConfig{

		User: fmt.Sprint(credentialProfile["username"]),

		Auth: []ssh.AuthMethod{ssh.Password(fmt.Sprint(credentialProfile["password"]))},

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		Timeout: time.Second * 4,
	}

	ip := fmt.Sprint(discoveryProfile["ip"], ":", discoveryProfile["port"])

	connection, err = ssh.Dial("tcp", ip, config)

	return
}
