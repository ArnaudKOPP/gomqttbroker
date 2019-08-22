package acl

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	// PUB publish fct
	PUB = 1
	// SUB subscribe fct
	SUB = 2
	// PUBSUB publish & sub fct
	PUBSUB = 3
	// CLIENTID client id name
	CLIENTID = "clientid"
	// USERNAME username name
	USERNAME = "username"
	// IP ip name
	IP = "ip"
	// ALLOW name
	ALLOW = "allow"
	// DENY name
	DENY = "deny"
)

// AuthInfo struct for auth info
type AuthInfo struct {
	Auth   string
	Typ    string
	Val    string
	PubSub int
	Topics []string
}

// Config struct for ACL config
type Config struct {
	File string
	Info []*AuthInfo
}

// ConfigLoad load ACL config
func ConfigLoad(file string) (*Config, error) {
	if file == "" {
		file = "./conf/acl.conf"
	}
	aclconifg := &Config{
		File: file,
		Info: make([]*AuthInfo, 0, 4),
	}
	err := aclconifg.Parse()
	if err != nil {
		return nil, err
	}
	return aclconifg, err
}

// Parse parse ACL config file
func (c *Config) Parse() error {
	f, err := os.Open(c.File)
	defer f.Close()
	if err != nil {
		return err
	}
	buf := bufio.NewReader(f)
	var parseErr error
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if isCommentOut(line) {
			continue
		}
		if line == "" {
			return parseErr
		}
		// fmt.Println(line)
		tmpArr := strings.Fields(line)
		if len(tmpArr) != 5 {
			parseErr = errors.New("\"" + line + "\" format is error")
			break
		}
		if tmpArr[0] != ALLOW && tmpArr[0] != DENY {
			parseErr = errors.New("\"" + line + "\" format is error")
			break
		}
		if tmpArr[1] != CLIENTID && tmpArr[1] != USERNAME && tmpArr[1] != IP {
			parseErr = errors.New("\"" + line + "\" format is error")
			break
		}
		var pubsub int
		pubsub, err = strconv.Atoi(tmpArr[3])
		if err != nil {
			parseErr = errors.New("\"" + line + "\" format is error")
			break
		}
		topicStr := strings.Replace(tmpArr[4], " ", "", -1)
		topicStr = strings.Replace(topicStr, "\n", "", -1)
		topics := strings.Split(topicStr, ",")
		tmpAuth := &AuthInfo{
			Auth:   tmpArr[0],
			Typ:    tmpArr[1],
			Val:    tmpArr[2],
			Topics: topics,
			PubSub: pubsub,
		}
		c.Info = append(c.Info, tmpAuth)
		if err != nil {
			if err != io.EOF {
				parseErr = err
			}
			break
		}
	}
	return parseErr
}
func isCommentOut(line string) bool {
	if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "*") {
		return true
	}
	return false
}
