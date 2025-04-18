package main

import (
	"time"

	rb "github.com/go-vgo/robotgo"
)

func main() {

	// width, height := rb.GetScreenSize()
	// fmt.Printf("屏幕分辨率：%d x %d\n", width, height)
	rb.Scale = true
	// time.Sleep(50 * time.Second)
	rb.Move(600, 1388)
	rb.Click()
	// st1 := "我系香香宝宝龙"
	// st2 := "号，我又复活啦"
	res := "我是黑客"
	for i := 0; i < 50; i++ {
		// res := st1 + strconv.Itoa(i+1) + st2
		// time.Sleep(5 * time.Minute)
		rb.TypeStr(res)
		time.Sleep(400 * time.Millisecond)

		rb.KeyToggle("enter")
	}
}
