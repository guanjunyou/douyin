package utils

import (
	"sync"
	"time"
)

// Snowflake 结构体
type Snowflake struct {
	mu         sync.Mutex
	timestamp  int64 // 时间戳部分
	machineID  int64 // 机器 ID 部分
	sequenceID int64 // 序列号部分
}

// NewSnowflake 函数，返回一个新的 Snowflake 实例
func NewSnowflake(machineID int64) *Snowflake {
	return &Snowflake{
		timestamp:  0,
		machineID:  machineID,
		sequenceID: 0,
	}
}

// NextID 方法，生成下一个唯一的 ID
func (sf *Snowflake) NextID() int64 {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	// 获取当前时间戳，单位为毫秒
	now := time.Now().UnixNano() / int64(time.Millisecond)

	// 如果当前时间戳与上次生成的时间戳相同，则序列号递增
	if sf.timestamp == now {
		sf.sequenceID++
	} else {
		// 否则，重置序列号为 0
		sf.sequenceID = 0
	}

	// 更新时间戳为当前时间戳
	sf.timestamp = now

	// 生成 ID，包括时间戳、机器 ID 和序列号部分
	ID := (now << 22) | (sf.machineID << 10) | sf.sequenceID
	return ID
}
