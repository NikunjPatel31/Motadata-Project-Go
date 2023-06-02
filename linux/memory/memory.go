package memory

import (
	"MotadataPlugin/linux/util"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
)

func GetStat(connection *ssh.Client) (statistics map[string]any, err error) {

	session, err := connection.NewSession()

	if err != nil {
		return
	}

	memoryStat, err := session.Output("free -b | awk 'NR>1 {print $2\" \" $3\" \" ((($2 - $7) * 100) / $2) \" \" $4 \" \" (($4 * 100) / $2) \" \" $7}'| head -n 1|tr '\\n' \" \" && free -b | awk 'NR>2 {print $2}'")

	if err != nil {
		return
	}

	memorySplit := strings.Split(string(memoryStat), " ")

	totalMemory, err := strconv.Atoi(memorySplit[0])

	if err != nil {
		return
	}

	statistics = make(map[string]any)

	statistics[util.SystemMemoryInstalledBytes] = totalMemory / 8

	statistics[util.SystemMemoryUsedBytes] = memorySplit[1]

	statistics[util.SystemMemoryUsedPercentage] = memorySplit[2]

	statistics[util.SystemMemoryFreeBytes] = memorySplit[3]

	statistics[util.SystemMemoryFreePercentage] = memorySplit[4]

	statistics[util.SystemMemoryAvailableBytes] = memorySplit[5]

	statistics[util.SystemMemorySwapBytes] = strings.TrimSpace(memorySplit[6])

	statistics[util.SystemMemoryTotalBytes] = memorySplit[0]

	return
}
