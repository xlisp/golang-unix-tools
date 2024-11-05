package main

// Go script to parse zshrc notes into a markdown outline
import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

func repeatString(input string, times int) string {
	var result strings.Builder

	for i := 0; i < times; i++ {
		result.WriteString(input)
	}

	return result.String()
}

func main() {
	str1, _ := ioutil.ReadFile("/home/xlisp/.zshrc")
	stri := string(str1)
	lines := strings.Split(stri, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "##") {
			outlineRe, _ := regexp.Compile(`=>|=》`)
			outlines := outlineRe.Split(strings.Replace(line, "##", "", 1), -1)
			for num, outline := range outlines {
				oline := strings.TrimSpace(outline)
				if (num == 0) { 
					fmt.Println(repeatString("    ", num) + "    - ^ " + oline)
				} else {
					fmt.Println(repeatString("    ", num) + "    - " + oline)
				}
			}
		}
	}
}
