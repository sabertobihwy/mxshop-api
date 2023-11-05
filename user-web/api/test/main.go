package main

import (
	"fmt"
	"golang.org/x/exp/rand"
	"strings"
	"time"
)

func generateSMSCode(len int) string {
	array := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Seed(uint64(time.Now().UnixNano()))
	var sb strings.Builder
	for i := 0; i < len; i++ {
		sb.WriteString(fmt.Sprintf("%d", array[rand.Intn(len)]))
	}
	return sb.String()
}

func main() {
	fmt.Println(generateSMSCode(5))

}
