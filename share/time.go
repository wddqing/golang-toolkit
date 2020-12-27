package share

import (
	"fmt"
	"time"
)

const (
	MINUTE_HAS_SECS int64 = 60
	HOUR_HAS_SECS   int64 = MINUTE_HAS_SECS * 60
)

//获取 当前年月份 +上个月份
func GetMonths() (string, string) {
	t := time.Now()
	year := t.Year()
	month := t.Month()
	c := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	cc := c.Format("2006_01")
	if month == 1 {
		year = year - 1
		month = 12
		l := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
		ll := l.Format("2006_01")
		return cc, ll
	}
	lastMonth := month - 1
	l := time.Date(year, lastMonth, 1, 0, 0, 0, 0, time.Local)
	ll := l.Format("2006_01")
	return cc, ll
}

//GetForeverTime GetForeverTime
func GetForeverTime() int64 {
	//统一100年时间为2017-01-02开始
	t, _ := time.Parse("2006-01-02", "2017-01-02")
	timestamp := t.Unix() + 100*365*86400
	//mysql int(11) 最大值为4294967295
	if timestamp > 4294967295 {
		timestamp = 4294967295
	}
	return timestamp
}

//GetForeverSeconds 100年
func GetForeverSeconds() int64 {
	return 3153600000 //100 * 365 * 86400
}

func GetTenYearsTimestamp() int64 {
	timestamp := time.Now().Unix() + 10*365*86400
	//mysql int(11) 最大值为4294967295
	if timestamp > 4294967295 {
		timestamp = 4294967295
	}
	return timestamp
}

func FormatSecs(secs int64) string {
	var (
		fmtStr string
		hour   int64
		min    int64
	)
	if hour = secs / HOUR_HAS_SECS; hour > 0 {
		fmtStr = fmt.Sprintf("%d小时", hour)
		secs -= hour * HOUR_HAS_SECS
	}
	if min = secs / MINUTE_HAS_SECS; min > 0 {
		fmtStr = fmt.Sprintf(fmtStr+"%d分", min)
		secs -= min * MINUTE_HAS_SECS
	} else if hour > 0 {
		fmtStr = fmt.Sprintf(fmtStr+"%d分", 0)
	}
	return fmt.Sprintf(fmtStr+"%d秒", secs)
}

func FormatTimeInChina(t int64) string {
	tf := time.Unix(t, 0)
	nt := tf.Format("2006年01月02日 15:04:05")
	return nt
}

// 获取当天剩余时间
func GetTodayRemainTime() int64 {
	return GetTomorrowZeroTime() - time.Now().Unix()
}

// 获取当天零点
func GetTodayZeroTime() int64 {
	t := time.Now().Format("2006-01-02")
	tm, _ := time.ParseInLocation("2006-01-02", t, time.Local)
	return tm.Unix()
}

// 获取第二天零点时间
func GetTomorrowZeroTime() int64 {
	return GetTodayZeroTime() + 86400
}

// 获取零点时间
func GetZeroTime(t int64) int64 {
	tt := time.Unix(t, 0).Format("2006-01-02")
	tm, _ := time.ParseInLocation("2006-01-02", tt, time.Local)
	return tm.Unix()
}

// 获取当天时间与 t 相差的自然日天数
func GetNaturalDays(t int64) int32 {
	todayZero := GetTodayZeroTime()
	tZero := GetZeroTime(t)
	if todayZero < tZero {
		return 0
	}

	day := (todayZero-tZero)/86400 + 1
	return int32(day)
}

// 启动的时候执行一次，然后等待零点再重复执行
func CacheZeroLoop(f func()) {
	go func() {
		for {
			f()
			now := time.Now()
			// 计算下一个零点
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
		}
	}()
}

// 调用是执行一次，每隔sec 秒重复执行
func CacheLoop(f func(), sec int) {

	go func(sec int) {
		ticker := time.NewTicker(time.Duration(sec) * time.Second)
		f()
		for range ticker.C {
			f()
		}
	}(sec)
}
