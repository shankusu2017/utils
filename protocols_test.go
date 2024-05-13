package utils

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"testing"
)

func TestProtoDecode(t *testing.T) {
	// Ethernet + IPv4 + UDP + DNS
	etherBuf := []byte{0x08, 0x9b, 0x4b, 0x00, 0xfb, 0x27, 0xd8, 0xf2, 0xca, 0x56, 0x6f, 0x7b, 0x08, 0x00, 0x45, 0x00,
		0x00, 0x3f, 0x7d, 0xeb, 0x00, 0x00, 0x40, 0x11, 0x45, 0x8a, 0x0a, 0x0a, 0x0a, 0x9d, 0xda, 0x06,
		0xc8, 0x8b, 0xdd, 0x2f, 0x00, 0x35, 0x00, 0x2b, 0x8e, 0xfb, 0x45, 0x97, 0x01, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x78, 0x08, 0x72, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63,
		0x68, 0x02, 0x71, 0x71, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01}

	packet := gopacket.NewPacket(etherBuf, layers.LayerTypeEthernet, gopacket.Default)

	// 判断数据包是否为以太网数据包，可解析出源mac地址、目的mac地址、以太网类型（如ip类型）等
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		fmt.Println("检测到以太网数据包")
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
		fmt.Println("源 Mac 地址为：", ethernetPacket.SrcMAC)
		fmt.Println("目的 Mac地址为：", ethernetPacket.DstMAC)
		//以太网类型通常是 IPv4，但也可以是 ARP 或其他
		fmt.Println("以太网类型为：", ethernetPacket.EthernetType)
		fmt.Println(ethernetPacket.Payload)
		fmt.Println()
	}

	// 判断数据包是否为IP数据包，可解析出源ip、目的ip、协议号等
	ipLayer := packet.Layer(layers.LayerTypeIPv4) //这里抓取ipv4的数据包
	if ipLayer != nil {
		fmt.Println("检测到IP层数据包")
		ip, _ := ipLayer.(*layers.IPv4)
		// IP layer variables:
		// Version (Either 4 or 6)
		// IHL (IP Header Length in 32-bit words)
		// TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
		//Checksum,SrcIP, DstIP
		fmt.Printf("源ip为 %d 目的ip为 %d\n", ip.SrcIP, ip.DstIP)
		fmt.Println("协议版本为：", ip.Version)
		fmt.Println("首部长度为:", ip.IHL)
		fmt.Println("区分服务为：", ip.TOS)
		fmt.Println("总长度为:", ip.Length)
		fmt.Println("标识id为：", ip.Id)
		fmt.Println("标志为：", ip.Flags)
		fmt.Println("片偏移", ip.FragOffset)
		fmt.Println("TTL", ip.TTL)
		fmt.Println("协议Protocol为：", ip.Protocol)
		fmt.Println("校验和为：", ip.Checksum)
		fmt.Println("基础层", ip.BaseLayer)
		fmt.Println("内容contents：", ip.Contents)
		fmt.Println("可选字段为：", ip.Options)
		fmt.Println("填充为：", ip.Padding)
		fmt.Println()
	}

	// 判断数据包是否为TCP数据包，可解析源端口、目的端口、seq序列号、tcp标志位等
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		fmt.Println("检测到tcp数据包")
		tcp, _ := tcpLayer.(*layers.TCP)
		// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
		// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
		fmt.Printf("源端口为 %d 目的端口为 %d\n", tcp.SrcPort, tcp.DstPort)
		fmt.Println("序列号为：", tcp.Seq)
		fmt.Println("确认号为：", tcp.Ack)
		fmt.Println("数据偏移量：", tcp.DataOffset)
		fmt.Println("标志位：", tcp.CWR, tcp.ECE, tcp.URG, tcp.ACK, tcp.PSH, tcp.RST, tcp.SYN, tcp.FIN)
		fmt.Println("窗口大小：", tcp.Window)
		fmt.Println("校验值：", tcp.Checksum)
		fmt.Println("紧急指针：", tcp.Urgent)
		fmt.Println("tcp选项：", tcp.Options)
		fmt.Println("填充：", tcp.Padding)
		fmt.Println()
	}

	// 判断数据包是否为TCP数据包，可解析源端口、目的端口、seq序列号、tcp标志位等
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		fmt.Println("检测到udp数据包")
		udp, _ := udpLayer.(*layers.UDP)
		// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
		// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
		fmt.Printf("源端口为 %d 目的端口为 %d\n", udp.SrcPort, udp.DstPort)
		fmt.Println()
	}
}

func TestCalSessionID(t *testing.T) {
	etherBuf := []byte{0x08, 0x9b, 0x4b, 0x00, 0xfb, 0x27, 0xd8, 0xf2, 0xca, 0x56, 0x6f, 0x7b, 0x08, 0x00, 0x45, 0x00,
		0x00, 0x3f, 0x7d, 0xeb, 0x00, 0x00, 0x40, 0x11, 0x45, 0x8a, 0x0a, 0x0a, 0x0a, 0x9d, 0xda, 0x06,
		0xc8, 0x8b, 0xdd, 0x2f, 0x00, 0x35, 0x00, 0x2b, 0x8e, 0xfb, 0x45, 0x97, 0x01, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x78, 0x08, 0x72, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63,
		0x68, 0x02, 0x71, 0x71, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01}
	ipBuf := []byte{0x45, 0x00,
		0x00, 0x3f, 0x7d, 0xeb, 0x00, 0x00, 0x40, 0x11, 0x45, 0x8a, 0x0a, 0x0a, 0x0a, 0x9d, 0xda, 0x06,
		0xc8, 0x8b, 0xdd, 0x2f, 0x00, 0x35, 0x00, 0x2b, 0x8e, 0xfb, 0x45, 0x97, 0x01, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x78, 0x08, 0x72, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63,
		0x68, 0x02, 0x71, 0x71, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01}

	id, err := CalSessionID(etherBuf)
	if err != nil {
		log.Printf("err:%s", err)
	}

	id, _ = CalSessionID(ipBuf)
	log.Printf(id)
}
