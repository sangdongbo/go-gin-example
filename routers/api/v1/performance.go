package v1

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/gin-gonic/gin"
)

// ===== 1. CPU 性能问题：slice 未预分配 =====

// TestCPUSlow 演示未预分配 slice 导致的性能问题
// 访问: GET /api/v1/perf/cpu-slow
func TestCPUSlow(c *gin.Context) {
	appG := app.Gin{C: c}

	start := time.Now()

	// ❌ 问题代码：未预分配容量，会导致多次内存分配和拷贝
	var result []int
	for i := 0; i < 100000; i++ {
		result = append(result, i)
	}

	elapsed := time.Since(start)

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"method":   "未预分配 slice",
		"count":    len(result),
		"duration": elapsed.String(),
		"tip":      "多次扩容和内存拷贝，性能较差",
	})
}

// TestCPUFast 演示预分配 slice 的优化效果
// 访问: GET /api/v1/perf/cpu-fast
func TestCPUFast(c *gin.Context) {
	appG := app.Gin{C: c}

	start := time.Now()

	// ✅ 优化代码：预分配容量，避免多次扩容
	result := make([]int, 0, 100000)
	for i := 0; i < 100000; i++ {
		result = append(result, i)
	}

	elapsed := time.Since(start)

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"method":   "预分配 slice",
		"count":    len(result),
		"duration": elapsed.String(),
		"tip":      "预分配容量，避免扩容，性能提升明显",
	})
}

// ===== 2. 内存/GC 问题：interface{} 滥用 =====

// TestGCPressure 演示 interface{} 导致的 GC 压力
// 访问: GET /api/v1/perf/gc-pressure
func TestGCPressure(c *gin.Context) {
	appG := app.Gin{C: c}

	start := time.Now()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	gcBefore := m.NumGC

	// ❌ 问题代码：大量使用 interface{} 装箱，增加 GC 压力
	data := make([]interface{}, 100000)
	for i := 0; i < 100000; i++ {
		data[i] = i // int 装箱为 interface{}
	}

	// 强制触发 GC 来观察效果
	runtime.GC()
	runtime.ReadMemStats(&m)
	gcAfter := m.NumGC

	elapsed := time.Since(start)

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"method":      "使用 interface{} 装箱",
		"count":       len(data),
		"duration":    elapsed.String(),
		"gc_count":    gcAfter - gcBefore,
		"alloc_mb":    m.Alloc / 1024 / 1024,
		"total_alloc": m.TotalAlloc / 1024 / 1024,
		"tip":         "interface{} 装箱增加堆分配和 GC 压力",
	})
}

// TestGCOptimized 演示使用具体类型减少 GC 压力
// 访问: GET /api/v1/perf/gc-optimized
func TestGCOptimized(c *gin.Context) {
	appG := app.Gin{C: c}

	start := time.Now()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	gcBefore := m.NumGC

	// ✅ 优化代码：使用具体类型，避免装箱
	data := make([]int, 100000)
	for i := 0; i < 100000; i++ {
		data[i] = i
	}

	runtime.GC()
	runtime.ReadMemStats(&m)
	gcAfter := m.NumGC

	elapsed := time.Since(start)

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"method":      "使用具体类型",
		"count":       len(data),
		"duration":    elapsed.String(),
		"gc_count":    gcAfter - gcBefore,
		"alloc_mb":    m.Alloc / 1024 / 1024,
		"total_alloc": m.TotalAlloc / 1024 / 1024,
		"tip":         "避免装箱，减少堆分配和 GC 压力",
	})
}

// ===== 3. 锁竞争问题 =====

var (
	mutexCounter  int64
	mu            sync.Mutex
	atomicCounter int64
)

// TestMutexSlow 演示 Mutex 锁竞争
// 访问: GET /api/v1/perf/mutex-slow
func TestMutexSlow(c *gin.Context) {
	appG := app.Gin{C: c}

	start := time.Now()
	mutexCounter = 0

	// ❌ 问题代码：大量协程竞争同一个 Mutex
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mu.Lock()
				mutexCounter++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	elapsed := time.Since(start)

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"method":   "Mutex 锁",
		"counter":  mutexCounter,
		"duration": elapsed.String(),
		"tip":      "高并发下锁竞争严重，性能较差",
	})
}

// TestAtomicFast 演示原子操作优化
// 访问: GET /api/v1/perf/atomic-fast
func TestAtomicFast(c *gin.Context) {
	appG := app.Gin{C: c}

	start := time.Now()
	atomic.StoreInt64(&atomicCounter, 0)

	// ✅ 优化代码：使用原子操作，无锁并发
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&atomicCounter, 1)
			}
		}()
	}
	wg.Wait()

	elapsed := time.Since(start)
	counter := atomic.LoadInt64(&atomicCounter)

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"method":   "原子操作",
		"counter":  counter,
		"duration": elapsed.String(),
		"tip":      "原子操作无锁，性能显著提升",
	})
}

// ===== 4. Goroutine 泄漏 =====

