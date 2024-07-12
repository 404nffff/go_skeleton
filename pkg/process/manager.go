package process

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/sevlyar/go-daemon"
)

var validCommands = []string{"start", "start debug", "stop", "restart"}

// GetValidCommands 返回有效的命令列表
func getValidCommands() []string {
	return validCommands
}

// isValidCommand 检查提供的命令是否在有效命令列表中
func isValidCommand(cmd string, validCommands []string) bool {
	for _, validCommand := range validCommands {
		if cmd == validCommand {
			return true
		}
	}
	return false
}

// getPIDFile 根据服务名称生成 PID 文件路径
func getPIDFile(service string) string {
	return fmt.Sprintf("./logs/%s.pid", service)
}

// getLogFile 根据服务名称生成日志文件路径
func getLogFile(service string) string {
	return fmt.Sprintf("./logs/%s.log", service)
}

// setProcessName 设置进程名称
func setProcessName(name string) error {
	// nameBytes := []byte(name)
	// if len(nameBytes) >= 16 {
	// 	return fmt.Errorf("name too long")
	// }

	// mib := []int32{1, 12, 9}
	// _, _, errno := syscall.Syscall6(
	// 	syscall.SYS___SYSCTL,
	// 	uintptr(unsafe.Pointer(&mib[0])),
	// 	3,
	// 	uintptr(unsafe.Pointer(&nameBytes[0])),
	// 	uintptr(len(nameBytes)),
	// 	0,
	// 	0,
	// )
	// if errno != 0 {
	// 	return errno
	// }
	return nil
}

// StartServerInBackground 启动后台服务器
// 使用守护进程的方式启动服务并记录日志和 PID 文件
func startServerInBackground(pidFile, logFile, processName string, startFunc func()) {
	cntxt := &daemon.Context{
		PidFileName: pidFile,
		PidFilePerm: 0644,
		LogFileName: logFile,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}

	d, err := cntxt.Reborn()
	if err != nil {
		panic(fmt.Sprintf("Unable to run: %s", err))
	}
	if d != nil {
		return
	}
	defer func() {
		err := cntxt.Release()
		if err != nil {
			panic(fmt.Sprintf("Unable to release context: %s", err))
		}
	}()
	log.Print("Daemon started")

	// 设置进程名称
	// if err := setProcessName(processName); err != nil {
	// 	panic(fmt.Sprintf("Failed to set process name: %v", err)
	// }

	startFunc()

	err = daemon.ServeSignals()
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	log.Print("Daemon terminated")
}

// StopServer 停止运行的服务器
// 通过 PID 文件找到守护进程并发送停止信号
func stopServer(pidFile string) {
	cntxt := &daemon.Context{PidFileName: pidFile}
	d, err := cntxt.Search()
	if err != nil {
		panic(fmt.Sprintf("Unable to find the daemon: %s", err))
	}

	if d != nil {
		err := d.Signal(syscall.SIGTERM)
		if err != nil {
			panic(fmt.Sprintf("Unable to send signal to the daemon: %s", err))
		}

		log.Print("Daemon stopped")
	} else {
		log.Print("No daemon found")
	}
}

// RestartServer 重启服务器
// 先停止当前运行的服务器，然后重新启动
func restartServer(pidFile, logFile, processName string, startFunc func()) {
	stopServer(pidFile)
	time.Sleep(2 * time.Second)
	startServerInBackground(pidFile, logFile, processName, startFunc)
}

// 初始化服务管理器
func Initialize(service string, startFunc func()) {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <command>")
		fmt.Println("Commands:", validCommands)
		os.Exit(1)
	}

	command := os.Args[1]

	if len(os.Args) > 2 {
		command += " " + os.Args[2]
	}

	if !isValidCommand(command, validCommands) {
		fmt.Println("Invalid command. Use", validCommands)
		os.Exit(1)
	}

	pidFile := getPIDFile(service)
	logFile := getLogFile(service)
	processName := service // 使用服务名称作为进程名称

	switch command {
	case "start":
		startServerInBackground(pidFile, logFile, processName, startFunc)
	case "start debug":
		startFunc()
	case "stop":
		stopServer(pidFile)
	case "restart":
		restartServer(pidFile, logFile, processName, startFunc)
	default:
		panic(fmt.Sprintf("Invalid command. Use start, start debug, stop, or restart."))
	}
}
