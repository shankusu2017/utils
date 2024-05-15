package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestIPConv(t *testing.T) {
	fmt.Printf(os.Getwd())
	InitPac("../OpenVPN-Golang/config/cnIP.cfg")
	if IsLocalIP("47.103.204.157") != true {
		t.Fatalf("8a48b8f6 FATAL")
	}
	if IsLocalIP("8.8.8.8") {
		t.Fatalf("666d1174 FATAL")
	}
	if IsLocalIP("116.63.131.209") != true {
		t.Fatalf("0x08c817a8 FATAL")
	}
	if IsLocalIP("64.176.173.228") == true {
		t.Fatalf("0x1ecbf52d FATAL")
	}
	if IsLocalIP("64.176.85.253") == true {
		t.Fatalf("0x3eb80066 FATAL")
	}
	if IsLocalIP("182.92.205.225") != true {
		t.Fatalf("0x5c68b8eb FATAL")
	}
}
