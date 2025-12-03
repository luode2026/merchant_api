package main

import (
	"flag"
	"fmt"
	"merchant_api/pkg/config"
	"merchant_api/pkg/database"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gen"
)

func main() {
	// 定义命令行参数
	var (
		tables string
	)

	flag.StringVar(&tables, "table", "", "指定要生成的表名，多个表用逗号分隔，为空则生成所有表")
	flag.Parse()

	// 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	// 初始化数据库连接
	if err := database.InitMySQL(cfg.Database.MySQL); err != nil {
		panic(fmt.Sprintf("初始化数据库失败: %v", err))
	}

	fmt.Println("🚀 开始生成模型文件...")
	fmt.Println("📦 数据库:", cfg.Database.MySQL.Database)
	if tables != "" {
		fmt.Println("📋 指定表:", tables)
	} else {
		fmt.Println("📋 所有表")
	}
	fmt.Println("📂 Model 目录: ./internal/model")
	fmt.Println("📂 DAO 目录: ./internal/dao")
	fmt.Println("")

	// 创建生成器
	g := gen.NewGenerator(gen.Config{
		OutPath:           "./internal/dao",   // DAO 代码输出目录
		OutFile:           "gen.go",           // 查询代码文件名
		ModelPkgPath:      "./internal/model", // Model 代码输出目录
		WithUnitTest:      false,
		FieldNullable:     true,
		FieldCoverable:    false,
		FieldSignable:     false,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		Mode:              gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	// 使用数据库连接
	g.UseDB(database.GetDB())

	// 自定义 JSON 标签格式（使用蛇形命名）
	g.WithJSONTagNameStrategy(func(columnName string) string {
		return columnName // 保持原始列名
	})

	// 自定义文件名策略（移除 .gen 后缀）
	g.WithFileNameStrategy(func(tableName string) string {
		return tableName
	})

	// 确定要生成的模型
	var models []interface{}

	if tables != "" {
		// 生成指定表
		tableList := strings.Split(tables, ",")
		for _, tableName := range tableList {
			tableName = strings.TrimSpace(tableName)
			if tableName == "" {
				continue
			}
			models = append(models, g.GenerateModel(tableName))
		}
	} else {
		// 生成所有表
		models = g.GenerateAllTable()
	}

	// 应用模型配置
	g.ApplyBasic(models...)

	// 执行生成
	g.Execute()

	// 重命名文件（移除 .gen 后缀，DAO 文件添加 _dao 后缀）
	if err := removeGenSuffix("./internal/model"); err != nil {
		fmt.Printf("⚠️ 重命名 Model 文件失败: %v\n", err)
	}
	if err := renameDaoFiles("./internal/dao"); err != nil {
		fmt.Printf("⚠️ 重命名 DAO 文件失败: %v\n", err)
	}

	fmt.Println("")
	fmt.Println("✅ 模型生成成功！")
	fmt.Println("📁 Model 目录: ./internal/model")
	fmt.Println("📁 DAO 目录: ./internal/dao")
}

// renameDaoFiles 重命名 DAO 文件：移除 .gen 后缀并添加 _dao 后缀
func renameDaoFiles(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		// 跳过 gen.go 文件
		if name == "gen.go" {
			continue
		}

		if strings.HasSuffix(name, ".gen.go") {
			oldPath := filepath.Join(dir, name)
			// 移除 .gen.go，添加 _dao.go
			baseName := strings.TrimSuffix(name, ".gen.go")
			newName := baseName + "_dao.go"
			newPath := filepath.Join(dir, newName)

			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
			fmt.Printf("📝 重命名: %s -> %s\n", name, newName)
		}
	}
	return nil
}

// removeGenSuffix 移除目录下所有文件的 .gen 后缀
func removeGenSuffix(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if strings.HasSuffix(name, ".gen.go") {
			oldPath := filepath.Join(dir, name)
			newPath := filepath.Join(dir, strings.Replace(name, ".gen.go", ".go", 1))
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
			fmt.Printf("📝 重命名: %s -> %s\n", name, filepath.Base(newPath))
		}
	}
	return nil
}
