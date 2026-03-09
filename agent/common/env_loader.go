package common

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"fmt"
)

func init() {
	// 加载顺序：
	// 1. 当前目录 .env
	// 2. 上级目录 .env (通常在 agent/ 目录下)
	loaded1 := tryLoad(".env")
	loaded2 := tryLoad(filepath.Join("..", ".env"))
	
	if loaded1 || loaded2 {
		fmt.Printf("[ENV] 环境变量加载完成。模型 ID: %s\n", os.Getenv("VOLC_ENDPOINT_ID"))
	} else {
		fmt.Println("[ENV] 警告：未找到任何 .env 文件！")
	}
}

func tryLoad(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	
	fmt.Printf("[ENV] 正在从 %s 加载环境变量...\n", path)
	
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		if k == "" {
			continue
		}
		// 强制覆盖：如果文件里有新值，优先使用文件里的值
		// 之前的逻辑是 if os.Getenv(k) == "" { ... }，这意味着如果系统环境变量已存在（比如上次运行残留或手动设置），就不会更新
		// 现在改为总是覆盖，确保 .env 文件是最高优先级
		_ = os.Setenv(k, v)
	}
	return true
}

