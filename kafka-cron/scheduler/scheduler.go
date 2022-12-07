package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

type cronjob struct {
	crontab string
	command string
	str     string
	// cluster string
	// retries int
}

func main() {
	log.Info("Create new cron")
	c := cron.New()
	cronjobs := readCrontabfile("crontab.txt")
	for _, job := range cronjobs {
		myJob := job
		_, er := c.AddFunc(job.crontab, func() {
			queueJob(myJob.command, []string{myJob.str})
		})
		if er != nil {
			fmt.Println(er)
		}
		fmt.Printf("cronjobs: started cron for %+v\n", myJob)
	}
	c.Run()
}

func queueJob(command string, args []string) {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Command Successfully Executed")
	fmt.Println(string(stdout))
}

func readCrontabfile(path string) []cronjob {
	readFile, err := os.Open("cronfile.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer readFile.Close()
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	result := make([]cronjob, 0)
	var fileLines []string
	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}
	for _, line := range fileLines {
		val := strings.Split(line, ",")
		cj := cronjob{
			crontab: val[0],
			command: val[1],
			str:     val[2],
		}
		result = append(result, cj)
	}

	return result
}
