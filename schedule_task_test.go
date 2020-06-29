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
