package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func runInteractiveMenu() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("===============================================")
		fmt.Printf("  Zpanel 命令行管理工具\n  Version: %s\n", Version)
		fmt.Println("===============================================")
		fmt.Println("  1. 查看面板入口信息")
		fmt.Println("  2. 修改面板密码")
		fmt.Println("  3. 修改面板端口")
		fmt.Println("  4. 启动面板")
		fmt.Println("  5. 停止面板")
		fmt.Println("  6. 重启面板")
		fmt.Println("  7. 查看 LNMP 状态")
		fmt.Println("  8. 一键安装 LNMP")
		fmt.Println("  9. 查看站点列表")
		fmt.Println("  0. 退出")
		fmt.Println("===============================================")
		fmt.Print("请输入命令编号: ")
		line, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(line)
		var err error
		switch choice {
		case "1":
			cfg, e := loadConfig()
			if e != nil {
				err = e
			} else {
				printPanelURL(cfg)
			}
		case "2":
			fmt.Print("新密码: ")
			pw, _ := reader.ReadString('\n')
			err = newUserPasswordCmd().RunE(nil, []string{strings.TrimSpace(pw)})
		case "3":
			fmt.Print("新端口: ")
			p, _ := reader.ReadString('\n')
			err = newPortSetCmd().RunE(nil, []string{strings.TrimSpace(p)})
		case "4":
			err = newStartCmd().RunE(nil, nil)
		case "5":
			err = newStopCmd().RunE(nil, nil)
		case "6":
			err = newRestartCmd().RunE(nil, nil)
		case "7":
			newLNMPStatusCmd().Run(nil, nil)
		case "8":
			err = newLNMPInstallCmd().RunE(nil, nil)
		case "9":
			err = newSiteListCmd().RunE(nil, nil)
		case "0", "q", "exit":
			return
		default:
			fmt.Println("无效选项")
			continue
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		}
		fmt.Println()
	}
}
