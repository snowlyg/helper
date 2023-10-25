package global

import (
	"database/sql/driver"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

// DateTime 自定义事件类型
type DateTime time.Time

//  日期格式
const (
	DateLayout      = "2006-01-02"
	DateTimeLayout  = "2006-01-02 15:04:05"
	TimeLayout      = "15:04:05"
	BuildTimeLayout = "2006.0102.150405"
	TimestampLayout = "20060102150405"
)

var StartTime = time.Now()

func (dt *DateTime) UnmarshalJSON(data []byte) (err error) {
	value := strings.Trim(string(data), "\"")
	now, err := time.ParseInLocation(DateTimeLayout, value, time.Local)
	*dt = DateTime(now)
	return
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(DateTimeLayout)+2)
	b = append(b, '"')
	b = time.Time(dt).AppendFormat(b, DateTimeLayout)
	b = append(b, '"')
	return b, nil
}

func (dt DateTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	ti := time.Time(dt)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

func (dt *DateTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*dt = DateTime(value)
		return nil
	}
	return nil
}

func (dt DateTime) String() string {
	return time.Time(dt).Format(DateTimeLayout)
}

func UpTime() time.Duration {
	return time.Since(StartTime)
}

func UpTimeString() string {
	d := UpTime()
	days := d / (time.Hour * 24)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	d -= minutes * time.Minute
	seconds := d / time.Second
	return fmt.Sprintf("%d Days %d Hours %d Mins %d Secs", days, hours, minutes, seconds)
}

// 获取时区
func GetLocation() (*time.Location, error) {
	location, err := time.LoadLocation("Local")
	if err != nil {
		return location, nil
	}
	return location, nil
}

func GetDuration(t int64, v string) time.Duration {
	switch v {
	case "h":
		return time.Hour * time.Duration(t)
	case "m":
		return time.Minute * time.Duration(t)
	case "s":
		return time.Second * time.Duration(t)
	default:
		return time.Minute * time.Duration(t)
	}
}

func SleepRandomDuration() {
	ns := int64(5) * 1000000000
	// 以当前时间为随机数种子，如果所有 log-agent-updater 在同一时间启动，系统时间是相同的，那么随机种子就是一样的
	// 问题不大，批量ssh去启动 log-agent-updater 的话也是一个顺次的过程
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	d := time.Duration(r.Int63n(ns)) * time.Nanosecond
	time.Sleep(d)
}

// OverTimeNow 超时比较
func OverTimeNow(compareTime, subTime time.Time, sub int64) (bool, int64, error) {
	location, err := GetLocation()
	if err != nil {
		return false, 0, err
	}
	subT := compareTime.In(location).Sub(subTime)
	abs := int64(math.Abs(math.Ceil(subT.Minutes())))
	if abs > sub {
		return true, abs, nil
	}
	return false, abs, nil
}
