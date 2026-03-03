package test

import (
	"testing"
	"yiarce/core/yorm"
)

// TestPageCalculation 测试分页计算功能
func TestPageCalculation(t *testing.T) {
	// 创建一个测试用的配置
	config := yorm.Config{
		Host:      "127.0.0.1",
		Port:      "3306",
		Database:  "test_db",
		Username:  "root",
		Password:  "root",
		Character: "utf8mb4",
	}

	// 测试Page方法的分页计算
	db, err := yorm.ConnMysql(config)
	if err != nil {
		// 数据库连接失败，跳过测试
		t.Skipf("Database connection failed: %v", err)
	}
	defer db.Close()

	// 测试分页计算
	page := 2
	size := 10
	_ = db.Page(page, size)

	// 测试Page方法的边界条件
	_ = db.Page(0, 10) // 这应该会返回错误但不会panic
	_ = db.Page(1, 0)  // 这应该会返回错误但不会panic
}

// TestWhereMethod 测试Where方法的功能
func TestWhereMethod(t *testing.T) {
	// 创建一个测试用的配置
	config := yorm.Config{
		Host:      "127.0.0.1",
		Port:      "3306",
		Database:  "test_db",
		Username:  "root",
		Password:  "root",
		Character: "utf8mb4",
	}

	// 测试Where方法
	db, err := yorm.ConnMysql(config)
	if err != nil {
		// 数据库连接失败，跳过测试
		t.Skipf("Database connection failed: %v", err)
	}
	defer db.Close()

	// 测试Where方法的不同用法
	// 1. Where("id", 1) => id = 1
	_ = db.Where("id", 1)

	// 2. Where("id", ">", 1) => id > 1
	_ = db.Where("id", ">", 1)

	// 3. Where("id = 1 AND name = 'test'") => id = 1 AND name = 'test'
	_ = db.Where("id = 1 AND name = 'test'")

	// 4. 测试空条件
	_ = db.Where("")

	// 5. 测试无效的操作符类型
	_ = db.Where("id", 123, 456) // 这应该会返回错误
}

// TestLimitMethod 测试Limit方法的功能
func TestLimitMethod(t *testing.T) {
	// 创建一个测试用的配置
	config := yorm.Config{
		Host:      "127.0.0.1",
		Port:      "3306",
		Database:  "test_db",
		Username:  "root",
		Password:  "root",
		Character: "utf8mb4",
	}

	// 测试Limit方法
	db, err := yorm.ConnMysql(config)
	if err != nil {
		// 数据库连接失败，跳过测试
		t.Skipf("Database connection failed: %v", err)
	}
	defer db.Close()

	// 测试正常情况
	_ = db.Limit(10)

	// 测试边界条件
	_ = db.Limit(0) // 这应该会返回错误
}

// TestTableMethod 测试Table方法的功能
func TestTableMethod(t *testing.T) {
	// 创建一个测试用的配置
	config := yorm.Config{
		Host:      "127.0.0.1",
		Port:      "3306",
		Database:  "test_db",
		Username:  "root",
		Password:  "root",
		Character: "utf8mb4",
	}

	// 测试Table方法
	db, err := yorm.ConnMysql(config)
	if err != nil {
		// 数据库连接失败，跳过测试
		t.Skipf("Database connection failed: %v", err)
	}
	defer db.Close()

	// 测试正常情况
	_ = db.Table("users")
	_ = db.Table("users", "u")

	// 测试边界条件
	_ = db.Table("")          // 这应该会返回错误
	_ = db.Table("users", "") // 这应该会返回错误
}
