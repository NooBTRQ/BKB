package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	HttpIP         string
	HttpPort       string
	RpcIP          string
	RpcPort        string
	CandidateId    int
	ClusterMembers []ClusterMember
}

type ClusterMember struct {
	RpcIP   string
	RpcPort string
}

var CfgInstance *Config

func InitConfig() error {
	configStr, err := ioutil.ReadFile("config/config")
	if err != nil {
		return err
	}
	CfgInstance = &Config{}
	err = json.Unmarshal(configStr, CfgInstance)

	return err
}
