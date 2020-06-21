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

	NewDailySchedule(12, 0, 0, func(t time.Time) {
		fmt.Printf("现在时间是%s,定时执行中", t)
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
