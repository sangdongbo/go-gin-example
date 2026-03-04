package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {
	s := make([]int, 0)
	oldCap := cap(s)
	for i := 0; i < 2000; i++ {
		s = append(s, i)
		if newCap := cap(s); newCap != oldCap {
			fmt.Printf("扩容: %d -> %d\n", oldCap, newCap)
			oldCap = newCap
		}
	}

	var a *int = nil
	var b interface{} = a
	if b == nil {
		fmt.Println("b is nil")
	} else {
		fmt.Println("b is NOT nil") // 实际执行这一行，因为 b 包含了 *int 类型信息
	}
}