// TestGoroutineLeak 演示 goroutine 泄漏问题
// 访问: GET /api/v1/perf/goroutine-leak
func TestGoroutineLeak(c *gin.Context) {
	appG := app.Gin{C: c}

	beforeCount := runtime.NumGoroutine()

	// ❌ 问题代码：启动协程但 channel 永不关闭，导致协程泄漏
	ch := make(chan int)
	for i := 0; i < 10; i++ {
		go func() {
			for range ch { // 永远等待，协程泄漏
				// do nothing
			}
		}()
	}

	time.Sleep(100 * time.Millisecond) // 等待协程启动
	afterCount := runtime.NumGoroutine()

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"method":            "未关闭 channel",
		"goroutines_before": beforeCount,
		"goroutines_after":  afterCount,
		"leaked":            afterCount - beforeCount,
		"tip":               "⚠️ 协程永远阻塞，造成泄漏！生产环境要避免",
	})
}

// TestGoroutineClean 演示正确的协程管理
// 访问: GET /api/v1/perf/goroutine-clean
func TestGoroutineClean(c *gin.Context) {
	appG := app.Gin{C: c}

	beforeCount := runtime.NumGoroutine()

	// ✅ 优化代码：使用 context 或关闭 channel 来通知协程退出
	ch := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range ch {
				// process data
			}
		}()
	}

	// 模拟一些工作
	time.Sleep(100 * time.Millisecond)

	// 关闭 channel，通知所有协程退出
	close(ch)
	wg.Wait() // 等待所有协程退出

	afterCount := runtime.NumGoroutine()

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"method":            "正确管理协程生命周期",
		"goroutines_before": beforeCount,
		"goroutines_after":  afterCount,
		"leaked":            afterCount - beforeCount,
		"tip":               "✅ 使用 close(ch) 和 WaitGroup 优雅退出",
	})
}

// ===== 5. 综合性能对比 =====

// TestBenchmark 综合性能对比
// 访问: GET /api/v1/perf/benchmark
func TestBenchmark(c *gin.Context) {
	appG := app.Gin{C: c}

	results := make(map[string]interface{})

	// 1. Slice 性能对比
	results["slice"] = compareSlicePerformance()

	// 2. 锁性能对比
	results["lock"] = compareLockPerformance()

	// 3. 内存分配对比
	results["memory"] = compareMemoryAllocation()

	// 4. 协程信息
	results["goroutines"] = runtime.NumGoroutine()

	// 5. 内存统计
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	results["memory_stats"] = map[string]interface{}{
		"alloc_mb":      m.Alloc / 1024 / 1024,
		"total_alloc":   m.TotalAlloc / 1024 / 1024,
		"sys_mb":        m.Sys / 1024 / 1024,
		"num_gc":        m.NumGC,
		"goroutine_num": runtime.NumGoroutine(),
	}

	appG.Response(http.StatusOK, e.SUCCESS, results)
}

// 辅助函数：对比 slice 性能
func compareSlicePerformance() map[string]interface{} {
	n := 100000

	// 未预分配
	start1 := time.Now()
	var s1 []int
	for i := 0; i < n; i++ {
		s1 = append(s1, i)
	}
	time1 := time.Since(start1)

	// 预分配
	start2 := time.Now()
	s2 := make([]int, 0, n)
	for i := 0; i < n; i++ {
		s2 = append(s2, i)
	}
	time2 := time.Since(start2)

	return map[string]interface{}{
		"without_prealloc": time1.String(),
		"with_prealloc":    time2.String(),
		"speedup":          fmt.Sprintf("%.2fx", float64(time1)/float64(time2)),
	}
}

// 辅助函数：对比锁性能
func compareLockPerformance() map[string]interface{} {
	n := 10000

	// Mutex
	start1 := time.Now()
	var counter1 int64
	var mu1 sync.Mutex
	var wg1 sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg1.Add(1)
		go func() {
			defer wg1.Done()
			for j := 0; j < n; j++ {
				mu1.Lock()
				counter1++
				mu1.Unlock()
			}
		}()
	}
	wg1.Wait()
	time1 := time.Since(start1)

	// Atomic
	start2 := time.Now()
	var counter2 int64
	var wg2 sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			for j := 0; j < n; j++ {
				atomic.AddInt64(&counter2, 1)
			}
		}()
	}
	wg2.Wait()
	time2 := time.Since(start2)

	return map[string]interface{}{
		"mutex":   time1.String(),
		"atomic":  time2.String(),
		"speedup": fmt.Sprintf("%.2fx", float64(time1)/float64(time2)),
	}
}

// 辅助函数：对比内存分配
func compareMemoryAllocation() map[string]interface{} {
	n := 50000

	// interface{} 装箱
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)
	before1 := m1.TotalAlloc

	data1 := make([]interface{}, n)
	for i := 0; i < n; i++ {
		data1[i] = i
	}

	runtime.ReadMemStats(&m1)
	after1 := m1.TotalAlloc
	alloc1 := after1 - before1

	// 具体类型
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)
	before2 := m2.TotalAlloc

	data2 := make([]int, n)
	for i := 0; i < n; i++ {
		data2[i] = i
	}

	runtime.ReadMemStats(&m2)
	after2 := m2.TotalAlloc
	alloc2 := after2 - before2

	return map[string]interface{}{
		"interface_alloc_kb": alloc1 / 1024,
		"concrete_alloc_kb":  alloc2 / 1024,
		"reduction":          fmt.Sprintf("%.2fx", float64(alloc1)/float64(alloc2)),
	}
}

