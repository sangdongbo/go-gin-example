package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func TestMap(c *gin.Context) {
	// 1. 创建 map 的几种方式
	// 方式1：使用 make 创建空 map
	userAges := make(map[string]int)
	userAges["Alice"] = 25
	userAges["Bob"] = 30
	userAges["Charlie"] = 35

	// 方式2：直接初始化 map
	fruits := map[string]float64{
		"apple":  2.5,
		"banana": 1.8,
		"orange": 3.2,
	}

	// 7. 嵌套 map 的正确使用方法
	// 方法1：逐步初始化
	testMap1 := make(map[string]map[int]int)

	// 必须先初始化内层 map
	testMap1["classOne"] = make(map[int]int)
	testMap1["classOne"][1] = 10
	testMap1["classOne"][2] = 20

	testMap1["classTwo"] = make(map[int]int)
	testMap1["classTwo"][1] = 30
	testMap1["classTwo"][2] = 40

	fmt.Println("嵌套 map 内容:", testMap1)

	// 方法2：直接初始化
	testMap2 := map[string]map[int]int{
		"classOne": {
			1: 10,
			2: 20,
		},
		"classTwo": {
			1: 30,
			2: 40,
		},
	}

	// 方法3：安全的添加方法（
	testMap3 := make(map[string]map[int]int)

	// 安全添加函数
	addToNestedMap := func(outerKey string, innerKey, value int) {
		if testMap3[outerKey] == nil {
			testMap3[outerKey] = make(map[int]int)
		}
		testMap3[outerKey][innerKey] = value
	}

	addToNestedMap("math", 1, 95)
	addToNestedMap("math", 2, 87)
	addToNestedMap("english", 1, 92)
	addToNestedMap("english", 2, 88)

	// 遍历嵌套 map
	fmt.Println("\n=== 遍历嵌套 map ===")
	for className, scores := range testMap3 {
		fmt.Printf("科目: %s\n", className)
		for studentId, score := range scores {
			fmt.Printf("  学生%d: %d分\n", studentId, score)
		}
	}

	// 方式3：创建空 map 并添加元素
	colors := map[int]string{}
	colors[1] = "red"
	colors[2] = "green"
	colors[3] = "blue"

	testColors := map[int]string{}
	testColors[1] = "testOne"
	testColors[2] = "testTwo"
	testColors[3] = "testThree"

	fmt.Println(testColors)

	// 2. 读取 map 值
	aliceAge := userAges["Alice"]

	// 安全读取（检查键是否存在）
	var bobStatus string
	var bobAge int
	if age, ok := userAges["Bob"]; ok {
		bobAge = age
		bobStatus = "found"
	} else {
		bobStatus = "not found"
	}

	// 3. 修改 map 值
	userAges["Alice"] = 26 // 修改已存在的键
	userAges["David"] = 28 // 添加新键

	// 4. 删除 map 元素
	delete(userAges, "Charlie")

	// 5. 遍历 map
	var userList []map[string]interface{}
	for name, age := range userAges {
		userList = append(userList, map[string]interface{}{
			"name": name,
			"age":  age,
		})
	}

	mapTestThree := make(map[int]string)
	mapTestThree[1] = "testOne"
	mapTestThree[2] = "testTwo"
	mapTestThree[3] = "testThree"

	if threeOneValue, ok := mapTestThree[2]; ok {
		fmt.Println("Key 2 exists with value:", threeOneValue)
	} else {
		fmt.Println("Key 2 does not exist.")
	}

	// 6. 检查 map 长度
	mapLength := len(userAges)

	// 7. 嵌套 map 示例
	students := map[string]map[string]interface{}{
		"student1": {
			"name":  "张三",
			"age":   20,
			"grade": "A",
		},
		"student2": {
			"name":  "李四",
			"age":   21,
			"grade": "B",
		},
	}

	// 8. map 作为计数器使用
	words := []string{"apple", "banana", "apple", "orange", "banana", "apple"}
	wordCount := make(map[string]int)
	for _, word := range words {
		wordCount[word]++
	}

	// 9. 清空 map（重新赋值为空 map）
	emptyMap := make(map[string]int)
	emptyMap["test"] = 1
	emptyMap = make(map[string]int) // 清空

	c.JSON(200, gin.H{
		"message": "Go Map 常用操作示例",
		"examples": gin.H{
			// 基本操作
			"1_basic_maps": gin.H{
				"user_ages": userAges,
				"fruits":    fruits,
				"colors":    colors,
			},

			// 读取操作
			"2_reading": gin.H{
				"alice_age":  aliceAge,
				"bob_age":    bobAge,
				"bob_status": bobStatus,
			},

			// 遍历结果
			"3_iteration": gin.H{
				"user_list":  userList,
				"map_length": mapLength,
			},

			// 嵌套 map 示例
			"4_nested_maps": gin.H{
				"students":   students,
				"test_map_1": testMap1,
				"test_map_2": testMap2,
				"test_map_3": testMap3,
			},

			// 计数器用法
			"5_counter": gin.H{
				"original_words": words,
				"word_count":     wordCount,
			},

			// 实用技巧
			"6_tips": gin.H{
				"empty_map_length": len(emptyMap),
				"map_is_reference": "map 是引用类型，传递时不会复制",
				"zero_value":       "map 的零值是 nil",
				"key_types":        "键必须是可比较的类型（string, int, bool 等）",
			},
		},

		// 常见错误和解决方案
		"common_mistakes": gin.H{
			"mistake_1": gin.H{
				"error":        "panic: assignment to entry in nil map",
				"wrong_code":   "m := make(map[string]map[int]int); m[\"key\"][1] = 1",
				"correct_code": "m := make(map[string]map[int]int); m[\"key\"] = make(map[int]int); m[\"key\"][1] = 1",
			},
			"mistake_2": gin.H{
				"error":        "读取不存在的键",
				"wrong_code":   "value := m[\"nonexistent\"]",
				"correct_code": "value, exists := m[\"nonexistent\"]; if exists { ... }",
			},
		},

		// 常用模式总结
		"common_patterns": []string{
			"1. 使用 make(map[K]V) 创建空 map",
			"2. 使用 value, ok := map[key] 安全读取",
			"3. 使用 delete(map, key) 删除元素",
			"4. 使用 for key, value := range map 遍历",
			"5. 使用 len(map) 获取元素数量",
			"6. map 可以作为计数器、缓存、索引使用",
		},
	})
}
