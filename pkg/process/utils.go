package process

import (
	"os"
	"strconv"
	"syscall"
)

func GetPIDFromFile(pidFile string) (int, error) {
	data, err := os.ReadFile(pidFile) // 使用 os.ReadFile 替代 ioutil.ReadFile
	if err != nil {
		return 0, err
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}

	return pid, nil
}

func IsProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))
	return err == nil
}
