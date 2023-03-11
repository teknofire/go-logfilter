package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Name    string `yaml:"name"`
	Match   string `yaml:"match"`
	Command string `yaml:"command"`
	Skip    bool   `yaml:"skip"`
}

func main() {
	configFile := os.Args[1]
	config := readConfig(configFile)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		logline := scanner.Text()
		ok := handleLine(logline, config.Rules)
		if ok {
			fmt.Printf("%s\n", logline)
		}
	}
}

func handleLine(line string, rules []Rule) bool {
	// Skip empty lines
	if len(line) == 0 {
		return false
	}

	for _, r := range rules {
		if strings.Contains(line, r.Match) {
			if r.Skip {
				return false
			}
			runHook(r.Command, line)
		}
	}

	return true
}

func readConfig(config string) Config {
	filename, _ := filepath.Abs(config)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	c := Config{}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		panic(err)
	}

	return c
}

func runHook(cmdFormat string, logline string) {
	cmd := fmt.Sprintf(cmdFormat, logline)
	process := exec.Command("/bin/bash", "-c", cmd)
	if err := process.Run(); err != nil {
		log.Fatal(err)
	}
}
