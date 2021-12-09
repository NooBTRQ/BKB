package app

import (
	infra "BlackKingBar/infrastructure"
	"fmt"
	"sync"
	"syscall"
	"time"
)

// raft状态机

const (
	Follower  = 0
	Candidate = 1
	Leader    = 2
)

type StateMachine struct {
	CandidateId       int8
	State             int8
	CurrentTerm       int64
	VoteFor           int8
	CommitIndex       int64
	LastApplied       int64
	NextIndex         map[int8]int64
	MatchIndex        map[int8]int64
	Log               []RaftLog
	ElectionTimeout   time.Duration
	HeartBeatDuration time.Duration
	ElectionTimer     *time.Ticker
	HeartBeatTimer    *time.Ticker
	lock              sync.RWMutex
	wg                sync.WaitGroup
	data              map[string]string
}

var Raft *StateMachine

// 恢复持久化数据
// 初始化timer
// 并将节点初始化为follower
func InitStateMachine() error {
	cfg := infra.CfgInstance
	Raft = &StateMachine{}
	Raft.CandidateId = cfg.CandidateId
	Raft.ElectionTimer = time.NewTicker(time.Second * 3600)
	Raft.HeartBeatTimer = time.NewTicker(time.Second * 3600)
	Raft.RecoverMachine()
	Raft.RecoverLog()
	Raft.RecoverData()
	Raft.BecomeFollower()
	return nil
}

func (raft *StateMachine) BecomeFollower() {
	fmt.Println(time.Now())
	fmt.Println("成为follower")
	raft.State = Follower

	raft.ElectionTimeout = time.Duration(infra.GetRandNumber(150, 300, int(raft.CandidateId))) * time.Millisecond
	raft.ElectionTimer.Reset(raft.ElectionTimeout)
	raft.stopHeartBeat()
	go func() {
		for {
			<-raft.ElectionTimer.C
			raft.StartElection()
		}
	}()
}

func (raft *StateMachine) BecomeCandicate() {
	fmt.Println(time.Now())
	fmt.Println("成为candidate")
	raft.State = Candidate
	raft.CurrentTerm++
	raft.VoteFor = raft.CandidateId
}

func (raft *StateMachine) BecomeLeader() {
	fmt.Println(time.Now())
	fmt.Println("成为leader")
	raft.State = Leader
	raft.ElectionTimer.Stop()
	raft.HeartBeatDuration = time.Duration(10) * time.Millisecond
	raft.HeartBeatTimer.Reset(raft.HeartBeatDuration)
	raft.MatchIndex = make(map[int8]int64)
	raft.NextIndex = make(map[int8]int64)
	cfg := infra.CfgInstance

	for _, member := range cfg.ClusterMembers {
		raft.NextIndex[member.CandidateId] = int64(len(raft.Log))
	}
	go func() {
		for {
			<-raft.HeartBeatTimer.C
			raft.SendHeartBeat()
		}
	}()
}

func (raft *StateMachine) stopHeartBeat() {
	raft.HeartBeatTimer.Stop()
}

func (raft *StateMachine) ResetElectionTimeout() {
	raft.ElectionTimer.Stop()
	raft.ElectionTimeout = time.Duration(infra.GetRandNumber(150, 300, int(raft.CandidateId))) * time.Millisecond
	raft.ElectionTimer.Reset(raft.ElectionTimeout)
}

func (raft *StateMachine) SaveMachine() {
	data := make([]byte, 0)
	data = append(data, infra.IntToBytes(int64(raft.CurrentTerm))...)
	data = append(data, infra.IntToBytes(int64(raft.VoteFor))...)
	cfg := infra.CfgInstance
	infra.WriteFile(data, cfg.MachineFilePath)
}

func (raft *StateMachine) RecoverMachine() error {
	cfg := infra.CfgInstance
	data, err := infra.ReadFile(cfg.MachineFilePath)

	if err == syscall.ERROR_FILE_NOT_FOUND {
		return nil
	}

	if err != nil {
		return err
	}

	infra.BytesToInt(data[0:8], &raft.CurrentTerm)
	infra.BytesToInt(data[8:], &raft.VoteFor)
	return nil
}

func (raft *StateMachine) SaveLog() {
	cfg := infra.CfgInstance
	data := make([]byte, 0)
	for _, log := range raft.Log {
		data = append(data, infra.IntToBytes(log.Index)...)
		data = append(data, infra.IntToBytes(log.Term)...)
		data = append(data, infra.IntToBytes(log.Cmd.Operation)...)
		kByte := []byte(log.Cmd.Key)
		vByte := []byte(log.Cmd.Value)

		data = append(data, infra.IntToBytes(int16(len(kByte)))...)
		data = append(data, infra.IntToBytes(int16(len(vByte)))...)
		data = append(data, kByte...)
		data = append(data, vByte...)

		infra.WriteFile(data, cfg.LogFilePath)
	}
}

func (raft *StateMachine) RecoverLog() error {

	raft.Log = make([]RaftLog, 0)

	cfg := infra.CfgInstance
	data, err := infra.ReadFile(cfg.LogFilePath)
	if err == syscall.ERROR_FILE_NOT_FOUND {
		return nil
	}

	if err != nil {
		return err
	}

	for len(data) > 0 {
		log := RaftLog{}
		infra.BytesToInt(data[:8], &log.Index)
		data = data[8:]
		infra.BytesToInt(data[:8], &log.Term)
		data = data[8:]
		log.Cmd = Command{}
		infra.BytesToInt(data[:1], &log.Cmd.Operation)
		data = data[1:]
		var kLen int16
		var vLen int16
		infra.BytesToInt(data[:2], kLen)
		data = data[2:]
		infra.BytesToInt(data[:2], vLen)
		data = data[2:]
		log.Cmd.Key = string(data[:kLen])
		data = data[kLen:]
		log.Cmd.Value = string(data[:kLen])
		data = data[kLen:]
		raft.Log = append(raft.Log, log)
	}

	return nil
}

func (raft *StateMachine) SaveData() {

	cfg := infra.CfgInstance
	data := make([]byte, 0)
	for k, v := range raft.data {
		kByte := []byte(k)
		vByte := []byte(v)
		data = append(data, infra.IntToBytes(int16(len(kByte)))...)
		data = append(data, infra.IntToBytes(int16(len(vByte)))...)
		data = append(data, kByte...)
		data = append(data, vByte...)

		infra.WriteFile(data, cfg.DataFilePath)
	}
}

func (raft *StateMachine) RecoverData() error {
	raft.data = make(map[string]string)

	cfg := infra.CfgInstance
	data, err := infra.ReadFile(cfg.LogFilePath)
	if err == syscall.ERROR_FILE_NOT_FOUND {
		return nil
	}

	if err != nil {
		return err
	}

	for len(data) > 0 {
		var kLen int16
		var vLen int16
		infra.BytesToInt(data[:2], kLen)
		data = data[2:]
		infra.BytesToInt(data[:2], vLen)
		data = data[2:]
		raft.data[string(data[:kLen])] = string(data[kLen : kLen+vLen])
		data = data[:kLen+vLen]
	}

	return nil
}
