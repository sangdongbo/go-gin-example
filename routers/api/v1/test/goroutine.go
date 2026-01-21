package v1

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func TestGoroutine(c *gin.Context) {
	// 1. 基本 goroutine 使用
	results := []string{}
	results = append(results, "=== 基本 Goroutine ===")

	// 创建一个 channel 来收集结果
	basicChan := make(chan string, 5)

	// 启动多个 goroutine
	for i := 1; i <= 3; i++ {
		go func(id int) {
			time.Sleep(time.Millisecond * 100)
			basicChan <- fmt.Sprintf("Goroutine %d 完成", id)
		}(i)
	}

	// 收集结果
	for i := 0; i < 3; i++ {
		results = append(results, <-basicChan)
	}

	basicChan1 := make(chan int, 10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			time.Sleep(time.Millisecond * 100)
			basicChan1 <- i
		}(i)
	}

	for i := 0; i < 10; i++ {
		fmt.Println(<-basicChan1)
	}

	// 2. 无缓冲 channel（同步）
	unbufferedChan := make(chan string)

	go func() {
		unbufferedChan <- "无缓冲通道消息"
	}()

	unbufferedMsg := <-unbufferedChan

	// 3. 有缓冲 channel（异步）
	bufferedChan := make(chan int, 3)
	bufferedChan <- 1
	bufferedChan <- 2
	bufferedChan <- 3

	bufferedResults := []int{}
	bufferedResults = append(bufferedResults, <-bufferedChan)
	bufferedResults = append(bufferedResults, <-bufferedChan)
	bufferedResults = append(bufferedResults, <-bufferedChan)

	// 4. channel 关闭和 range 遍历
	numberChan := make(chan int, 5)

	go func() {
		for i := 1; i <= 5; i++ {
			numberChan <- i
		}
		close(numberChan) // 关闭 channel
	}()

	rangeResults := []int{}
	for num := range numberChan {
		rangeResults = append(rangeResults, num)
	}

	// 5. select 多路复用
	chan1 := make(chan string, 1)
	chan2 := make(chan string, 1)

	go func() {
		time.Sleep(time.Millisecond * 50)
		chan1 <- "来自 channel 1"
	}()

	go func() {
		time.Sleep(time.Millisecond * 100)
		chan2 <- "来自 channel 2"
	}()

	selectResults := []string{}
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-chan1:
			selectResults = append(selectResults, msg1)
		case msg2 := <-chan2:
			selectResults = append(selectResults, msg2)
		case <-time.After(time.Millisecond * 200):
			selectResults = append(selectResults, "超时")
		}
	}

	// 6. WaitGroup 等待多个 goroutine
	var wg sync.WaitGroup
	wgResults := make(chan string, 5)

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(time.Millisecond * 50)
			wgResults <- fmt.Sprintf("Worker %d 完成任务", id)
		}(i)
	}

	// 等待所有 goroutine 完成后关闭 channel
	go func() {
		wg.Wait()
		close(wgResults)
	}()

	wgResultList := []string{}
	for result := range wgResults {
		wgResultList = append(wgResultList, result)
	}

	// 7. Worker Pool 模式
	jobs := make(chan int, 10)
	poolResults := make(chan string, 10)

	// 创建 3 个 worker
	numWorkers := 3
	var poolWg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		poolWg.Add(1)
		go func(workerID int) {
			defer poolWg.Done()
			for job := range jobs {
				time.Sleep(time.Millisecond * 20)
				poolResults <- fmt.Sprintf("Worker %d 处理了 Job %d", workerID, job)
			}
		}(w)
	}

	// 发送任务
	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs)

	// 等待所有 worker 完成
	go func() {
		poolWg.Wait()
		close(poolResults)
	}()

	poolResultList := []string{}
	for result := range poolResults {
		fmt.Println("------", result)
		poolResultList = append(poolResultList, result)
	}

	// 8. 单向 channel（只读/只写）
	dataChan := make(chan int, 3)

	// 只写 channel
	go func(ch chan<- int) {
		for i := 1; i <= 3; i++ {
			ch <- i * 10
		}
		close(ch)
	}(dataChan)

	// 只读 channel
	readOnlyResults := []int{}
	func(ch <-chan int) {
		for val := range ch {
			readOnlyResults = append(readOnlyResults, val)
		}
	}(dataChan)

	// 9. 互斥锁 Mutex， 竞态条件
	var mutex sync.Mutex
	counter := 0
	var counterWg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		counterWg.Add(1)
		go func() {
			defer counterWg.Done()
			mutex.Lock()
			counter++
			mutex.Unlock()
		}()
	}

	counterWg.Wait()

	// 10. 读写锁 RWMutex
	var rwMutex sync.RWMutex
	data := make(map[string]int)
	data["count"] = 0

	var rwWg sync.WaitGroup

	// 1 个写入者
	rwWg.Add(1)
	go func() {
		defer rwWg.Done()
		for i := 0; i < 5; i++ {
			rwMutex.Lock()
			data["count"]++
			rwMutex.Unlock()
			time.Sleep(time.Millisecond * 10)
		}
	}()

	// 3 个读取者
	readResults := make(chan int, 15)
	for r := 0; r < 3; r++ {
		rwWg.Add(1)
		go func(readerID int) {
			defer rwWg.Done()
			for i := 0; i < 5; i++ {
				rwMutex.RLock()
				val := data["count"]
				rwMutex.RUnlock()
				readResults <- val
				time.Sleep(time.Millisecond * 5)
			}
		}(r)
	}

	rwWg.Wait()
	close(readResults)

	readResultList := []int{}
	for val := range readResults {
		readResultList = append(readResultList, val)
	}

	// 11. Once - 只执行一次
	var once sync.Once
	onceResults := []string{}
	onceChan := make(chan string, 5)

	for i := 0; i < 5; i++ {
		go func(id int) {
			once.Do(func() {
				onceChan <- fmt.Sprintf("只有 Goroutine %d 执行了这段代码", id)
			})
		}(i)
	}

	time.Sleep(time.Millisecond * 100)
	close(onceChan)

	for msg := range onceChan {
		onceResults = append(onceResults, msg)
	}

	// 12. Context 超时控制（简化示例）
	timeoutChan := make(chan string, 1)
	doneChan := make(chan bool, 1)

	go func() {
		time.Sleep(time.Millisecond * 50)
		select {
		case <-doneChan:
			return
		default:
			timeoutChan <- "任务完成"
		}
	}()

	contextResult := ""
	select {
	case res := <-timeoutChan:
		contextResult = res
	case <-time.After(time.Millisecond * 100):
		contextResult = "任务超时"
		doneChan <- true
	}

	timeoutChane1 := make(chan string, 1)
	doneChan1 := make(chan bool, 1)

	go func() {
		time.Sleep(time.Millisecond * 50)
		select {
		case <-doneChan1:
			return
		default:
			timeoutChane1 <- "任务完成"
		}
	}()

	contextResult1 := ""
	select {
	case res := <-timeoutChane1:
		contextResult1 = res
	case <-time.After(time.Millisecond * 100):
		contextResult1 = "任务超时"
		doneChan1 <- true
	}

	fmt.Println(contextResult1)

	// 13. Pipeline 模式
	// 第一阶段：生成数字
	generator := func() <-chan int {
		out := make(chan int)
		go func() {
			for i := 1; i <= 5; i++ {
				out <- i
			}
			close(out)
		}()
		return out
	}

	// 第二阶段：平方
	square := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			for num := range in {
				out <- num * num
			}
			close(out)
		}()
		return out
	}

	// 第三阶段：转换为字符串
	toString := func(in <-chan int) <-chan string {
		out := make(chan string)
		go func() {
			for num := range in {
				out <- fmt.Sprintf("结果: %d", num)
			}
			close(out)
		}()
		return out
	}

	// 构建 pipeline
	pipelineResults := []string{}
	for result := range toString(square(generator())) {
		pipelineResults = append(pipelineResults, result)
	}

	// 14. Fan-Out Fan-In 模式
	fanInSource := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		fanInSource <- i
	}
	close(fanInSource)

	// Fan-Out: 分发到多个 worker
	fanOutWorker := func(id int, in <-chan int) <-chan string {
		out := make(chan string)
		go func() {
			for num := range in {
				out <- fmt.Sprintf("Worker %d 处理: %d", id, num)
			}
			close(out)
		}()
		return out
	}

	// Fan-In: 合并多个 channel
	fanIn := func(channels ...<-chan string) <-chan string {
		out := make(chan string)
		var fanWg sync.WaitGroup

		for _, ch := range channels {
			fanWg.Add(1)
			go func(c <-chan string) {
				defer fanWg.Done()
				for msg := range c {
					out <- msg
				}
			}(ch)
		}

		go func() {
			fanWg.Wait()
			close(out)
		}()

		return out
	}

	// 创建 3 个 worker
	w1 := fanOutWorker(1, fanInSource)
	w2 := fanOutWorker(2, fanInSource)
	w3 := fanOutWorker(3, fanInSource)

	// 合并结果
	fanInResults := []string{}
	for result := range fanIn(w1, w2, w3) {
		fanInResults = append(fanInResults, result)
	}

	fmt.Println("并发操作演示完成")

	c.JSON(200, gin.H{
		"message": "Go Goroutine & Channel 并发操作示例",
		"examples": gin.H{
			// 基本使用
			"1_basic_goroutine": gin.H{
				"results": results,
				"note":    "使用 go 关键字启动 goroutine",
			},

			// Channel 类型
			"2_channels": gin.H{
				"unbuffered": gin.H{
					"message": unbufferedMsg,
					"note":    "无缓冲 channel 是同步的，发送和接收必须同时准备好",
				},
				"buffered": gin.H{
					"results": bufferedResults,
					"note":    "有缓冲 channel 可以在缓冲区未满时异步发送",
				},
				"range": gin.H{
					"results": rangeResults,
					"note":    "使用 range 遍历 channel，直到 channel 关闭",
				},
			},

			// Select
			"3_select": gin.H{
				"results": selectResults,
				"note":    "select 可以同时等待多个 channel 操作",
			},

			// WaitGroup
			"4_waitgroup": gin.H{
				"results": wgResultList,
				"note":    "WaitGroup 用于等待一组 goroutine 完成",
			},

			// Worker Pool
			"5_worker_pool": gin.H{
				"workers":    numWorkers,
				"total_jobs": 9,
				"results":    poolResultList,
				"note":       "Worker Pool 模式用于限制并发数量",
			},

			// 单向 Channel
			"6_directional_channel": gin.H{
				"results": readOnlyResults,
				"note":    "chan<- 只写，<-chan 只读，增强类型安全",
			},

			// 互斥锁
			"7_mutex": gin.H{
				"final_counter": counter,
				"expected":      100,
				"note":          "Mutex 用于保护共享资源，防止竞态条件",
			},

			// 读写锁
			"8_rwmutex": gin.H{
				"final_count":  data["count"],
				"read_results": readResultList,
				"note":         "RWMutex 允许多个读操作并发，但写操作互斥",
			},

			// Once
			"9_once": gin.H{
				"results": onceResults,
				"note":    "sync.Once 确保函数只执行一次，常用于单例初始化",
			},

			// Context
			"10_context": gin.H{
				"result": contextResult,
				"note":   "使用 timeout 控制 goroutine 的执行时间",
			},

			// Pipeline
			"11_pipeline": gin.H{
				"results": pipelineResults,
				"note":    "Pipeline 模式通过 channel 串联多个处理阶段",
			},

			// Fan-Out Fan-In
			"12_fan_in_out": gin.H{
				"results": fanInResults,
				"note":    "Fan-Out 分发任务，Fan-In 聚合结果",
			},
		},

		// 常见错误
		"common_mistakes": gin.H{
			"mistake_1": gin.H{
				"error":        "panic: send on closed channel",
				"wrong_code":   "close(ch); ch <- 1",
				"correct_code": "发送完成后再 close，或使用 defer close(ch)",
			},
			"mistake_2": gin.H{
				"error":        "goroutine 泄漏",
				"wrong_code":   "启动 goroutine 但永不退出",
				"correct_code": "使用 context 或 done channel 通知退出",
			},
			"mistake_3": gin.H{
				"error":        "deadlock: all goroutines are asleep",
				"wrong_code":   "ch := make(chan int); ch <- 1; <-ch",
				"correct_code": "使用 goroutine 或缓冲 channel",
			},
			"mistake_4": gin.H{
				"error":        "竞态条件 (race condition)",
				"wrong_code":   "多个 goroutine 直接修改共享变量",
				"correct_code": "使用 Mutex 或 channel 同步",
			},
			"mistake_5": gin.H{
				"error":        "忘记关闭 channel",
				"note":         "range 遍历 channel 时，如果不关闭会死锁",
				"correct_code": "发送方负责关闭 channel",
			},
		},

		// 最佳实践
		"best_practices": []string{
			"1. 由发送方关闭 channel，接收方不应关闭",
			"2. 不要向已关闭的 channel 发送数据",
			"3. 可以从已关闭的 channel 接收数据（返回零值）",
			"4. 使用 WaitGroup 等待多个 goroutine 完成",
			"5. 使用 select 的 default 分支避免阻塞",
			"6. 使用 defer 确保 Mutex 解锁和 channel 关闭",
			"7. goroutine 必须有明确的退出机制",
			"8. 使用 -race 标志检测竞态条件: go run -race main.go",
			"9. channel 传递数据，Mutex 保护数据",
			"10. 优先使用 channel，它更符合 Go 的并发哲学",
		},

		// 性能提示
		"performance_tips": gin.H{
			"tip_1": "使用缓冲 channel 减少 goroutine 阻塞",
			"tip_2": "Worker Pool 模式控制并发数，避免创建过多 goroutine",
			"tip_3": "读多写少场景使用 RWMutex 提高性能",
			"tip_4": "避免在循环中创建大量 goroutine",
			"tip_5": "使用 sync.Pool 复用对象，减少 GC 压力",
		},

		// 常用模式
		"common_patterns": gin.H{
			"worker_pool":  "限制并发数量，处理大量任务",
			"pipeline":     "数据流式处理，每个阶段由 goroutine 处理",
			"fan_out":      "并行处理，提高吞吐量",
			"fan_in":       "聚合多个来源的数据",
			"timeout":      "使用 select + time.After 实现超时",
			"quit_channel": "使用 done channel 通知 goroutine 退出",
		},
	})
}
