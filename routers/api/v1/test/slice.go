package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func TestSlice(c *gin.Context) {
	// 1. 创建 slice 的几种方式
	// 方式1：字面量创建
	array1 := [5]int{1, 2, 3, 4, 5}
	fmt.Println("Array:", array1)

	fruits := []string{"apple", "banana", "orange"}
	fruits = append(fruits, "mango")
	fruits[2] = "kiwi"
	fmt.Println("Fruits Slice:", fruits)

	// 方式2：使用 make 创建（指定长度）
	numbers := make([]int, 5)    // 长度为5，容量为5，元素默认为0
	scores := make([]int, 3, 10) // 长度为3，容量为10
	scores[0] = 85
	scores[1] = 90
	scores[2] = 95

	numbers1 := make([]int, 10, 20) // 长度为10，容量为20
	scores1 := make([]int, 2, 4)
	scores1[0] = 100
	scores1[1] = 98
	fmt.Println(numbers1, scores1)

	// 方式3：使用 make 创建空 slice（长度为0）
	names := make([]string, 0, 5) // 长度为0，容量为5
	names = append(names, "Alice", "Bob", "Charlie")

	// 方式4：从数组创建 slice
	arr := [5]int{1, 2, 3, 4, 5}
	arrSlice := arr[1:4] // [2, 3, 4]

	// 方式5：nil slice
	var nilSlice []int // nil slice，长度和容量都是0

	// 2. 读取和修改 slice
	firstFruit := fruits[0] // 读取第一个元素
	fruits[1] = "mango"     // 修改第二个元素

	// 安全访问（避免越界）
	var lastFruit string
	if len(fruits) > 0 {
		lastFruit = fruits[len(fruits)-1]
	}

	// 3. append 操作
	// 单个元素追加
	fruits = append(fruits, "grape")

	// 多个元素追加
	fruits = append(fruits, "pear", "peach")

	// 追加另一个 slice
	moreFruits := []string{"kiwi", "melon"}
	fruits = append(fruits, moreFruits...)

	// 4. 切片操作
	sliceDemo := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	slice1 := sliceDemo[2:5] // [2, 3, 4]
	slice2 := sliceDemo[:3]  // [0, 1, 2]
	slice3 := sliceDemo[7:]  // [7, 8, 9]
	slice4 := sliceDemo[:]   // 完整副本

	// 5. copy 操作
	source := []int{1, 2, 3, 4, 5}
	dest := make([]int, len(source))
	copyCount := copy(dest, source) // 复制元素

	// 部分复制
	partialDest := make([]int, 3)
	copy(partialDest, source) // 只复制前3个

	// 6. len 和 cap
	testSlice := make([]int, 5, 10)
	sliceLen := len(testSlice) // 5
	sliceCap := cap(testSlice) // 10

	// 7. 遍历 slice
	var fruitList []map[string]interface{}
	for index, fruit := range fruits {
		fruitList = append(fruitList, map[string]interface{}{
			"index": index,
			"name":  fruit,
		})
	}

	for index, fruit := range fruits {
		fmt.Printf("Fruit %d: %s\n", index, fruit)
	}

	for i := 0; i < len(fruits); i++ {
		fmt.Printf("Fruit %d: %s\n", i, fruits[i])
	}

	for _, fruit1 := range fruits {
		fmt.Printf("Fruit %s\n", fruit1)
	}

	// 只遍历值
	total := 0
	for _, num := range numbers {
		total += num
	}

	// 8. 删除元素
	// 删除索引为2的元素
	deleteIndex := 2
	deleteDemo := []string{"a", "b", "c", "d", "e"}
	deleteDemo = append(deleteDemo[:deleteIndex], deleteDemo[deleteIndex+1:]...)

	// 9. 插入元素
	insertIndex := 2
	insertValue := "x"
	insertDemo := []string{"a", "b", "c", "d"}
	insertDemo = append(insertDemo[:insertIndex], append([]string{insertValue}, insertDemo[insertIndex:]...)...)

	// 10. 清空 slice
	clearDemo := []int{1, 2, 3, 4, 5}
	clearDemo = clearDemo[:0] // 清空但保留容量
	// 或 clearDemo = nil            // 彻底清空

	// 11. 嵌套 slice
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	// 动态构建嵌套 slice
	students := make([][]string, 3)
	students[0] = []string{"Alice", "Math", "A"}
	students[1] = []string{"Bob", "English", "B"}
	students[2] = []string{"Charlie", "Physics", "A+"}

	students1 := make([][]string, 3)
	students1[0] = []string{"David", "Chemistry", "B+"}
	students1[1] = []string{"Eva", "Biology", "A"}
	fmt.Println(students1)

	// 12. slice 扩容机制演示
	capacityDemo := make([]int, 0)
	capacityChanges := []map[string]int{}
	for i := 0; i < 20; i++ {
		oldCap := cap(capacityDemo)
		capacityDemo = append(capacityDemo, i)
		newCap := cap(capacityDemo)
		if oldCap != newCap {
			capacityChanges = append(capacityChanges, map[string]int{
				"length":  len(capacityDemo),
				"old_cap": oldCap,
				"new_cap": newCap,
			})
		}
	}

	// 13. slice 作为函数参数（引用传递）
	modifySlice := func(s []int) {
		if len(s) > 0 {
			s[0] = 999
		}
	}
	testModify := []int{1, 2, 3}
	beforeModify := testModify[0]
	modifySlice(testModify)
	afterModify := testModify[0]
	fmt.Println("Before Modify:", beforeModify, "After Modify:", afterModify)

	// 14. 过滤 slice
	allNumbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	evenNumbers := []int{}
	for _, num := range allNumbers {
		if num%2 == 0 {
			evenNumbers = append(evenNumbers, num)
		}
	}

	// 15. 反转 slice
	reverseDemo := []string{"a", "b", "c", "d", "e"}
	for i, j := 0, len(reverseDemo)-1; i < j; i, j = i+1, j-1 {
		reverseDemo[i], reverseDemo[j] = reverseDemo[j], reverseDemo[i]
	}

	fmt.Println("Slice 操作演示完成")

	c.JSON(200, gin.H{
		"message": "Go Slice 常用操作示例",
		"examples": gin.H{
			// 创建操作
			"1_creation": gin.H{
				"literal_slice": fruits,
				"make_slice":    numbers,
				"make_with_cap": scores,
				"empty_slice":   names,
				"from_array":    arrSlice,
				"nil_slice":     nilSlice,
			},

			// 读取和修改
			"2_access": gin.H{
				"first_fruit": firstFruit,
				"last_fruit":  lastFruit,
				"modified":    fruits,
			},

			// 切片操作
			"3_slicing": gin.H{
				"original":     sliceDemo,
				"slice_2_to_5": slice1,
				"slice_to_3":   slice2,
				"slice_from_7": slice3,
				"full_slice":   slice4,
			},

			// copy 操作
			"4_copy": gin.H{
				"source":       source,
				"dest":         dest,
				"copy_count":   copyCount,
				"partial_dest": partialDest,
			},

			// len 和 cap
			"5_len_and_cap": gin.H{
				"length":           sliceLen,
				"capacity":         sliceCap,
				"capacity_demo":    capacityDemo,
				"capacity_changes": capacityChanges,
			},

			// 遍历
			"6_iteration": gin.H{
				"fruit_list": fruitList,
				"total":      total,
			},

			// 删除和插入
			"7_delete_insert": gin.H{
				"after_delete": deleteDemo,
				"after_insert": insertDemo,
				"cleared":      clearDemo,
			},

			// 嵌套 slice
			"8_nested": gin.H{
				"matrix":   matrix,
				"students": students,
			},

			// 引用传递
			"9_reference": gin.H{
				"before_modify": beforeModify,
				"after_modify":  afterModify,
				"note":          "slice 是引用类型，修改会影响原始数据",
			},

			// 过滤和反转
			"10_operations": gin.H{
				"all_numbers":  allNumbers,
				"even_numbers": evenNumbers,
				"reversed":     reverseDemo,
			},
		},

		// 常见错误和解决方案
		"common_mistakes": gin.H{
			"mistake_1": gin.H{
				"error":        "panic: runtime error: index out of range",
				"wrong_code":   "s := []int{1,2,3}; x := s[5]",
				"correct_code": "if len(s) > 5 { x := s[5] }",
			},
			"mistake_2": gin.H{
				"error":        "删除元素后忘记重新赋值",
				"wrong_code":   "append(s[:i], s[i+1:]...)",
				"correct_code": "s = append(s[:i], s[i+1:]...)",
			},
			"mistake_3": gin.H{
				"error": "nil slice 和 empty slice 的区别",
				"note":  "var s []int 是 nil，[]int{} 和 make([]int, 0) 是 empty",
				"tip":   "通常使用 len(s) == 0 判断空，而不是 s == nil",
			},
			"mistake_4": gin.H{
				"error": "slice 扩容导致底层数组变化",
				"note":  "append 可能创建新数组，原 slice 和新 slice 不再共享底层数组",
				"tip":   "使用 make 预分配容量可以避免频繁扩容",
			},
		},

		// 常用模式总结
		"common_patterns": []string{
			"1. 使用 make([]T, len, cap) 创建指定容量的 slice",
			"2. 使用 append 追加元素，记得重新赋值: s = append(s, elem)",
			"3. 使用 copy(dst, src) 复制 slice，避免共享底层数组",
			"4. 使用 s[:0] 清空 slice 但保留容量",
			"5. 删除元素: s = append(s[:i], s[i+1:]...)",
			"6. 插入元素: s = append(s[:i], append([]T{elem}, s[i:]...)...)",
			"7. 使用 len(s) == 0 判断 slice 是否为空",
			"8. 遍历时使用 range，需要修改元素时使用索引",
		},

		// 性能提示
		"performance_tips": gin.H{
			"tip_1": "预分配容量可以减少内存分配和复制",
			"tip_2": "大 slice 传参时传指针可以避免复制",
			"tip_3": "频繁 append 时提前估算容量",
			"tip_4": "使用 copy 比循环赋值更高效",
		},
	})
}
