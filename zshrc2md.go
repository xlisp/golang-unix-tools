package main
// go写的zshrc笔记解析为markdown大纲
import (
	"fmt"
	"strings"
	"regexp"
	"io/ioutil"
)

func repeatString(input string, times int) string {
	var result strings.Builder

	for i := 0; i < times; i++ {
		result.WriteString(input)
	}

	return result.String()
}

func main() {
	str1, _ := ioutil.ReadFile("/Users/emacspy/.xonshrc")
	stri := string(str1)
	lines := strings.Split(stri, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "##") {
			outlineRe, _ := regexp.Compile(`=>|=》`)
			outlines := outlineRe.Split(strings.Replace(line, "##", "", 1), -1)
			for num, outline := range outlines {
				oline := strings.TrimSpace(outline)
				fmt.Println(repeatString("    ", num) + "- " + oline)
			}
		}
	}

}

