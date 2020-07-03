package coming

import (
	"time"
)

type schedule struct {
	Month, Day, Weekday  int8
	Hour, Minute, Second int8
	Task                 func(time.Time)
}

const ANY = -1 // mod by MDR

var tasks []schedule

// 这个function新建一个简单定时任务到后台任务集,schedule可以通过month, day, weekday, hour, minute, second.
// 这6个字段匹配到一年中的任意时刻,最小精度为秒
// this function create a new scheduled task and append to tasks,whitch will be checked every second.
// this schedule task can match any time in a year with minimum precision second
func NewSchedule(month, day, weekday, hour, minute, second int8, task func(time.Time)) {
	cj := schedule{month, day, weekday, hour, minute, second, task}
	tasks = append(tasks, cj)
}

// 创建一个定时任务,匹配一个月中的任意时间,但时间必须是明确的
// this creates a new scheduled task matched any exactly time in a month.
func NewMonthlySchedule(day, hour, minute, second int8, task func(time.Time)) {
	NewSchedule(ANY, day, ANY, hour, minute, second, task)
}

// 创建一个定时任务,匹配一周中的任意时间,但时间必须是明确的
// this creates a new scheduled task matched any exactly time in a week.
func NewWeeklySchedule(weekday, hour, minute, second int8, task func(time.Time)) {
	NewSchedule(ANY, ANY, weekday, hour, minute, second, task)
}

// 创建一个定时任务,匹配一天中的任意时间,但时间必须是明确的
// this creates a new scheduled task matched any exactly time of day.
func NewDailySchedule(hour, minute, second int8, task func(time.Time)) {
	NewSchedule(ANY, ANY, ANY, hour, minute, second, task)
}

func (sc schedule) Matches(t time.Time) (ok bool) {
	ok = (sc.Month == ANY || sc.Month == int8(t.Month())) &&
		(sc.Day == ANY || sc.Day == int8(t.Day())) &&
		(sc.Weekday == ANY || sc.Weekday == int8(t.Weekday())) &&
		(sc.Hour == ANY || sc.Hour == int8(t.Hour())) &&
		(sc.Minute == ANY || sc.Minute == int8(t.Minute())) &&
		(sc.Second == ANY || sc.Second == int8(t.Second()))
	return ok
}

func processSchedules() {
	for {
		now := time.Now()
		for _, j := range tasks {
			// execute all our cron tasks asynchronously
			if j.Matches(now) {
				go j.Task(now)
			}
		}
		time.Sleep(time.Second)
	}
}

func init() {
	go processSchedules()
}
