package v1

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func TestOneContext(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	ch := make(chan string, 1)
	go func() {
		time.Sleep(3 * time.Second) // 模拟慢操作
		ch <- "操作完成"
	}()

	select {
	case result := <-ch:
		c.JSON(200, gin.H{"message": result})
	case <-ctx.Done():
		c.JSON(504, gin.H{"error": "请求超时"})
	}
}

func LimitRate(c *gin.Context) {
	limiter := rate.NewLimiter(2, 3)

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for i := 1; i <= 10; i++ {
		<-ticker.C // 模拟每 200ms 有一个请求到达

		ctx, cancel := context.WithTimeout(c.Request.Context(), 300*time.Millisecond)
		err := limiter.Wait(ctx)
		cancel()

		if err != nil {
			fmt.Printf("请求 %d 被限流或超时：%v\n", i, err)
		} else {
			fmt.Printf("处理请求 %d at %s\n", i, time.Now().Format("15:04:05.000"))
		}
	}
}

func LimitRate1(c *gin.Context) {
	limiter := rate.NewLimiter(2, 3)

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for i := 1; i <= 20; i++ {
		<-ticker.C

		ctx, cancel := context.WithTimeout(c.Request.Context(), 300*time.Millisecond)
		err := limiter.Wait(ctx)
		cancel()

		if err != nil {
			fmt.Printf("请求 %d 被限流或超时：%v\n", i, err)
		} else {
			fmt.Printf("处理请求 %d at %s\n", i, time.Now().Format("15:04:05.000"))
		}
	}
}

func callServiceA(resultChan chan<- string) {
	time.Sleep(1 * time.Second)
	resultChan <- "serviceA"
}

func callServiceB(resultChan chan<- string) {
	time.Sleep(1500 * time.Millisecond)
	resultChan <- "serviceB"
}

func Aggregate(c *gin.Context) {
	startTime := time.Now()
	resultChan := make(chan string)
	go callServiceA(resultChan)
	go callServiceB(resultChan)

	resultA := <-resultChan
	resultB := <-resultChan

	endTime := time.Now()

	wasteTime := endTime.Sub(startTime).Seconds()
	formatted := fmt.Sprintf("%.2f", wasteTime)

	// 聚合结果并返回
	c.JSON(http.StatusOK, gin.H{
		"serviceA":  resultA,
		"serviceB":  resultB,
		"wasteTime": formatted,
	})
}

func externalService(
	ctx context.Context,
	name string,
	delay time.Duration,
	succeed bool,
	resultChan chan<- string,
	errChan chan<- error,
) {
	select {
	case <-ctx.Done():
		// 主动取消或超时
		errChan <- fmt.Errorf("%s canceled", name)
	case <-time.After(delay):
		if succeed {
			resultChan <- fmt.Sprintf("%s 成功返回", name)
		} else {
			errChan <- fmt.Errorf("%s 失败", name)
		}
	}
}

func Fastest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	resultChan := make(chan string, 1)
	errChan := make(chan error, 1)

	go externalService(ctx, "服务A", 1*time.Second, false, resultChan, errChan)
	go externalService(ctx, "服务B", 800*time.Millisecond, true, resultChan, errChan)
	go externalService(ctx, "服务C", 1*time.Second, false, resultChan, errChan)

	var result string
	var errCount int

loop:
	for {
		select {
		case res := <-resultChan:
			result = res
			break loop
		case <-errChan:
			errCount++
			if errCount == 3 {
				break loop
			}
		case <-ctx.Done():
			result = "超时未返回任何结果"
			break loop
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

func goWorker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s 收到取消信号，退出\n", name)
			return
		default:
			fmt.Printf("%s 正在工作...\n", name)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func TestContextTimeout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	go goWorker(ctx, "workOne")
	go goWorker(ctx, "workTwo")

	time.Sleep(3 * time.Second)
	fmt.Println("worker 进程结束")
}

func fetchData(ctx context.Context, name string, delay time.Duration, resultChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	select {
	case <-ctx.Done():
		fmt.Printf("[%s] 超时或取消，放弃执行\n", name)
	case <-time.After(delay):
		result := fmt.Sprintf("[%s] 数据返回（耗时 %.1fs）", name, delay.Seconds())
		resultChan <- result
	}
}

func Aggregate2(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	resultChan := make(chan string, 3)

	wg.Add(3)
	go fetchData(ctx, "用户服务", 1*time.Second, resultChan, &wg)
	go fetchData(ctx, "订单服务", 3*time.Second, resultChan, &wg)
	go fetchData(ctx, "积分服务", 500*time.Millisecond, resultChan, &wg)

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var results []string
	for res := range resultChan {
		results = append(results, res)
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
	})
}

func doTask(id int) {
	sleepTime := time.Duration(rand.Intn(1000)+500) * time.Millisecond
	time.Sleep(sleepTime)
	fmt.Printf("任务 #%d 完成，耗时 %.2fs\n", id, sleepTime.Seconds())
}

func DoTask1(c *gin.Context) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5)

	taskCount := 6
	for i := 0; i < taskCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() {
				<-semaphore
			}()
			doTask(id)
		}(i)
	}
	wg.Wait()
	fmt.Println("所有任务执行完毕")
}

func Lock1(c *gin.Context) {
	var mutex sync.Mutex
	var count int

	for i := 0; i < 1000; i++ {
		go func() {
			mutex.Lock()
			count++
			mutex.Unlock()
		}()
	}

	fmt.Println("count:", count)
}
