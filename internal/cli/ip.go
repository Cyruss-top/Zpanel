package cli

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"

	"github.com/zex/zpanel/internal/config"
)

var ipv4Re = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)

var publicIPServices = []string{
	"https://api.ipify.org",
	"https://ifconfig.me/ip",
	"https://ip.sb",
	"https://icanhazip.com",
	"https://myip.ipip.net",
}

func queryPublicIP() string {
	for _, u := range publicIPServices {
		out, err := exec.Command("curl", "-fsSL", "--connect-timeout", "3", u).Output()
		if err != nil {
			continue
		}
		for _, m := range ipv4Re.FindAllString(string(out), -1) {
			if ip := net.ParseIP(m); ip != nil && !isPrivateIP(ip) {
				return m
			}
		}
	}
	return ""
}

func localIPs() []string {
	out, err := exec.Command("hostname", "-I").Output()
	if err != nil {
		return nil
	}
	return strings.Fields(string(out))
}

func firstLocalIP() string {
	ips := localIPs()
	if len(ips) > 0 {
		return ips[0]
	}
	return ""
}

func isPrivateIP(ip net.IP) bool {
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast()
}

func serverIP() string {
	if ip := queryPublicIP(); ip != "" {
		return ip
	}
	for _, s := range localIPs() {
		if ip := net.ParseIP(s); ip != nil && !isPrivateIP(ip) {
			return s
		}
	}
	if ip := firstLocalIP(); ip != "" {
		return ip
	}
	return "127.0.0.1"
}

func printPanelURL(cfg *config.Config) {
	publicIP := queryPublicIP()
	localIP := firstLocalIP()
	displayIP := publicIP
	label := "面板地址"
	if displayIP == "" {
		displayIP = localIP
		if displayIP == "" {
			displayIP = "127.0.0.1"
		}
	} else {
		label = "面板地址(公网)"
	}

	entry := cfg.EntryPrefix()
	path := ""
	if entry != "" {
		path = entry + "/"
	}
	fmt.Printf("%s: http://%s:%d%s\n", label, displayIP, cfg.Panel.Port, path)
	if publicIP != "" && localIP != "" && publicIP != localIP {
		fmt.Printf("面板地址(内网): http://%s:%d%s\n", localIP, cfg.Panel.Port, path)
	}
	fmt.Printf("用户名:   %s\n", cfg.Auth.Username)
	if cfg.Panel.Entry != "" {
		fmt.Printf("安全入口: /%s/\n", cfg.Panel.Entry)
	}
	fmt.Printf("配置文件: %s\n", configPath())
}
