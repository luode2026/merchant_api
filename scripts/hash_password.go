package main

import (
	"fmt"
	"merchant_api/internal/pkg/utils"
)

func main() {
	// 测试密码加密和验证
	password := "admin888"

	// 生成密码哈希
	hash, err := utils.HashPassword(password)
	if err != nil {
		fmt.Printf("加密失败: %v\n", err)
		return
	}

	fmt.Printf("原始密码: %s\n", password)
	fmt.Printf("加密后的哈希: %s\n", hash)
	fmt.Printf("哈希长度: %d\n", len(hash))

	// 验证密码
	isValid := utils.CheckPassword(password, hash)
	fmt.Printf("密码验证结果: %v\n", isValid)

	// 测试另一个密码
	password2 := "password123"
	hash2, _ := utils.HashPassword(password2)
	fmt.Printf("\n原始密码: %s\n", password2)
	fmt.Printf("加密后的哈希: %s\n", hash2)
	fmt.Printf("哈希长度: %d\n", len(hash2))
}
