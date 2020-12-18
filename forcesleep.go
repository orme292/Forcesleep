package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func main() {

	var args string

	if len(os.Args) > 1 {
		args = os.Args[1]
	} else {
		fmt.Println("You have to specify the delay in minutes or 'cancel'.")
		os.Exit(0)
	}

	if args == "cancel" {
		cmd1 := "ps -ef | grep -i shutdown | grep -v grep | awk '{print $2}'"
		out, err := exec.Command("bash", "-c", cmd1).Output()
		if err != nil {
			fmt.Println("Failed.")
			os.Exit(0)
		}
		if len(string(out)) <= 0 {
			fmt.Println("No shutdown process found.")
			os.Exit(0)
		}
		fmt.Println("Ending shutdown process " + string(out))
		cmd2 := "sudo kill -9 " + string(out)
		out, err = exec.Command("bash", "-c", cmd2).Output()
		fmt.Println(string(out))
		os.Exit(0)
	}

	_, err := strconv.Atoi(args)
	if err != nil {
		fmt.Println("You have to provide a number of minutes.")
		os.Exit(0)
	}

	cmd := exec.Command("sudo", "shutdown", "-s", "+"+args)

	err = cmd.Run()
	if err != nil {
		fmt.Println("Could not execute command.\n", err)
		os.Exit(0)
	}
}
