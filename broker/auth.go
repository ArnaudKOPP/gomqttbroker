package broker

import (
	"strings"

	"../lib/acl"
)

const (
	// PUB acl status
	PUB = 1
	// SUB acl status
	SUB = 2
)

func (c *client) CheckTopicAuth(typ int, topic string) bool {
	if !c.broker.config.ACL {
		return true
	}
	if strings.HasPrefix(topic, "$queue/") {
		topic = string([]byte(topic)[7:])
		if topic == "" {
			return false
		}
	}
	ip := c.info.remoteIP
	username := string(c.info.username)
	clientid := string(c.info.clientID)
	aclInfo := c.broker.ACLConfig
	return acl.CheckTopicAuth(aclInfo, typ, ip, username, clientid, topic)

}

var (
	watchList = []string{"./conf"}
)
