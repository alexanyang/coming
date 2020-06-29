package coming

import (
	"fmt"
	"testing"
	"time"
)

/**
 * @author anyang
 * Email: 1300378587@qq.com
 * Created Date:2020-06-21 11:35
 */

func TestNewSchedule(t *testing.T) {

	NewDailySchedule(ANY, ANY, ANY, func(t time.Time) {
		fmt.Printf("现在时间是%s,定时执行中\n", t)
	})
	time.Sleep(5 * time.Second)
}

func ExampleNewDailySchedule() {
	NewDailySchedule(12, 12, 0, func(t time.Time) {
		fmt.Println("Example of NewDailySchedule(...) ")
	})
}

func ExampleNewMonthlySchedule() {
	NewMonthlySchedule(1, 0, 0, 0, func(t time.Time) {
		fmt.Println("This beginning of a month")
	})
}

func ExampleNewWeeklySchedule() {
	NewWeeklySchedule(0, 0, 0, 0, func(t time.Time) {
		fmt.Println("It is Sunday comming")
	})
}

func ExampleNewSchedule() {
	NewSchedule(1, ANY, ANY, 1, 1, 1, func(t time.Time) {
		fmt.Println("It will run at 01:01:01 every day in January ")
	})
}
