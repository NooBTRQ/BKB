package raft

import (
	infra "BlackKingBar/infrastructure"
	"context"
	"fmt"
	"io/fs"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/errors"
)

// raft状态机

const (
	Follower  = 0
	Candidate = 1
	Leader    = 2
)

type Raft struct {
	NodeId               int8
	State                int8
	CurrentTerm          int64
	LeaderId             int8
	VoteFor              int8
	CommitIndex          int64
	LastApplied          int64
	NextIndex            sync.Map
	MatchIndex           sync.Map
	Log                  []RaftLog
	ElectionTimeout      time.Duration
	ElectionTicker       *time.Ticker
	HeartBeatDuration    time.Duration
	HeartBeatTicker      *time.Ticker
	ReceiveHeartBeatTime time.Time
	lock                 sync.RWMutex
	con                  sync.Cond
	Data                 map[string]string
	EventCh              chan interface{}
	Done                 chan interface{}
}

var R *Raft

// 恢复持久化数据
// 将节点初始化为follower
// 启动ticker
func InitStateMachine(ctx context.Context) error {
	cfg := infra.CfgInstance
	R = &Raft{}
	R.NodeId = cfg.NodeId
	R.EventCh = make(chan interface{})
	R.Done = make(chan interface{})
	R.ElectionTicker = time.NewTicker(1000 * time.Second)
	R.HeartBeatTicker = time.NewTicker(1000 * time.Second)
	R.con = *sync.NewCond(&sync.Mutex{})
	err := R.RecoverMachine()
	if err != nil {
		return errors.Wrap(err, "recover machine err")
	}
	err = R.RecoverLog()
	if err != nil {
		return errors.Wrap(err, "recover log err")
	}
	err = R.RecoverData()
	if err != nil {
		return errors.Wrap(err, "recover data err")
	}
	R.BecomeFollower()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-R.ElectionTicker.C:
				d := time.Since(R.ReceiveHeartBeatTime)
				//fmt.Println(d)
				//fmt.Println(R.ElectionTimeout)
				if d > R.ElectionTimeout {
					R.EventCh <- Election{}
				}

			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-R.HeartBeatTicker.C:
				R.EventCh <- HeartBeat{}
			}
		}
	}()
	go R.ApplyLog(ctx)
	return nil
}

func (raft *Raft) LastLog() *RaftLog {

	l := len(raft.Log)
	return &raft.Log[l-1]
}

func (raft *Raft) Process(ctx context.Context) {
	var event interface{}
	for {
		select {
		case event = <-raft.EventCh:
			raft.handleEvent(event)
		case <-ctx.Done():
			return
		}
	}
}

func (raft *Raft) handleEvent(msg interface{}) {
	switch e := msg.(type) {
	case Election:
		//fmt.Println("处理发起选举1")
		raft.StartElection()
	case *VoteRequest:
		//fmt.Println("处理投票请求2")
		res := raft.HandleVoteRequest(e)
		e.Done <- res
	case HeartBeat:
		//fmt.Println("处理发起心跳3")
		raft.HeartBeat()
	case *Command:
		res := raft.HandleClientOpeartion(e)
		e.Done <- res
	case *AppendEntriesRequest:
		//fmt.Println("处理收到appendlog")
		res := raft.HandleAppendEntries(e)
		e.Done <- res
	default:
	}
}

func (raft *Raft) BecomeFollowerWithLogAppend(req *AppendEntriesRequest) {
	raft.State = Follower
	if req != nil {
		raft.LeaderId = req.LeaderId
		raft.CurrentTerm = req.Term
	}
	raft.ResetElectionTicker()
	raft.StopHeartBeatTicker()
}

func (raft *Raft) BecomeFollower() {
	raft.BecomeFollowerWithLogAppend(nil)
}

func (raft *Raft) BecomeCandicate() {
	raft.State = Candidate
	raft.CurrentTerm++
	raft.VoteFor = raft.NodeId
}

