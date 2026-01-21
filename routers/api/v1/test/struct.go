package v1

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// 1. 基本 struct 定义
type User struct {
	ID       int
	Name     string
	Email    string
	Age      int
	IsActive bool
}

// 2. 带标签的 struct（用于 JSON 序列化）
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description,omitempty"` // omitempty: 空值时不输出
	CreatedAt   string  `json:"created_at"`
}

// 3. 嵌套 struct
type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	ZipCode string `json:"zip_code"`
}

type Employee struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Address Address `json:"address"` // 嵌套 struct
	Salary  float64 `json:"salary"`
}

// 4. 匿名字段（继承）
type Person struct {
	Name string
	Age  int
}

type Student struct {
	Person    // 匿名字段，Student "继承" Person 的字段
	StudentID string
	Grade     string
}

// 5. 方法 - 值接收者
func (u User) GetInfo() string {
	return fmt.Sprintf("%s (%d years old)", u.Name, u.Age)
}

// 6. 方法 - 指针接收者（可以修改 struct）
func (u *User) UpdateEmail(newEmail string) {
	u.Email = newEmail
}

func (u *User) IncrementAge() {
	u.Age++
}

// 7. 构造函数
func NewUser(id int, name, email string, age int) *User {
	return &User{
		ID:       id,
		Name:     name,
		Email:    email,
		Age:      age,
		IsActive: true,
	}
}

