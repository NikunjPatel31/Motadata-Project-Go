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

	if totalMemory, err := strconv.Atoi(memorySplit[0]); err == nil {

		response[util.SystemMemoryInstalledBytes] = totalMemory / 8

	}

	if totalBytes, err := strconv.Atoi(memorySplit[0]); err == nil {

		response[util.SystemMemoryTotalBytes] = totalBytes

	}

	if usedBytes, err := strconv.Atoi(memorySplit[1]); err == nil {

		response[util.SystemMemoryUsedBytes] = usedBytes

	}

	if usedPer, err := strconv.ParseFloat(memorySplit[2], 64); err == nil {

		response[util.SystemMemoryUsedPercentage] = usedPer

	}

	if freeBytes, err := strconv.Atoi(memorySplit[3]); err == nil {

		response[util.SystemMemoryFreeBytes] = freeBytes

	}

	if freePer, err := strconv.ParseFloat(memorySplit[4], 64); err == nil {

		response[util.SystemMemoryFreePercentage] = freePer

	}

	if availableBytes, err := strconv.Atoi(memorySplit[5]); err == nil {

		response[util.SystemMemoryAvailableBytes] = availableBytes

	}

	response[util.SystemMemorySwapBytes] = strings.TrimSpace(memorySplit[6])

	return
}
