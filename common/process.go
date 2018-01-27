package common

import (
	"fmt"
	"os"
)

// TagList list all configured tags
func TagList() {
	fmt.Println("Shortcut command configured are below:")
	for tg, cmd := range C.Tags {
		fmt.Println(" ", tg, ":", cmd)
	}
	os.Exit(0)
}

// TagPrint print specified tag configure
func TagPrint(t string) {
	found := false
	for tg, cmd := range C.Tags {
		if tg == t {
			fmt.Println(cmd)
			found = true
			break
		}
	}
	if !found {
		fmt.Println("No such tag: ", t)
	}
	os.Exit(0)
}
