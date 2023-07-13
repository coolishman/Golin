package port

import (
	"bufio"
	"fmt"
	"golin/port/Protocol"
	"net"
	"strings"
	"time"
)

var portProtocols = map[string]string{
	"20":    "文件传输协议（数据）:「FTP-data」",
	"25":    "邮件传输协议:「SMTP」",
	"53":    "域名解析:「DNS」",
	"110":   "邮件传输协议:「POP3」",
	"135":   "RPC 服务",
	"137":   "NetBIOS 名称服务",
	"138":   "NetBIOS 数据报服务",
	"139":   "NetBIOS 会话服务",
	"161":   "网络管理协议:「SNMP」",
	"162":   "网络管理报警和事件提醒:「SNMP-trap」",
	"143":   "邮件传输协议:「IMAP」",
	"445":   "Microsoft 的 SMB 协议",
	"465":   "带有 SSL 安全的 SMTP：「SMTPS」",
	"514":   "系统日志服务:「Syslog」",
	"587":   "邮件提交协议（MSA）:「Submission」",
	"993":   "带有 SSL 安全的 IMAP：「IMAPS」",
	"995":   "带有 SSL 安全的 POP3：「POP3S」",
	"1024":  "保留给操作系统的端口（通常用于动态或私有端口）",
	"1080":  "SOCKS 代理",
	"1194":  "开放VPN",
	"1433":  "数据库:「SQL Server」",
	"1521":  "数据库:「Oracle」",
	"2049":  "网络文件系统:「NFS」",
	"3389":  "远程桌面协议:「RDP」",
	"5601":  "ES数据库管理系统:「Kibana」",
	"5900":  "虚拟网络计算:「VNC」",
	"9200":  "数据库:「ElasticSearch」",
	"9300":  "数据库:「ElasticSearch集群管理端口」",
	"11211": "缓存服务:「Memcached」",
	"27077": "数据库:「MondoDB」",
}

// parseProtocol 协议/组件分析：有的基于默认端口去对应服务
func parseProtocol(conn net.Conn, host, port string) string {

	if protocol, ok := portProtocols[port]; ok {
		return protocol
	}

	if err := conn.SetReadDeadline(time.Now().Add(time.Duration(Timeout) * time.Second)); err != nil {
		return ""
	}
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		line = ""
	}

	switch {
	case Protocol.IsSSHProtocol(line):
		return strings.ReplaceAll(strings.ReplaceAll(line, "\r", ""), "\n", "")

	case strings.HasPrefix(line, "220"):
		return "文件传输协议（控制）:「FTP」"

	case Protocol.IsRedisProtocol(conn):
		return "数据库:「Redis」"

	case Protocol.IsTelnet(conn):
		return "Telnet"

	case Protocol.IsMySqlProtocol(host, port):
		return "数据库:「MySQL」"

	case Protocol.IsPgsqlProtocol(host, port):
		return "数据库:「PostgreSQL」"

	default:
		isWeb := Protocol.IsWeb(host, port)
		if isWeb != "" {
			return fmt.Sprintf("Web应用 %s", isWeb)
		}
	}

	return ""
}