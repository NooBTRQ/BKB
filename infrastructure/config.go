package infrastructure

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	MachineFilePath   string
	DataFilePath      string
	LogFilePath       string
	LogInfoPath       string
	HttpIP            string
	HttpPort          string
	RpcIP             string
	RpcPort           string
	NodeId            int8
	ElectionTimeout   int16
	HeartBeatDuration int16
	Peers             []ClusterMember
}

type ClusterMember struct {
	NodeId  int8
	RpcIP   string
	RpcPort string
}

var CfgInstance *Config

func InitConfig() error {
	configStr, err := ioutil.ReadFile("./config")
	if err != nil {
		return err
	}
	CfgInstance = &Config{}
	err = json.Unmarshal(configStr, CfgInstance)

	return err
}

func (cfg *Config) HalfCount() int64 {

	return int64(len(cfg.Peers) / 2)
}
