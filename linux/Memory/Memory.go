package Memory

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
)

func GetStat(connection *ssh.Client) (statistics map[string]any, err error) {
	defer func() {
		if err := recover(); err != nil {
			// send appropriate response from here
			fmt.Println("Error in getStat: ", err.(error).Error())
		}
	}()

	session, err := connection.NewSession()

	if err != nil {
		panic(err)
	}

	memoryStat, err := session.Output("free -b | awk 'NR>1 {print $2\" \" $3\" \" ((($2 - $7) * 100) / $2) \" \" $4 \" \" (($4 * 100) / $2) \" \" $7}'| head -n 1|tr '\\n' \" \" && free -b | awk 'NR>2 {print $2}'")

	if err != nil {
		panic(err)
	}

	fmt.Println("Memory stat: ", string(memoryStat))

	memorySplit := strings.Split(string(memoryStat), " ")

	totalMemory, err := strconv.Atoi(memorySplit[0])

	if err != nil {
		panic(err)
	}

	statistics = make(map[string]any)

	statistics["system.memory.installed"] = totalMemory / 8

	statistics["system.memory.used.bytes"] = memorySplit[1]

	statistics["system.memory.used.percentage"] = memorySplit[2]

	statistics["system.memory.free.bytes"] = memorySplit[3]

	statistics["system.memory.free.percentage"] = memorySplit[4]

	statistics["system.memory.available.bytes"] = memorySplit[5]

	statistics["system.memory.swap"] = memorySplit[6]

	statistics["system.memory.total.bytes"] = memorySplit[0]

	return
}
