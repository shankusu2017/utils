package utils

// 参考资料 https://blog.csdn.net/youuzi/article/details/134059230

import (
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// ProtoHead packet 中各层协议的头
type protoHead struct {
	IPv4 *layers.IPv4
	UDP  *layers.UDP
	TCP  *layers.TCP
}

// DecodePacket 解码 buf, 返回可能的 ether, ip, icmp udp, tcp, 层数据
func decodePacket(ipBuf []byte) *protoHead {
	var head protoHead

	// 假设是 Ipv4开头
	packet := gopacket.NewPacket(ipBuf, layers.LayerTypeIPv4, gopacket.Default)

	// 判断数据包是否为IP数据包，可解析出源ip、目的ip、协议号等
	ipLayer := packet.Layer(layers.LayerTypeIPv4) //这里抓取ipv4的数据包
	if ipLayer != nil {
		//fmt.Println("检测到IP层数据包")
		ip, _ := ipLayer.(*layers.IPv4)
		if ip.Version != 4 {
			return &head
		}
		head.IPv4 = ip
		// IP layer variables:
		// Version (Either 4 or 6)
		// IHL (IP Header Length in 32-bit words)
		// TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
		//Checksum,SrcIP, DstIP
		//fmt.Printf("源ip为 %d 目的ip为 %d\n", ip.SrcIP, ip.DstIP)
		//fmt.Println("协议版本为：", ip.Version)
		//fmt.Println("首部长度为:", ip.IHL)
		//fmt.Println("区分服务为：", ip.TOS)
		//fmt.Println("总长度为:", ip.Length)
		//fmt.Println("标识id为：", ip.Id)
		//fmt.Println("标志为：", ip.Flags)
		//fmt.Println("片偏移", ip.FragOffset)
		//fmt.Println("TTL", ip.TTL)
		//fmt.Println("协议Protocol为：", ip.Protocol)
		//fmt.Println("校验和为：", ip.Checksum)
		//fmt.Println("基础层", ip.BaseLayer)
		//fmt.Println("内容contents：", ip.Contents)
		//fmt.Println("可选字段为：", ip.Options)
		//fmt.Println("填充为：", ip.Padding)
		//fmt.Println()
	}

	// 判断数据包是否为TCP数据包，可解析源端口、目的端口、seq序列号、tcp标志位等
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		//fmt.Println("检测到tcp数据包")
		tcp, _ := tcpLayer.(*layers.TCP)
		head.TCP = tcp
		// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
		// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
		//fmt.Printf("源端口为 %d 目的端口为 %d\n", tcp.SrcPort, tcp.DstPort)
		//fmt.Println("序列号为：", tcp.Seq)
		//fmt.Println("确认号为：", tcp.Ack)
		//fmt.Println("数据偏移量：", tcp.DataOffset)
		//fmt.Println("标志位：", tcp.CWR, tcp.ECE, tcp.URG, tcp.ACK, tcp.PSH, tcp.RST, tcp.SYN, tcp.FIN)
		//fmt.Println("窗口大小：", tcp.Window)
		//fmt.Println("校验值：", tcp.Checksum)
		//fmt.Println("紧急指针：", tcp.Urgent)
		//fmt.Println("tcp选项：", tcp.Options)
		//fmt.Println("填充：", tcp.Padding)
		//fmt.Println()
	}

	// 判断数据包是否为TCP数据包，可解析源端口、目的端口、seq序列号、tcp标志位等
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		//fmt.Println("检测到udp数据包")
		udp, _ := udpLayer.(*layers.UDP)
		head.UDP = udp
		//// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
		//// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
		//fmt.Printf("源端口为 %d 目的端口为 %d\n", udp.SrcPort, udp.DstPort)
		//fmt.Println()
	}

	return &head
}

// CalSrcSessionID 计算一个IP会话的ID（类似TCP、UDP的五元组）
// 从源端发往接收端
func CalSrcSessionID(ipBuf []byte) (string, error) {
	var id = ""

	head := decodePacket(ipBuf)
	if head.IPv4 == nil {
		return id, errors.New(fmt.Sprintf("0x46d6f691 ipv4 is nil"))
	}

	if head.UDP != nil {
		id = fmt.Sprintf("ipSrc_%s:ipDst_%s:protocol_%d:srcPort_%d:dstPort_%d",
			head.IPv4.SrcIP.String(), head.IPv4.DstIP.String(), head.IPv4.Protocol,
			head.UDP.SrcPort, head.UDP.DstPort)
	} else if head.TCP != nil {
		id = fmt.Sprintf("ipSrc_%s:ipDst_%s:protocol_%d:srcPort_%d:dstPort_%d",
			head.IPv4.SrcIP.String(), head.IPv4.DstIP.String(), head.IPv4.Protocol,
			head.TCP.SrcPort, head.TCP.DstPort)
	} else {
		id = fmt.Sprintf("ipSrc_%s:ipDst_%s:protocol_%d",
			head.IPv4.SrcIP.String(), head.IPv4.DstIP.String(), head.IPv4.Protocol)
	}

	return id, nil
}

// CalDstSessionID 计算一个IP会话的ID（类似TCP、UDP的五元组）
// 从接收端发往源端
func CalDstSessionID(ipBuf []byte) (string, error) {
	var id = ""

	head := decodePacket(ipBuf)
	if head.IPv4 == nil {
		return id, errors.New(fmt.Sprintf("0x46d6f691 ipv4 is nil"))
	}

	if head.UDP != nil {
		id = fmt.Sprintf("ipSrc_%s:ipDst_%s:protocol_%d:srcPort_%d:dstPort_%d",
			head.IPv4.DstIP.String(), head.IPv4.SrcIP.String(), head.IPv4.Protocol,
			head.UDP.DstPort, head.UDP.SrcPort)
	} else if head.TCP != nil {
		id = fmt.Sprintf("ipSrc_%s:ipDst_%s:protocol_%d:srcPort_%d:dstPort_%d",
			head.IPv4.DstIP.String(), head.IPv4.SrcIP.String(), head.IPv4.Protocol,
			head.TCP.DstPort, head.TCP.SrcPort)
	} else {
		id = fmt.Sprintf("ipSrc_%s:ipDst_%s:protocol_%d",
			head.IPv4.DstIP.String(), head.IPv4.SrcIP.String(), head.IPv4.Protocol)
	}

	return id, nil
}
