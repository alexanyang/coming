package coming

import (
	"sync"
	"time"
)

type schedule struct {
	Month, Day, Weekday  int8
	Hour, Minute, Second int8
	Task                 func(time.Time)
}

const ANY = -1 // mod by MDR

var (
	tasks  map[string]schedule
	locker sync.RWMutex
)

// NewSchedule create a new scheduled task and append to tasks,whitch will be checked every second.
// this schedule task can match any time in a year with minimum precision second
// 这个function新建一个简单定时任务到后台任务集,schedule可以通过month, day, weekday, hour, minute, second.
// 这6个字段匹配到一年中的任意时刻,最小精度为秒
func NewSchedule(month, day, weekday, hour, minute, second int8, name string, task func(time.Time)) {
	cj := schedule{month, day, weekday, hour, minute, second, task}
	locker.Lock()
	defer locker.Unlock()
	tasks[name] = cj
}

// RemoveSchedule can remove schedule task safely
func RemoveSchedule(name string) {
	locker.Lock()
	defer locker.Unlock()
	delete(tasks, name)
}

// NewMonthlySchedule creates a new schedule task matched any exactly time in a month.
// 创建一个定时任务,匹配一个月中的任意时间,但时间必须是明确的
func NewMonthlySchedule(day, hour, minute, second int8, name string, task func(time.Time)) {
	NewSchedule(ANY, day, ANY, hour, minute, second, name, task)
}

// NewWeeklySchedule creates a new scheduled task matched any exactly time in a week.
// 创建一个定时任务,匹配一周中的任意时间,但时间必须是明确的
func NewWeeklySchedule(weekday, hour, minute, second int8, name string, task func(time.Time)) {
	NewSchedule(ANY, ANY, weekday, hour, minute, second, name, task)
}

// NewDailySchedule creates a new scheduled task matched any exactly time of day.
// 创建一个定时任务,匹配一天中的任意时间,但时间必须是明确的
func NewDailySchedule(hour, minute, second int8, name string, task func(time.Time)) {
	NewSchedule(ANY, ANY, ANY, hour, minute, second, name, task)
}

//Matches will match now time to execute task function
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
		locker.RLock()
		for _, j := range tasks {
			// execute all our cron tasks asynchronously
			if j.Matches(now) {
				go j.Task(now)
			}
		}
		locker.RUnlock()
		time.Sleep(time.Second)
	}
}

func init() {
	tasks = make(map[string]schedule)
	go processSchedules()
}
