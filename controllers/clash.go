package controllers

import (
	"bufio"
	"fmt"
	"github.com/morningfish/beluga/tools"
	"io"
	"os"
)

var (
	RuleFile string
	Rules    []string
)

func InitRule() error {
	rules, err := GetRuleFromFile()
	if err != nil {
		return err
	}
	if rules == nil {
		for _, host := range BindHost {
			Rules = append(Rules, fmt.Sprintf("DOMAIN-SUFFIX,%s,Proxy", host))
		}
		return nil
	}
	Rules = rules
	return nil
}
func GetRuleFromFile() ([]string, error) {
	var newRule []string
	if tools.Exists(RuleFile) {
		file, err := os.Open(RuleFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		br := bufio.NewReader(file)
		for {
			line, _, err := br.ReadLine()
			if err == io.EOF {
				break
			}
			newRule = append(newRule, string(line[:]))
		}
		return newRule, nil
	}
	return nil, nil // 不存在时，返回nil，不保存，走默认 rule
}
