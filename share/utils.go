package share

import (
	"encoding/binary"
	"net"
	"runtime/debug"
	"sync"
	"time"

	"go.uber.org/zap"
)

//获取某个时刻所属星期的星期一
func GetWeekMonday(t time.Time) time.Time {
	weekday := time.Duration(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	year, month, day := t.Date()
	currentZeroDay := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return currentZeroDay.Add(-1 * (weekday - 1) * 24 * time.Hour)
}

func GetTodayBeginAt() time.Time {
	year, month, day := time.Now().Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func IP2Uint(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func Uint2IP(p uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, p)
	return ip
}

// 在独立的协程安全地执行给定的函数，执行过程中会捕获panic
func SafeExecFunc(fn func(...interface{}), waitGroup *sync.WaitGroup, args ...interface{}) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				zap.L().Error("从异常中恢复", zap.Any("e", err), zap.String("stack", string(debug.Stack())))
			}
			if waitGroup != nil {
				waitGroup.Done()
			}
		}()

		fn(args...)
	}()
}

func ContainsInt32(slice []int32, i int32) bool {
	for _, v := range slice {
		if v == i {
			return true
		}
	}
	return false
}
