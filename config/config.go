package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	HttpIP        string
	HttpPort      string
	RpcIP         string
	RpcPort       string
	CandidateId   int
	ClusterMember []Config
}

var CfgInstance *Config

func InitConfig() error {

	configStr, err := ioutil.ReadFile("config.json")

	if err != nil {
		return err
	}

	err = json.Unmarshal(configStr, CfgInstance)

	return err
}