func NewProduct(id int, name string, price float64) *Product {
	return &Product{
		ID:        id,
		Name:      name,
		Price:     price,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
}

// 8. 接口实现
type Describer interface {
	Describe() string
}

func (p Product) Describe() string {
	return fmt.Sprintf("%s - $%.2f", p.Name, p.Price)
}

func (e Employee) Describe() string {
	return fmt.Sprintf("%s works at %s", e.Name, e.Address.City)
}

func TestStruct(c *gin.Context) {
	// 1. 创建 struct 的几种方式
	// 方式1：字面量创建
	user1 := User{
		ID:       1,
		Name:     "Alice",
		Email:    "alice@example.com",
		Age:      25,
		IsActive: true,
	}

	// 方式2：省略字段名（必须按顺序）
	user2 := User{2, "Bob", "bob@example.com", 30, true}

	// 方式3：部分字段初始化（其他字段为零值）
	user3 := User{
		ID:   3,
		Name: "Charlie",
	}

	// 方式4：使用 new（返回指针）
	user4 := new(User)
	user4.ID = 4
	user4.Name = "David"
	user4.Email = "david@example.com"

	// 方式5：使用构造函数
	user5 := NewUser(5, "Eve", "eve@example.com", 28)

	// 2. 访问和修改字段
	userName := user1.Name
	user1.Age = 26
	user1.IsActive = false

	// 3. 指针访问（自动解引用）
	userPtr := &user1
	ptrName := userPtr.Name // 等同于 (*userPtr).Name
	userPtr.Age = 27

	// 4. 调用方法
	userInfo := user1.GetInfo()
	user1.UpdateEmail("alice.new@example.com")
	user1.IncrementAge()

	// 5. 嵌套 struct
	emp := Employee{
		ID:   101,
		Name: "John Doe",
		Address: Address{
			Street:  "123 Main St",
			City:    "New York",
			ZipCode: "10001",
		},
		Salary: 75000.50,
	}

	// 访问嵌套字段
	empCity := emp.Address.City
	emp.Address.Street = "456 Oak Ave"

	// 6. 匿名字段（继承）
	student := Student{
		Person: Person{
			Name: "Tom",
			Age:  20,
		},
		StudentID: "S12345",
		Grade:     "A",
	}

	// 可以直接访问匿名字段的属性
	studentName := student.Name // 等同于 student.Person.Name
	student.Age = 21

	// 7. struct 比较
	userA := User{ID: 1, Name: "Alice"}
	userB := User{ID: 1, Name: "Alice"}
	userC := User{ID: 2, Name: "Bob"}
	isEqual := userA == userB    // true
	isNotEqual := userA == userC // false

	// 8. struct 复制（值类型）
	originalUser := User{ID: 1, Name: "Original"}
	copiedUser := originalUser
	copiedUser.Name = "Modified"
	// originalUser.Name 仍然是 "Original"

	// 9. 匿名 struct
	config := struct {
		Host string
		Port int
	}{
		Host: "localhost",
		Port: 8080,
	}

	// 10. struct 切片
	users := []User{user1, user2, user3}
	userPtrs := []*User{user4, user5}

	// 11. struct map
	userMap := map[int]User{
		1: user1,
		2: user2,
		3: user3,
	}

	// 12. JSON 序列化和反序列化
	product := NewProduct(1, "Laptop", 999.99)
	product.Description = "High performance laptop"

	// 序列化为 JSON
	jsonData, _ := json.Marshal(product)
	jsonString := string(jsonData)

	// 反序列化
	var decodedProduct Product
	json.Unmarshal(jsonData, &decodedProduct)

	// 13. 空 struct
	type Empty struct{}
	emptySize := int(0) // Empty struct 大小为 0

	// 14. 遍历 struct 字段（使用反射的简单示例）
	userFields := []string{"ID", "Name", "Email", "Age", "IsActive"}

	// 15. struct 标签使用
	productWithTags := Product{
		ID:    1,
		Name:  "Phone",
		Price: 599.99,
		// Description 为空，JSON 中会省略（omitempty）
		CreatedAt: time.Now().Format("2006-01-02"),
	}
	taggedJSON, _ := json.Marshal(productWithTags)

	// 16. 实现接口
	var describers []Describer
	describers = append(describers, *product)
	describers = append(describers, emp)

	descriptions := []string{}
	for _, d := range describers {
		descriptions = append(descriptions, d.Describe())
	}

	// 17. struct 组合示例
	type Author struct {
		Name  string
		Email string
	}

	type Article struct {
		Title   string
		Content string
		Author  Author // 组合
		Tags    []string
	}

	article := Article{
		Title:   "Go Struct Tutorial",
		Content: "Learn about Go structs...",
		Author: Author{
			Name:  "Jane Smith",
			Email: "jane@example.com",
		},
		Tags: []string{"go", "programming", "tutorial"},
	}

	fmt.Println("Struct 操作演示完成")

	c.JSON(200, gin.H{
		"message": "Go Struct 常用操作示例",
		"examples": gin.H{
			// 创建方式
			"1_creation": gin.H{
				"literal_full":      user1,
				"literal_ordered":   user2,
				"literal_partial":   user3,
				"using_new":         user4,
				"using_constructor": user5,
			},

			// 访问和修改
			"2_access": gin.H{
				"user_name":  userName,
				"modified":   user1,
				"ptr_access": ptrName,
				"ptr_struct": userPtr,
			},

			// 方法调用
			"3_methods": gin.H{
				"get_info":        userInfo,
				"after_update":    user1.Email,
				"after_increment": user1.Age,
			},

			// 嵌套 struct
			"4_nested": gin.H{
				"employee":    emp,
				"emp_city":    empCity,
				"emp_address": emp.Address,
			},

			// 匿名字段
			"5_anonymous": gin.H{
				"student":       student,
				"student_name":  studentName,
				"direct_access": "student.Name 等同于 student.Person.Name",
			},

			// struct 比较和复制
			"6_comparison": gin.H{
				"is_equal":      isEqual,
				"is_not_equal":  isNotEqual,
				"original_user": originalUser,
				"copied_user":   copiedUser,
				"note":          "struct 是值类型，赋值时会完整复制",
			},

			// 匿名 struct
			"7_anonymous_struct": gin.H{
				"config": config,
				"usage":  "适用于临时数据结构",
			},

			// struct 集合
			"8_collections": gin.H{
				"user_slice": users,
				"user_ptrs":  userPtrs,
				"user_map":   userMap,
			},

			// JSON 操作
			"9_json": gin.H{
				"original":       product,
				"json_string":    jsonString,
				"decoded":        decodedProduct,
				"with_omitempty": string(taggedJSON),
			},

			// 接口实现
			"10_interface": gin.H{
				"descriptions": descriptions,
				"note":         "struct 可以实现接口",
			},

			// struct 组合
			"11_composition": gin.H{
				"article":     article,
				"author_name": article.Author.Name,
				"tags":        article.Tags,
			},

			// 其他特性
			"12_others": gin.H{
				"empty_struct_size": emptySize,
				"field_names":       userFields,
			},
		},

		// 常见错误和解决方案
		"common_mistakes": gin.H{
			"mistake_1": gin.H{
				"error":        "值接收者无法修改 struct",
				"wrong_code":   "func (u User) Update() { u.Name = \"new\" }",
				"correct_code": "func (u *User) Update() { u.Name = \"new\" }",
				"note":         "需要修改 struct 时使用指针接收者",
			},
			"mistake_2": gin.H{
				"error":        "未导出的字段无法 JSON 序列化",
				"wrong_code":   "type User struct { name string }",
				"correct_code": "type User struct { Name string }",
				"note":         "字段名首字母大写才能导出",
			},
			"mistake_3": gin.H{
				"error":    "比较包含不可比较字段的 struct",
				"note":     "struct 包含 slice、map 等不可比较类型时无法直接比较",
				"solution": "使用自定义比较函数或 reflect.DeepEqual",
			},
			"mistake_4": gin.H{
				"error":        "忘记初始化嵌套 struct",
				"wrong_code":   "var e Employee; e.Address.City = \"NYC\"",
				"correct_code": "e := Employee{Address: Address{City: \"NYC\"}}",
			},
		},

		// 常用模式总结
		"common_patterns": []string{
			"1. 使用构造函数 NewXxx() 创建 struct，确保正确初始化",
			"2. 导出的字段名首字母大写，私有字段小写",
			"3. 使用指针接收者修改 struct，值接收者读取数据",
			"4. 使用 struct 标签定义 JSON、数据库等映射关系",
			"5. 通过匿名字段实现组合而非继承",
			"6. 使用 omitempty 标签处理可选字段",
			"7. 大 struct 传参时使用指针避免复制",
			"8. 使用匿名 struct 处理临时数据",
		},

		// 最佳实践
		"best_practices": gin.H{
			"naming":       "struct 名称使用大驼峰，如 UserProfile",
			"constructors": "提供 NewXxx() 构造函数初始化复杂 struct",
			"methods":      "相关操作定义为方法，需要修改时用指针接收者",
			"composition":  "优先使用组合而不是继承",
			"tags":         "合理使用 struct 标签简化序列化",
			"zero_value":   "设计时考虑零值是否有意义",
		},
	})
}
