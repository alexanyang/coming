package coming

import (
	"fmt"
	"testing"
	"time"
)

func TestNewSchedule(t *testing.T) {

	NewDailySchedule(ANY, ANY, ANY, "小任务", func(t time.Time) {
		fmt.Printf("现在时间是%s,定时执行中\n", t)
	})
	time.Sleep(5 * time.Second)
}

func ExampleNewDailySchedule() {
	NewDailySchedule(12, 12, 0, "DailySchedule", func(t time.Time) {
		fmt.Println("Example of NewDailySchedule(...) ")
	})
}

func ExampleNewMonthlySchedule() {
	NewMonthlySchedule(1, 0, 0, 0, "MonthlySchedule", func(t time.Time) {
		fmt.Println("This beginning of a month")
	})
}

func ExampleNewWeeklySchedule() {
	NewWeeklySchedule(0, 0, 0, 0, "WeeklySchedule", func(t time.Time) {
		fmt.Println("It is Sunday comming")
	})
}

func ExampleNewSchedule() {
	NewSchedule(1, ANY, ANY, 1, 1, 1, "ExampleNewSchedule", func(t time.Time) {
		fmt.Println("It will run at 01:01:01 every day in January ")
	})
}
