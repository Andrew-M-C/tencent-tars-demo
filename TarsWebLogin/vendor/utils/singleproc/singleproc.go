package singleproc
import (
	"github.com/Andrew-M-C/go-tools/log"
	"github.com/9466/daemon"
	"io/ioutil"
	"utils/config"
	"strconv"
	"syscall"
	"os"
	"fmt"
	"strings"
)

func GetRunningProcessID() (int, error) {
	file_pach := config.PID_FILE_PATH
	pid_bytes, err := ioutil.ReadFile(file_pach)
	if err != nil {
		err_str := err.Error()
		if err_str == "no such file or directory" {
			return 0, nil
		} else {
			log.Error("Failed to get pid file: %s", err.Error())
			return 0, err
		}
	}
	pid_str := string(pid_bytes)
	pid_str = strings.Trim(pid_str, " \t\r\n")
	// log.Debug("Get pid string: >> %s <<", pid_str)

	// parse pid
	pid, err := strconv.ParseInt(pid_str, 10, 64)
	if err != nil {
		log.Error("Failed to parse config: %s", err.Error())
		return 0, err
	}
	if pid <= 0 {
		log.Error("Invalid pid string: >> %s <<", pid_str)
		return 0, nil
	}

	// check if process exists
	// reference: https://stackoverflow.com/questions/15204162/check-if-a-process-exists-in-go-way
	process, err := os.FindProcess(int(pid))
	if err != nil {
		log.Error("Failed to find process %d: %s", pid, err.Error())
		return 0, nil
	}

	err = process.Signal(syscall.Signal(0))
	if err != nil {
		// log.Info("got process %d status: %s", pid, err.Error())
		return 0, nil
	} else {
		log.Info("process %d is running", pid)
		return int(pid), nil
	}
}

func DaemonizeAndLogPid() (int, error) {
	daemon.Daemon(0, 1)

	pid := os.Getpid()
	file_pach := config.PID_FILE_PATH
	file_content := fmt.Sprintf("%d     ", pid)
	err := ioutil.WriteFile(file_pach, []byte(file_content), 0644)
	if err != nil {
		log.Error("Failed to write pid file: %s", err.Error())
	}

	return int(pid), nil
}
