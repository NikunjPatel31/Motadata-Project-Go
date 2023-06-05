package disk

import (
	"MotadataPlugin/linux/util"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
)

const (
	cmd = "iostat -x | awk 'NR>6 {print $1 \" \" $3 \" \" $2 \" \" $8 \" \" $9}'"

	SystemDiskBytesPerSec = "system.disk.bytes.per.sec"

	SystemDiskReadBytesPerSec = "system.disk.read.bytes.per.sec"

	SystemDiskReadOpsPerSec = "system.disk.read.ops.per.sec"

	SystemDiskWriteBytesPerSec = "system.disk.write.bytes.per.sec"

	SystemDiskWriteOpsPerSec = "system.disk.write.ops.per.sec"
)

func GetStat(connection *ssh.Client) (response map[string]interface{}, err error) {

	response = map[string]interface{}{}

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()

	session, err := connection.NewSession()

	if err != nil {
		return
	}

	diskStat, err := session.Output(cmd)

	if err != nil {
		return
	}

	var allDisk []map[string]interface{}

	cmdInfo := strings.Split(strings.TrimSpace(string(diskStat)), util.NewLineSeparator)

	for index := 0; index < len(cmdInfo); index++ {

		disk := make(map[string]interface{})

		diskInfo := strings.Split(cmdInfo[index], util.SpaceSeparator)

		disk[util.SystemDisk] = diskInfo[0]

		var readBytes float64

		if readBytes, err = strconv.ParseFloat(diskInfo[1], 64); err == nil {

			disk[SystemDiskReadBytesPerSec] = readBytes * 1000.0

		}

		if readOpsPerSec, err := strconv.ParseFloat(diskInfo[2], 64); err == nil {

			disk[SystemDiskReadOpsPerSec] = readOpsPerSec

		}

		var writeBytes float64

		if writeBytes, err = strconv.ParseFloat(diskInfo[3], 64); err == nil {

			disk[SystemDiskWriteBytesPerSec] = writeBytes * 1000.0

		}

		if writeOpsPerSec, err := strconv.ParseFloat(diskInfo[4], 64); err == nil {

			disk[SystemDiskWriteOpsPerSec] = writeOpsPerSec

		}

		disk[SystemDiskBytesPerSec] = (writeBytes + readBytes) * 1000.0

		allDisk = append(allDisk, disk)
	}

	response[util.SystemDisk] = allDisk

	return

}
