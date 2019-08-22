package broker

import (
	"strings"

	"../lib/acl"
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
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

func (b *Broker) handleFsEvent(event fsnotify.Event) error {
	switch event.Name {
	case b.config.ACLConf:
		if event.Op&fsnotify.Write == fsnotify.Write ||
			event.Op&fsnotify.Create == fsnotify.Create {
			log.Info("text:handling acl config change event:", zap.String("filename", event.Name))
			aclconfig, err := acl.ConfigLoad(event.Name)
			if err != nil {
				log.Error("aclconfig change failed, load acl conf error: ", zap.Error(err))
				return err
			}
			b.ACLConfig = aclconfig
		}
	}
	return nil
}

// StartACLWatcher start service if file change
func (b *Broker) StartACLWatcher() {
	go func() {
		log.Debug("Start ACL monitoring")
		wch, e := fsnotify.NewWatcher()
		if e != nil {
			log.Error("Start monitor acl config file error,", zap.Error(e))
			return
		}
		defer wch.Close()

		for _, i := range watchList {
			if err := wch.Add(i); err != nil {
				log.Error("Start monitor acl config file error,", zap.Error(err))
				return
			}
		}
		log.Info("Watching acl config file change...")
		for {
			select {
			case evt := <-wch.Events:
				b.handleFsEvent(evt)
			case err := <-wch.Errors:
				log.Error("Error:", zap.Error(err))
			}
		}
	}()
}
