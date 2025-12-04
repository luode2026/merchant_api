package main

import (
	"fmt"
	"merchant_api/internal/pkg/utils"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: go run scripts/gen_hash.go <密码>")
		fmt.Println("示例: go run scripts/gen_hash.go admin888")
		fmt.Println("\n常用密码哈希:")

		passwords := []string{"admin888", "password123", "123456"}
		for _, pwd := range passwords {
			hash, _ := utils.HashPassword(pwd)
			fmt.Printf("  密码: %-15s => 哈希: %s\n", pwd, hash)
		}
		return
	}

	password := os.Args[1]
	hash, err := utils.HashPassword(password)
	if err != nil {
		fmt.Printf("❌ 加密失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 密码加密成功\n")
	fmt.Printf("原始密码: %s\n", password)
	fmt.Printf("Bcrypt 哈希: %s\n", hash)
	fmt.Printf("\nSQL 更新语句:\n")
	fmt.Printf("UPDATE mer_merchant_admin SET pwd = '%s' WHERE account = 'YOUR_ACCOUNT';\n", hash)
}
