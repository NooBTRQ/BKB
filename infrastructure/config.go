package infrastructure

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	HttpIP            string
	HttpPort          string
	RpcIP             string
	RpcPort           string
	CandidateId       int8
	ElectionTimeout   int16
	HeartBeatDuration int16
	ClusterMembers    []ClusterMember
}

type ClusterMember struct {
	CandidateId int8
	RpcIP       string
	RpcPort     string
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
