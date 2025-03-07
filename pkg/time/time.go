package use_time

import "time"

// FormatDateTime 格式化时间为标准日期时间字符串
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// ParseDateTime 解析标准日期时间字符串为时间对象
func ParseDateTime(timeStr string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", timeStr)
}

// FormatDate 格式化时间为标准日期字符串
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// ParseDate 解析标准日期字符串为时间对象
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// GetStartOfDay 获取指定时间的当天开始时间
func GetStartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// GetEndOfDay 获取指定时间的当天结束时间
func GetEndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

// GetStartOfWeek 获取指定时间所在周的开始时间（周一）
func GetStartOfWeek(t time.Time) time.Time {
	weekday := t.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	return GetStartOfDay(t.AddDate(0, 0, -int(weekday-1)))
}

// GetEndOfWeek 获取指定时间所在周的结束时间（周日）
func GetEndOfWeek(t time.Time) time.Time {
	weekday := t.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	return GetEndOfDay(t.AddDate(0, 0, 7-int(weekday)))
}

// GetStartOfMonth 获取指定时间所在月的开始时间
func GetStartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// GetEndOfMonth 获取指定时间所在月的结束时间
func GetEndOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	lastDay := time.Date(year, month+1, 0, 23, 59, 59, 999999999, t.Location())
	return lastDay
}
