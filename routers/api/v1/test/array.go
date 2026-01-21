package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func TestArray(c *gin.Context) {
	var testArray1 [5]int
	testArray1[0] = 1

	fmt.Println(testArray1)

	weekDays := [7]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	numbers := [5]int{10, 20, 30, 40, 50}
	powersOfTwo := [5]int{1, 2, 4, 8, 16}
	copiedPowers := powersOfTwo
	copiedPowers[0] = 99
	emptyNumbers := [5]int{}

	schedule := [2][3]string{{"Math", "English", "History"}, {"Physics", "Chemistry", "PE"}}

	sliceWindow := numbers[1:4]
	weekSlice := weekDays[2:5]
	sumNumbers := func(arr [5]int) int {
		total := 0
		for _, v := range arr {
			total += v
		}
		return total
	}(numbers)
	averageNumbers := float64(sumNumbers) / float64(len(numbers))

	menuSlice := make([]string, 0, 5)
	menuSlice = append(menuSlice, "Coffee", "Tea")
	menuSlice = append(menuSlice, "Juice")
	menuSlice = append(menuSlice, "Soda")
	slicedMenu := menuSlice[1:]

	dayList := make([]gin.H, 0, len(weekDays))
	for idx, day := range weekDays {
		dayList = append(dayList, gin.H{"index": idx, "day": day})
	}

	c.JSON(200, gin.H{
		"message": "Go Array/Slice 常用操作示例",
		"examples": gin.H{
			"1_creation": gin.H{
				"week_days_manual": weekDays,
				"number_literal":   numbers,
				"powers_of_two":    powersOfTwo,
				"zero_array":       emptyNumbers,
			},
			"2_access": gin.H{
				"first_day":       weekDays[0],
				"last_number":     numbers[len(numbers)-1],
				"numbers_length":  len(numbers),
				"weekdays_length": len(weekDays),
			},
			"3_iteration": gin.H{
				"day_list":    dayList,
				"sum_numbers": sumNumbers,
				"avg_numbers": averageNumbers,
			},
			"4_slices": gin.H{
				"window_from_numbers": sliceWindow,
				"week_slice":          weekSlice,
				"menu_slice":          menuSlice,
				"sliced_menu":         slicedMenu,
			},
			"5_value_semantics": gin.H{
				"original_powers": powersOfTwo,
				"copied_powers":   copiedPowers,
			},
			"6_multi_dimensional": gin.H{"schedule": schedule},
		},
		"tips": gin.H{
			"array_value":  "数组是值类型，赋值会复制底层元素",
			"fixed_length": "数组长度在定义时固定，运行时不可变",
			"slice_view":   "切片共享底层数组，但可以设置不同的长度和容量",
			"index_bounds": "访问时必须在 [0,len) 范围内，否则会 panic",
			"slice_origin": "切片可以从数组、其他切片或字面量直接创建",
			"slice_make":   "make(len, cap) 可以预分配底层数组并初始化长度",
		},
	})
}
