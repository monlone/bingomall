package datetime

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	/** 设置每周的起始时间 */
	WeekStartDay = time.Sunday

	/** 指定日期和时间的默认转换格式 */
	dateTimeFormats = []string{"1/2/2006", "1/2/2006 15:4:5", "2006", "2006-1-2", "2006-01-02 15:04:05", "20060102150405", "15:4:5 Jan 2, 2006 MST"}
)

const (
	DefaultFormat  = "2006-01-02 15:04:05"
	CompressFormat = "20060102150405"
)

// DateTime 结构体
type DateTime struct {
	time.Time
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf(`"%s"`, time.Now().Format("2006-01-02 15:04:05"))
	return []byte(str), nil
}

func (t DateTime) MarshalJSONSecond() ([]byte, error) {
	//格式化秒
	seconds := t.Unix()
	return []byte(strconv.FormatInt(seconds, 10)), nil
}

func (t DateTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *DateTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = DateTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// 当前秒
func CurrentSecond() time.Time {
	return time.Now().Truncate(time.Second)
}

// 当前分钟
func CurrentMinute() time.Time {
	return time.Now().Truncate(time.Minute)
}

// 当前小时
func CurrentHour() time.Time {
	return time.Now().Truncate(time.Hour)
}

// 返回今天的日期
func Today() time.Time {
	year, month, day := time.Now().Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())
}

// 返回今天的最后一刻
func TodayEndMoment() time.Time {
	year, month, day := time.Now().Date()
	return time.Date(year, month, day, 23, 59, 59, int(time.Second-time.Nanosecond), time.Now().Location())
}

// 返回本周的第一刻
func BeginThisWeek() time.Time {
	today := Today()
	weekday := int(today.Weekday())
	if WeekStartDay != time.Sunday {
		weekStartDayInt := int(WeekStartDay)
		if weekday < weekStartDayInt {
			weekday = weekday + 7 - weekStartDayInt
		} else {
			weekday = weekday - weekStartDayInt
		}
	}
	return today.AddDate(0, 0, -weekday)
}

// 返回本周最后一刻
func EndThisWeek() time.Time {
	return BeginThisWeek().AddDate(0, 0, 7).Add(-time.Nanosecond)
}

// 返回本月的第一刻
func BeginThisMonth() time.Time {
	year, month, _ := time.Now().Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, time.Now().Location())
}

// 返回本月最后一刻
func EndThisMonth() time.Time {
	return BeginThisMonth().AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// 返回本年的第一刻
func BeginThisYear() time.Time {
	year, _, _ := time.Now().Date()
	return time.Date(year, 1, 1, 0, 0, 0, 0, time.Now().Location())
}

// 返回本年的最后一刻
func EndThisYear() time.Time {
	return BeginThisYear().AddDate(1, 0, 0).Add(-time.Nanosecond)
}

// 字符串转时间
func Parse(str string) (t time.Time, err error) {
	for _, format := range dateTimeFormats {
		t, err = time.Parse(format, str)
		if err == nil {
			return t, err
		}
	}
	err = errors.New("Can't parse string as time: " + str)
	return t, err
}

// 获取当前时间字符串 - yyyy-MM-dd HH:mm:ss
func (t *DateTime) CurrentDefault() string {
	return time.Now().In(time.Local).Format(DefaultFormat)
}

func String(time DateTime) string {
	return time.Format(DefaultFormat)
}

func (t *DateTime) CurrentTime() time.Time {
	return t.In(time.Local)
	//return time.Now().In(time.Local)	//原来是这个样子的。
}

// 获取当前时间字符串 - yyyyMMddHHmmss
func CurrentCompress() string {
	return time.Now().In(time.Local).Format(CompressFormat)
}