// ===== 额外的性能测试示例 =====

// TestStringConcat 字符串拼接性能对比
func TestStringConcat(c *gin.Context) {
	appG := app.Gin{C: c}

	n := 10000
	results := make(map[string]interface{})

	// 方式1: += 操作符（最慢）
	start1 := time.Now()
	s1 := ""
	for i := 0; i < n; i++ {
		s1 += "x"
	}
	time1 := time.Since(start1)
	results["operator_plus"] = time1.String()

	// 方式2: fmt.Sprintf（中等）
	start2 := time.Now()
	s2 := ""
	for i := 0; i < n; i++ {
		s2 = fmt.Sprintf("%s%s", s2, "x")
	}
	time2 := time.Since(start2)
	results["fmt_sprintf"] = time2.String()

	// 方式3: strings.Builder（最快）
	start3 := time.Now()
	var builder strings.Builder
	builder.Grow(n) // 预分配
	for i := 0; i < n; i++ {
		builder.WriteString("x")
	}
	_ = builder.String() // 使用 _ 忽略返回值
	time3 := time.Since(start3)
	results["strings_builder"] = time3.String()

	results["recommendation"] = "使用 strings.Builder 性能最优"
	results["speedup"] = fmt.Sprintf("Builder 比 += 快 %.0fx", float64(time1)/float64(time3))

	appG.Response(http.StatusOK, e.SUCCESS, results)
}

// TestMapPrealloc Map 预分配测试
func TestMapPrealloc(c *gin.Context) {
	appG := app.Gin{C: c}

	n := 100000

	// 未预分配
	start1 := time.Now()
	m1 := make(map[int]int)
	for i := 0; i < n; i++ {
		m1[i] = i
	}
	time1 := time.Since(start1)

	// 预分配
	start2 := time.Now()
	m2 := make(map[int]int, n)
	for i := 0; i < n; i++ {
		m2[i] = i
	}
	time2 := time.Since(start2)

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"without_prealloc": time1.String(),
		"with_prealloc":    time2.String(),
		"speedup":          fmt.Sprintf("%.2fx", float64(time1)/float64(time2)),
		"tip":              "Map 也支持预分配容量",
	})
}

// TestChannelBuffered Channel 缓冲区测试
func TestChannelBuffered(c *gin.Context) {
	appG := app.Gin{C: c}

	n := 10000

	// 无缓冲 channel
	start1 := time.Now()
	ch1 := make(chan int)
	go func() {
		for i := 0; i < n; i++ {
			ch1 <- i
		}
		close(ch1)
	}()
	count1 := 0
	for range ch1 {
		count1++
	}
	time1 := time.Since(start1)

	// 有缓冲 channel
	start2 := time.Now()
	ch2 := make(chan int, 100)
	go func() {
		for i := 0; i < n; i++ {
			ch2 <- i
		}
		close(ch2)
	}()
	count2 := 0
	for range ch2 {
		count2++
	}
	time2 := time.Since(start2)

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"unbuffered": time1.String(),
		"buffered":   time2.String(),
		"speedup":    fmt.Sprintf("%.2fx", float64(time1)/float64(time2)),
		"tip":        "合理使用缓冲 channel 可减少阻塞",
	})
}

// TestDeferOverhead Defer 开销测试
func TestDeferOverhead(c *gin.Context) {
	appG := app.Gin{C: c}

	n := 1000000

	// 使用 defer
	start1 := time.Now()
	for i := 0; i < n; i++ {
		func() {
			defer func() {}()
		}()
	}
	time1 := time.Since(start1)

	// 不使用 defer
	start2 := time.Now()
	for i := 0; i < n; i++ {
		func() {
			// 直接执行
		}()
	}
	time2 := time.Since(start2)

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"with_defer":    time1.String(),
		"without_defer": time2.String(),
		"overhead":      fmt.Sprintf("defer 有 %.1f%% 的开销", (float64(time1)-float64(time2))/float64(time2)*100),
		"tip":           "defer 有一定开销，但在大多数场景下可以忽略，优先保证代码清晰",
	})
}

// TestRangeVsIndex Range vs Index 性能对比
func TestRangeVsIndex(c *gin.Context) {
	appG := app.Gin{C: c}

	data := make([]int, 100000)
	for i := range data {
		data[i] = rand.Intn(100)
	}

	// range 方式
	start1 := time.Now()
	sum1 := 0
	for _, v := range data {
		sum1 += v
	}
	time1 := time.Since(start1)

	// index 方式
	start2 := time.Now()
	sum2 := 0
	for i := 0; i < len(data); i++ {
		sum2 += data[i]
	}
	time2 := time.Since(start2)

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"range_time": time1.String(),
		"index_time": time2.String(),
		"sum":        sum1,
		"tip":        "range 和 index 性能基本一致，range 代码更简洁",
	})
}
