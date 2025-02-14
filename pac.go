package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/shankusu2017/constant"
)

/*
 * 先从已知的内、外网地址 map 中查找
 * 找不到，再从特殊的外网地址池查找
 * 最后从普通的内网地址池查找(找不到则属于外网地址，找得到则属于内网地址)
 */
type pacMgrT struct {
	localIPMap map[string]bool // 已知的内网地址 map[ip.addr]=true
	outIPMap   map[string]bool // 已知的外网地址

	outAddrMap   [32]map[string]bool // 特殊的外网地址池信息（命中则属于网外地址，即使符合 localIPSubNet 的条件)
	localAddrMap [32]map[string]bool // 已知的内网地址池信息 [maskNetLen](map[netNo]=exist)
	mtxIPMap     sync.RWMutex
}

var (
	pacMgr *pacMgrT
)

func InitPac(configPath, outAddrCfg string) {
	p := &pacMgrT{}
	p.localIPMap = make(map[string]bool, constant.Size128K)
	p.outIPMap = make(map[string]bool, constant.Size4K)

	initOutCfg(outAddrCfg, p)
	initLocalCfg(configPath, p)

	pacMgr = p
}

func initOutCfg(configPath string, p *pacMgrT) {
	for i := 0; i < len(p.outAddrMap); i++ {
		p.outAddrMap[i] = make(map[string]bool, constant.Size256)
	}

	fd, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("FATAL f981d07e open config file(%s) err(%s)", configPath, err.Error())
	}
	defer fd.Close()

	br := bufio.NewReader(fd)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		strLine := string(a)
		//1.24.0.0/13
		addrInfo := strings.Split(strLine, "/")
		if len(addrInfo) == 2 {
			netNo := addrInfo[0]   /* 网络号 */
			netmask := addrInfo[1] /* 子网掩码 */
			subMask, err := strconv.Atoi(netmask)
			if err != nil {
				log.Fatalf("FATAL 0x14018e82 netmask invalid(%s)\n", netmask)
			} else {
				if subMask >= 0 && subMask < len(p.outAddrMap) {
					p.outAddrMap[subMask][netNo] = true
					//fmt.Printf("%s %s\n", netNo, netmask)
				} else {
					fmt.Printf("WARN 0x13c11d2d netmask invalid(%s)\n", netmask)
				}
			}
		} else {
			log.Fatal(fmt.Sprintf("valid gateWay addr:%s", addrInfo))
		}
	}
}

func initLocalCfg(configPath string, p *pacMgrT) {
	for i := 0; i < len(p.localAddrMap); i++ {
		p.localAddrMap[i] = make(map[string]bool, constant.Size256)
	}

	fd, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("FATAL f981d07e open config file(%s) err(%s)", configPath, err.Error())
	}
	defer fd.Close()

	br := bufio.NewReader(fd)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		strLine := string(a)
		//1.24.0.0/13
		addrInfo := strings.Split(strLine, "/")
		if len(addrInfo) == 2 {
			netNo := addrInfo[0]   /* 网络号 */
			netmask := addrInfo[1] /* 子网掩码 */
			subMask, err := strconv.Atoi(netmask)
			if err != nil {
				log.Fatalf("FATAL 62b5190e netmask invalid(%s)\n", netmask)
			} else {
				if subMask >= 0 && subMask < len(p.localAddrMap) {
					p.localAddrMap[subMask][netNo] = true
					//fmt.Printf("%s %s\n", netNo, netmask)
				} else {
					fmt.Printf("WARN 6b4d85ca netmask invalid(%s)\n", netmask)
				}
			}
		} else {
			log.Fatal(fmt.Sprintf("valid gateWay addr:%s", addrInfo))
		}
	}
}

// 初步判断 ip 的类型，看能否确定是 local-ip，outside-ip,
// 若不能确定则，进一步调用 slowFind 来判断
// IsLocalIP 是本地 ip 吗，若是则不需要发往 proxySrv，反之需要发往代理服务器
func quickFind(ip string) (bool, bool) {
	pacMgr.mtxIPMap.RLock()
	defer pacMgr.mtxIPMap.RUnlock()

	var done = true

	// 已经明确了是 外部 ip(先检查outer,再检查local)
	val := pacMgr.outIPMap[ip]
	if val == true {
		return false, done
	}

	// 已经明确了是 local ip
	val = pacMgr.localIPMap[ip]
	if val == true {
		return true, done
	}

	// 不知道是 local 还是 out ip ，还需要进一步查询
	done = false

	return false, done
}

// local-ip,outside-ip 的 cache均找不到此ip的消息
// 没办法，这里查询数据库来最终确定 ip 类型
func slowFind(ip string) bool {
	pacMgr.mtxIPMap.RLock()
	defer pacMgr.mtxIPMap.RUnlock()

	{ // 被特殊的外网地址池命中
		var out = false
		for i := 1; i < constant.Size32; i++ {
			addr := maskAddr(ip, i)
			out = pacMgr.outAddrMap[i][addr]
			if out == true {
				//log.Printf("subnet.bit:%d, netIP:%v\n", i, addr)
				break
			}
		}
		if out {
			pacMgr.outIPMap[ip] = true
			return false
		}
	}

	var loc = false
	for i := 1; i < constant.Size32; i++ {
		addr := maskAddr(ip, i)
		loc = pacMgr.localAddrMap[i][addr]
		if loc == true {
			//log.Printf("subnet.bit:%d, netIP:%v\n", i, addr)
			break
		}
	}

	if loc {
		pacMgr.localIPMap[ip] = true
	} else {
		pacMgr.outIPMap[ip] = true
	}

	return loc
}

// 192.168.1.2/16 --->192.168.0.0
func maskAddr(ip string, subNet int) string {
	ipv4Addr := net.ParseIP(ip)
	ipv4Mask := net.CIDRMask(subNet, 32)
	return ipv4Addr.Mask(ipv4Mask).To4().String()
}

func IsLocalIP(ip string) bool {
	loc, done := quickFind(ip)
	if done {
		return loc
	}

	loc = slowFind(ip)
	return loc
}
