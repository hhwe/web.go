// color.go make stdout colorful
// more infomation in https://my.oschina.net/dingdayu/blog/1537064
package main

import "fmt"

const (
	RED = uint8(iota + 91)
	GREEN
	YELLOW
	BLUE
)

func red(s string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", RED, s)
}

func green(s string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", GREEN, s)
}

func yellow(s string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", YELLOW, s)
}

func blue(s string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", BLUE, s)
}
