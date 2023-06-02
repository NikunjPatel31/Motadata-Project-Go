package linux

import (
	"MotadataPlugin/linux/cpu"
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

			err = errors.New(e.(string))

			response[util.Status] = "error"

			response[util.Message] = err.Error()
		}
	}()

	connection, err := connect(credentialProfile, discoveryProfile)

	_, _ = connection, err

	if err != nil {
		response[util.Status] = "error"

		response[util.Message] = err.Error()

		return
	}

	defer func() {
		e := connection.Close()
		if e != nil {
			err = e
		}
	}()

	session, err := connection.NewSession()

	_ = session

	if err != nil {
		response[util.Status] = "error"

		response[util.Message] = err.Error()

		return
	}

	command := "hostname"

	result, err := session.Output(command)

	_ = result

	if err != nil {
		response[util.Status] = "error"

		response[util.Message] = err.Error()

		return
	}

	response[util.Status] = "ok"

	response[util.Message] = "Connection established"

	return
}

func Collect(credentialProfile map[string]interface{}, discoveryProfile map[string]interface{}, matrices []interface{}) (response map[string]interface{}, err error) {

	response = make(map[string]interface{})

	defer func() {
		if e := recover(); e != nil {
			// send response from here
			//err = e.(error)
			err = errors.New(e.(string))
		}
	}()

	connection, e := connect(credentialProfile, discoveryProfile)

	if e != nil {
		response[util.Status] = "error"

		response[util.Message] = e.Error()

		fmt.Println("wefe")
		return
	}

	defer func() {
		e := connection.Close()
		if e != nil {
			err = e
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

	wg.Add(4)

	time.Sleep(time.Second * 10)

	// memory statistics
	go func() {
		memoryStat, err := memory.GetStat(connection)

		if err != nil {
			response[util.Status] = "error"

			return
		}

		for key, value := range memoryStat {
			response[key] = value
		}

		wg.Done()
	}()

	// process statistics
	go func() {
		processStat, err := process.GetStat(connection)

		if err != nil {
			response[util.Status] = "error"

			return
		}

		response[util.SystemProcess] = processStat[util.SystemProcess]

		wg.Done()
	}()

	// cpu statistics
	go func() {
		cpuStat, err := cpu.GetStat(connection)

		if err != nil {
			response[util.Status] = "error"

			return
		}

		response[util.SystemCPU] = cpuStat[util.SystemCPU]

		response[util.SystemCPUCore] = cpuStat[util.SystemCPUCore]

		response[util.SystemCPUIdlePercentage] = cpuStat[util.SystemCPUIdlePercentage]

		response[util.SystemCPUUserPercentage] = cpuStat[util.SystemCPUUserPercentage]

		response[util.SystemCPUPercentage] = cpuStat[util.SystemCPUPercentage]

		wg.Done()
	}()

	// system statistics
	go func() {
		systemStat, err := system.GetStat(connection)

		_ = systemStat

		if err != nil {
			response[util.Status] = "error"

			return
		}

		for key, value := range systemStat {
			response[key] = value
		}

		wg.Done()
	}()

	wg.Wait()

	return
}

func connect(credentialProfile map[string]interface{}, discoveryProfile map[string]interface{}) (connection *ssh.Client, err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(e.(string))
		}
	}()

	config := &ssh.ClientConfig{

		User: fmt.Sprint(credentialProfile["username"]),

		Auth: []ssh.AuthMethod{ssh.Password(fmt.Sprint(credentialProfile["password"]))},

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		Timeout: time.Second * 4,
	}

	ip := fmt.Sprint(discoveryProfile["ip"], ":", fmt.Sprint(discoveryProfile["port"]))

	connection, err = ssh.Dial("tcp", ip, config)

	_ = connection

	if err != nil {

		return

	}
	return
}