func (raft *Raft) BecomeLeader() {
	fmt.Println("成为leader")
	raft.State = Leader
	raft.LeaderId = raft.NodeId
	raft.HeartBeatDuration = time.Duration(10) * time.Millisecond
	raft.Log = append(raft.Log, RaftLog{Term: raft.CurrentTerm, Index: raft.LastLog().Index})
	cfg := infra.CfgInstance

	for _, member := range cfg.ClusterMembers {
		raft.NextIndex.Store(member.NodeId, raft.LastLog().Index+1)
	}
	raft.HeartBeat()
	raft.ResetHeartBeatTicker()
	raft.StopElectionTicker()
}

func (raft *Raft) ResetElectionTicker() {
	raft.ElectionTicker.Stop()
	timeout, _ := infra.GetRandNumber(150, 300, int(raft.NodeId))
	raft.ElectionTimeout = time.Duration(timeout) * time.Millisecond
	raft.ElectionTicker.Reset(raft.ElectionTimeout)
}

func (raft *Raft) StopElectionTicker() {
	raft.ElectionTicker.Stop()
}

func (raft *Raft) ResetHeartBeatTicker() {
	raft.HeartBeatDuration = time.Duration(10 * time.Millisecond)
	raft.HeartBeatTicker.Reset(raft.HeartBeatDuration)
}

func (raft *Raft) StopHeartBeatTicker() {
	raft.HeartBeatTicker.Stop()
}

func isFileNotFound(err error) bool {
	pathErr, ok := err.(*fs.PathError)
	return ok && pathErr.Err == syscall.ERROR_FILE_NOT_FOUND
}

func (raft *Raft) SaveMachine() error {
	data := make([]byte, 0)
	termByte, err := infra.IntToBytes(int64(raft.CurrentTerm))
	if err != nil {
		return err
	}
	voteByte, err := infra.IntToBytes(int64(raft.VoteFor))
	if err != nil {
		return err
	}
	lastApplyByte, err := infra.IntToBytes(int64(raft.LastApplied))
	if err != nil {
		return err
	}
	data = append(data, termByte...)
	data = append(data, voteByte...)
	data = append(data, lastApplyByte...)
	cfg := infra.CfgInstance
	err = infra.WriteFile(data, cfg.MachineFilePath)
	return err
}

func (raft *Raft) RecoverMachine() error {
	cfg := infra.CfgInstance
	data, err := infra.ReadFile(cfg.MachineFilePath)

	if isFileNotFound(err) {
		return nil
	}

	if err != nil {
		return err
	}

	err = infra.BytesToInt(data[0:8], &raft.CurrentTerm)
	if err != nil {
		return err
	}
	err = infra.BytesToInt(data[8:9], &raft.VoteFor)
	if err != nil {
		return err
	}

	err = infra.BytesToInt(data[9:], &raft.LastApplied)
	return err
}

func (raft *Raft) SaveLog(log *RaftLog) error {
	cfg := infra.CfgInstance
	data := make([]byte, 0)
	indexByte, err := infra.IntToBytes(log.Index)
	if err != nil {
		return err
	}
	termByte, err := infra.IntToBytes(log.Term)
	if err != nil {
		return err
	}
	data = append(data, indexByte...)
	data = append(data, termByte...)

	opByte, err := infra.IntToBytes(log.Cmd.Operation)
	if err != nil {
		return err
	}
	data = append(data, opByte...)
	kByte := []byte(log.Cmd.Key)
	vByte := []byte(log.Cmd.Value)
	klenByte, err := infra.IntToBytes(int16(len(kByte)))
	if err != nil {
		return err
	}
	vlenByte, err := infra.IntToBytes(int16(len(vByte)))
	if err != nil {
		return err
	}
	data = append(data, klenByte...)
	data = append(data, vlenByte...)
	data = append(data, kByte...)
	data = append(data, vByte...)
	err = infra.AppendFile(data, cfg.LogFilePath)
	return err
}

