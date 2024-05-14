package utils

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestIPConv(t *testing.T) {
	fmt.Printf(os.Getwd())
	InitPac("../config/cnIP.cfg")
	if IsLocalIP("47.103.204.157") != true {
		t.Fatalf("8a48b8f6 FATAL")
	}
	if IsLocalIP("8.8.8.8") {
		t.Fatalf("666d1174 FATAL")
	}
	log.Printf("isLocalIP:%v\n", IsLocalIP("182.92.205.225"))
}
