package primer

import (
	"fmt"
	"time"
)

func Repeat() {
	for i := 0; i < 5; i++ {
		fmt.Println("Hello from Docker!")
		time.Sleep(time.Second)
	}
}
