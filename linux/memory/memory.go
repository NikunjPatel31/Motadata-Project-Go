package memory

import (
	"MotadataPlugin/linux/util"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
)

func GetStat(connection *ssh.Client) (response map[string]interface{}, err error) {

	response = make(map[string]interface{})

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()

	session, err := connection.NewSession()

	if err != nil {
		return
	}

	memoryStat, err := session.Output("free -b | awk 'NR>1 {print $2\" \" $3\" \" ((($2 - $7) * 100) / $2) \" \" $4 \" \" (($4 * 100) / $2) \" \" $7}'| head -n 1|tr '\\n' \" \" && free -b | awk 'NR>2 {print $2}'")

	if err != nil {
		return
	}

	memorySplit := strings.Split(string(memoryStat), util.SpaceSeparator)

	totalMemory, err := strconv.Atoi(memorySplit[0])

	if err != nil {
		return
	}

	response = make(map[string]any)

	response[util.SystemMemoryInstalledBytes] = totalMemory / 8

	response[util.SystemMemoryUsedBytes] = memorySplit[1]

	response[util.SystemMemoryUsedPercentage] = memorySplit[2]

	response[util.SystemMemoryFreeBytes] = memorySplit[3]

	response[util.SystemMemoryFreePercentage] = memorySplit[4]

	response[util.SystemMemoryAvailableBytes] = memorySplit[5]

	response[util.SystemMemorySwapBytes] = strings.TrimSpace(memorySplit[6])

	response[util.SystemMemoryTotalBytes] = memorySplit[0]

	return
}
