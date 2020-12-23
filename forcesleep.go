package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func findProcesses(print bool) (err int, result []string) {
	cmd1 := "ps -ef | grep -i $(which shutdown) | grep -v grep | awk '{print $2}'"
	output, cmderr := exec.Command("bash", "-c", cmd1).Output()
	var processes []string
	if cmderr != nil {
		if print == true {
			fmt.Println("Failed when searching for processes.")
		}
		err = 1
	} else if len(string(output)) <= 0 {
		if print == true {
			fmt.Println("No shutdown process found.")
		}
		err = 1
	} else {
		processes = strings.Split(string(output), "\n")
		processes = processes[:len(processes)-1]
		if print == true {
			fmt.Println("Found", len(processes), "shutdown processes.")
		}
	}
	return err, processes
}

func cancelSleep() {
	err, processes := findProcesses(true)
	if err != 0 {
		os.Exit(0)
	}
	for _, process := range processes {
		fmt.Println("Ending process", process)
		cmd := "sudo kill -9 " + process
		_, errs := exec.Command("bash", "-c", cmd).Output()
		if errs != nil {
			fmt.Println(errs)
			os.Exit(0)
		}
	}
	err, processes = findProcesses(true)
	os.Exit(0)
}

func main() {
	var args string
	fmt.Println("")
	if len(os.Args) > 1 {
		args = os.Args[1]
	} else {
		fmt.Println("You have to specify the delay in minutes or 'cancel'.")
		os.Exit(0)
	}
	if strings.ToLower(args) == "cancel" {
		cancelSleep()
	}
	_, err := strconv.Atoi(args)
	if err != nil {
		fmt.Println("You have to provide a number of minutes.")
		os.Exit(0)
	}
	_, processes := findProcesses(false)
	if len(processes) > 0 {
		fmt.Println("")
		fmt.Println("**A shutdown process already exists. Another will be created.**")
		fmt.Println("- Use 'forcesleep cancel' to cancel.")
	}
	cmd1 := "which shutdown"
	output, cmderr := exec.Command("bash", "-c", cmd1).Output()
	outputAdj := string(output)
	outputAdj = strings.TrimSuffix(outputAdj, "\n")
	if cmderr != nil {
		fmt.Println("Error when trying to find shutdown binary.")
	}
	cmd := exec.Command("sudo", outputAdj, "-s", "+"+args)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Could not execute command.\n", err)
		os.Exit(0)
	}
}
