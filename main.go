package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

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
				runHook(r, line)
			}
		}
		if len(r.Match) > 0 {
			reg := regexp.MustCompile(r.Match)
			match := reg.FindString(line)
			if len(match) > 0 {
				runHook(r, match)
			}
		}
	}
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

func runHook(r Rule, logline string) {
	if !r.Skip {
		fmt.Printf("%s\n", logline)
	}
	cmd := fmt.Sprintf(r.Command, logline)

	process := exec.Command("/bin/bash", "-c", cmd)
	if err := process.Run(); err != nil {
		log.Fatal(err)
	}
}
