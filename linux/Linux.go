package linux

import (
	"MotadataPlugin/linux/Memory"
	"fmt"
	"golang.org/x/crypto/ssh"
)

func Discovery(credentialProfile map[string]any, discoveryProfile map[string]any) (response map[string]any, err error) {

	defer func() {
		if err := recover(); err != nil {
			response = make(map[string]any)

			response["status"] = "error"

			response["message"] = err
		}
	}()

	connection, err := connect(credentialProfile, discoveryProfile)

	_, _ = connection, err

	if err != nil {
		panic(err)
	}

	defer func() {
		err := connection.Close()
		if err != nil {
			fmt.Println("connection defer")
			panic(err)
		}
	}()

	session, err := connection.NewSession()

	_ = session

	if err != nil {
		panic(err)
	}

	command := "hostname"

	result, err := session.Output(command)

	_ = result

	if err != nil {
		panic(err)
	}

	response = make(map[string]any)

	response["status"] = "ok"

	response["message"] = "Connection established"

	return
}

func Collect(credentialProfile map[string]any, discoveryProfile map[string]any, matrices []any) (response map[string]any, err error) {

	defer func() {
		if err := recover(); err != nil {
			// send response from here
			fmt.Println("Error in collect: ", err.(error).Error())
		}
	}()

	connection, err := connect(credentialProfile, discoveryProfile)

	if err != nil {
		panic(err)
	}

	defer func() {
		err := connection.Close()
		if err != nil {
			panic(err)
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

	// memory statics
	stat, err := Memory.GetStat(connection)

	if err != nil {
		return nil, err
	}

	_ = stat

	response = make(map[string]any)

	response["status"] = "ok"

	response["memory"] = stat

	return
}

func connect(credentialProfile map[string]any, discoveryProfile map[string]any) (connection *ssh.Client, err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("error ", err)
		}
	}()

	config := &ssh.ClientConfig{

		User: fmt.Sprint(credentialProfile["username"]),

		Auth: []ssh.AuthMethod{ssh.Password(fmt.Sprint(credentialProfile["password"]))},

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	ip := fmt.Sprint(discoveryProfile["ip"], ":", fmt.Sprint(discoveryProfile["port"]))

	fmt.Println("IP: ", ip)

	connection, err = ssh.Dial("tcp", ip, config)

	_ = connection

	if err != nil {

		panic("Connection not established")

	}

	return
}
