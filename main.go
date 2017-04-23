package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

// run as Unprivileged user
func runUnprivileged() {
	cmd := exec.Command("/proc/self/exe", append([]string{"fork"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUSER | syscall.CLONE_NEWNS,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func runFork() {
	fmt.Printf("Running %v \n", os.Args[2:])

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := syscall.Chroot("rootfs")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chdir("/")
	if err != nil {
		log.Fatal(err)
	}

	// TODO: create and mount /proc

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// tcr run <args>
func main() {
	switch os.Args[1] {
	case "run":
		runUnprivileged()
	case "fork":
		runFork()
	default:
		fmt.Println("tcr run <args>")
	}

}