func (raft *Raft) RecoverLog() error {

	raft.Log = make([]RaftLog, 0)
	raft.Log = append(raft.Log, RaftLog{Cmd: Command{Operation: -1}})
	cfg := infra.CfgInstance
	data, err := infra.ReadFile(cfg.LogFilePath)
	if isFileNotFound(err) {
		return nil
	}

	if err != nil {
		return err
	}

	for len(data) > 0 {
		log := RaftLog{}
		err = infra.BytesToInt(data[:8], &log.Index)
		if err != nil {
			return err
		}
		data = data[8:]
		err = infra.BytesToInt(data[:8], &log.Term)
		if err != nil {
			return err
		}
		data = data[8:]
		log.Cmd = Command{}
		err = infra.BytesToInt(data[:1], &log.Cmd.Operation)
		if err != nil {
			return err
		}
		data = data[1:]
		var kLen int16
		var vLen int16
		err = infra.BytesToInt(data[:2], &kLen)
		if err != nil {
			return err
		}
		data = data[2:]
		err = infra.BytesToInt(data[:2], &vLen)
		if err != nil {
			return err
		}
		data = data[2:]
		log.Cmd.Key = string(data[:kLen])
		data = data[kLen:]
		log.Cmd.Value = string(data[:kLen])
		data = data[kLen:]
		raft.Log = append(raft.Log, log)
	}

	return nil
}

func (raft *Raft) SaveData() error {

	cfg := infra.CfgInstance
	data := make([]byte, 0)
	for k, v := range raft.Data {
		kByte := []byte(k)
		vByte := []byte(v)
		klenByte, err := infra.IntToBytes(int16(len(kByte)))
		if err != nil {
			return errors.Wrap(err, "save data err")
		}
		vlenByte, err := infra.IntToBytes(int16(len(vByte)))
		if err != nil {
			return errors.Wrap(err, "save data err")
		}
		data = append(data, klenByte...)
		data = append(data, vlenByte...)
		data = append(data, kByte...)
		data = append(data, vByte...)

		err = infra.WriteFile(data, cfg.DataFilePath)

		return err
	}

	return nil
}

func (raft *Raft) RecoverData() error {
	raft.Data = make(map[string]string)

	cfg := infra.CfgInstance
	data, err := infra.ReadFile(cfg.LogFilePath)
	if isFileNotFound(err) {
		return nil
	}

	if err != nil {
		return err
	}

	for len(data) > 0 {
		var kLen int16
		var vLen int16
		err = infra.BytesToInt(data[:2], &kLen)
		if err != nil {
			return err
		}
		data = data[2:]
		err = infra.BytesToInt(data[:2], &vLen)
		if err != nil {
			return err
		}
		data = data[2:]
		raft.Data[string(data[:kLen])] = string(data[kLen : kLen+vLen])
		data = data[:kLen+vLen]
	}

	return nil
}

func (raft *Raft) ApplyLog(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return
		default:
			for atomic.LoadInt64(&raft.CommitIndex) > raft.LastApplied {
				raft.lock.Lock()

				log := raft.Log[raft.LastApplied+1]
				if log.Cmd.Operation == Set {
					raft.Data[log.Cmd.Key] = log.Cmd.Value
				} else if log.Cmd.Operation == Delete {
					delete(raft.Data, log.Cmd.Key)
				}
				raft.lock.Unlock()
				raft.LastApplied++
				raft.SaveData()
				raft.con.Broadcast()
			}
		}
		time.Sleep(10 * time.Millisecond)
	}

}

func (raft *Raft) GetValue(key string) (string, bool) {

	for atomic.LoadInt64(&raft.LastApplied) < atomic.LoadInt64(&raft.CommitIndex) {
		raft.con.L.Lock()
		raft.con.Wait()
		raft.con.L.Unlock()
	}

	raft.lock.RLock()
	defer raft.lock.RUnlock()
	value, has := raft.Data[key]
	return value, has
}
