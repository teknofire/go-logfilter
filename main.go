package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Name     string `yaml:"name"`
	Contains string `yaml:"contains"`
	Match    string `yaml:"match"`
	Command  string `yaml:"command"`
	Skip     bool   `yaml:"skip"`
}

func main() {
	configFile := os.Args[1]
	config := readConfig(configFile)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		logline := scanner.Text()
		handleLine(logline, config.Rules)
	}
}

func handleLine(line string, rules []Rule) {
	// Skip empty lines
	if len(line) == 0 {
		return
	}

	for _, r := range rules {
		if len(r.Contains) > 0 {
			if strings.Contains(line, r.Contains) {
				logrus.Debug("contains rule: ", r.Contains)
				runHook(r, line)
				return 
			}
		}
		if len(r.Match) > 0 {
			reg := regexp.MustCompile(r.Match)
			match := reg.FindString(line)
			if len(match) > 0 {
				logrus.Debug("match rule: ", r.Match)
				runHook(r, match)
				return
			}
		}
	}
	logrus.Info(line)
	
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

func runHook(r Rule, logline string) error {
	if r.Skip || len(r.Command) == 0 {
		return nil
	}
	
	var cmd string
	if strings.Count(r.Command, "%s") > 0 {
		cmd = fmt.Sprintf(r.Command, logline)
	} else {
		cmd = r.Command
	}

	process := exec.Command("/bin/bash", "-c", cmd)
	if err := process.Run(); err != nil {
		logrus.WithError(err).Error("Error processing command: ", cmd)
		return err
	}

	return nil
}
