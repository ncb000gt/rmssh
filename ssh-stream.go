package main

import (
	"bufio"
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type config struct {
	Host     string
	Port     int64
	Command  string
	User     string
	Password string
}

func NewConfig() config {
	return config{Port: 22}
}

func main() {
	filterPtr := flag.String("filter", "", "Regexp pattern filter")
	flag.Parse()

	fmt.Print("Username: ")
	var user string
	fmt.Scanf("%s\n", &user)

	fmt.Print("Password: ")
	password, _ := terminal.ReadPassword(0)
	fmt.Println()

	argLen := len(flag.Args())
	var wg sync.WaitGroup
	wg.Add(argLen)

	for i := 0; i < argLen; i++ {
		c := parseConfig(flag.Arg(i))
		c.User = user
		c.Password = string(password[:len(password)])

		err := c.Stream(*filterPtr)
		if err != nil {
			panic("Problem creating stream: " + err.Error())
		}

	}

	wg.Wait()
	fmt.Println("Done.")
}

func parseConfig(configLine string) config {
	scmd := strings.Split(configLine, " ")

	c := NewConfig()
	c.Host = scmd[0]
	c.Command = strings.Join(scmd[1:], " ")

	return c
}

func (ssh_config *config) connect() (*ssh.Session, error) {
	sshClientConfig := &ssh.ClientConfig{
		User:            ssh_config.User,
		Auth:            []ssh.AuthMethod{ssh.Password(ssh_config.Password)}, // just passwords for now...
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", ssh_config.Host+":"+strconv.FormatInt(ssh_config.Port, 10), sshClientConfig)

	if err != nil {
		panic("Failed to dial: " + ssh_config.Host + " - " + err.Error())
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		panic("Failed to create session: " + err.Error())
	}

	return session, nil
}

// much inspiration from https://github.com/hypersleep/easyssh/blob/master/easyssh.go as I better learn Go
func (ssh_config *config) Stream(filter string) (err error) {
	session, err := ssh_config.connect()

	if err != nil {
		return err
	}

	// connect to both outputs (they are of type io.Reader)
	outReader, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	errReader, err := session.StderrPipe()
	if err != nil {
		return err
	}

	err = session.Start(ssh_config.Command)
	stdoutScanner := bufio.NewScanner(outReader)
	stderrScanner := bufio.NewScanner(errReader)

	// continuously send the command's output over the channel
	go func(stdoutScanner, stderrScanner *bufio.Scanner) {
		for stdoutScanner.Scan() {
			outline := stdoutScanner.Text()
			if outline != "" {
				if filter != "" {
					match, _ := regexp.MatchString(filter, outline)
					if match {
						fmt.Println(color.GreenString(ssh_config.User+"@"+ssh_config.Host), " - ", outline)
					}
				} else {
					fmt.Println(color.GreenString(ssh_config.User+"@"+ssh_config.Host), " - ", outline)
				}
			}
		}

		for stderrScanner.Scan() {
			errline := stderrScanner.Text()
			if errline != "" {
				fmt.Println(color.GreenString(ssh_config.User+"@"+ssh_config.Host), " - ", errline)
			}
		}
	}(stdoutScanner, stderrScanner)

	return err
}
