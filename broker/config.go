package broker

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"../logger"
	"go.uber.org/zap"
)

// Config hold config parameters json fmt
type Config struct {
	Worker  int     `json:"workerNum"`
	Host    string  `json:"host"`
	Port    string  `json:"port"`
	Router  string  `json:"router"`
	TLSHost string  `json:"tlsHost"`
	TLSPort string  `json:"tlsPort"`
	WsPath  string  `json:"wsPath"`
	WsPort  string  `json:"wsPort"`
	WsTLS   bool    `json:"wsTLS"`
	TLSInfo TLSInfo `json:"tlsInfo"`
	ACL     bool    `json:"acl"`
	ACLConf string  `json:"aclConf"`
	Debug   bool    `json:"-"`
}

// TLSInfo struct hold TLS info
type TLSInfo struct {
	Verify   bool   `json:"verify"`
	CaFile   string `json:"caFile"`
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

// DefaultConfig create default config
var DefaultConfig = &Config{
	Worker: 4096,
	Host:   "0.0.0.0",
	Port:   "1883",
	ACL:    false,
}

var (
	log *zap.Logger
)

func showHelp() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

// ConfigureConfig create config
func ConfigureConfig(args []string) (*Config, error) {
	config := &Config{}
	var (
		help       bool
		configFile string
	)
	fs := flag.NewFlagSet("Mqt-broker", flag.ExitOnError)
	fs.Usage = showHelp

	fs.BoolVar(&help, "h", false, "Show this message.")
	fs.BoolVar(&help, "help", false, "Show this message.")
	fs.IntVar(&config.Worker, "w", 1024, "worker num to process message, perfer (client num)/10.")
	fs.IntVar(&config.Worker, "worker", 1024, "worker num to process message, perfer (client num)/10.")
	fs.StringVar(&config.Port, "port", "1883", "Port to listen on.")
	fs.StringVar(&config.Port, "p", "1883", "Port to listen on.")
	fs.StringVar(&config.Host, "host", "0.0.0.0", "Network host to listen on")
	fs.StringVar(&config.Router, "r", "", "Router who maintenance cluster info")
	fs.StringVar(&config.Router, "router", "", "Router who maintenance cluster info")
	fs.StringVar(&config.WsPort, "ws", "", "port for ws to listen on")
	fs.StringVar(&config.WsPort, "wsport", "", "port for ws to listen on")
	fs.StringVar(&config.WsPath, "wsp", "", "path for ws to listen on")
	fs.StringVar(&config.WsPath, "wspath", "", "path for ws to listen on")
	fs.StringVar(&configFile, "config", "", "config file for broker")
	fs.StringVar(&configFile, "c", "", "config file for broker")
	fs.BoolVar(&config.Debug, "debug", false, "enable Debug logging.")
	fs.BoolVar(&config.Debug, "d", false, "enable Debug logging.")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if help {
		showHelp()
		return nil, nil
	}

	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "D":
			config.Debug = true
		}
	})

	logger.InitLogger(config.Debug)
	log = logger.Get().Named("Broker")

	if configFile != "" {
		tmpConfig, e := LoadConfig(configFile)
		if e != nil {
			return nil, e
		}
		config = tmpConfig
	}

	if err := config.check(); err != nil {
		return nil, err
	}

	return config, nil

}

// LoadConfig load from string(filename) and return a config
func LoadConfig(filename string) (*Config, error) {

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error("Read config file error: ", zap.Error(err))
		return nil, err
	}
	// log.Info(string(content))

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Error("Unmarshal config file error: ", zap.Error(err))
		return nil, err
	}

	return &config, nil
}

func (config *Config) check() error {

	if config.Worker == 0 {
		config.Worker = 1024
	}

	if config.Port != "" {
		if config.Host == "" {
			config.Host = "0.0.0.0"
		}
	}

	if config.TLSPort != "" {
		if config.TLSInfo.CertFile == "" || config.TLSInfo.KeyFile == "" {
			log.Error("tls config error, no cert or key file.")
			return errors.New("tls config error, no cert or key file")
		}
		if config.TLSHost == "" {
			config.TLSHost = "0.0.0.0"
		}
	}
	return nil
}

// NewTLSConfig create TLS config from TLS info
func NewTLSConfig(tlsInfo TLSInfo) (*tls.Config, error) {

	cert, err := tls.LoadX509KeyPair(tlsInfo.CertFile, tlsInfo.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("error parsing X509 certificate/key pair: %v", zap.Error(err))
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %v", zap.Error(err))
	}

	// Create TLSConfig
	// We will determine the cipher suites that we prefer.
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	// Require client certificates as needed
	if tlsInfo.Verify {
		config.ClientAuth = tls.RequireAndVerifyClientCert
	}
	// Add in CAs if applicable.
	if tlsInfo.CaFile != "" {
		rootPEM, err := ioutil.ReadFile(tlsInfo.CaFile)
		if err != nil || rootPEM == nil {
			return nil, err
		}
		pool := x509.NewCertPool()
		ok := pool.AppendCertsFromPEM([]byte(rootPEM))
		if !ok {
			return nil, fmt.Errorf("failed to parse root ca certificate")
		}
		config.ClientCAs = pool
	}

	return &config, nil
}
