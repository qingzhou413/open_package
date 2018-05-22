package file

import (
	"io/ioutil"
	"strings"
)

//把文件指定行替换为新行
func Replace(fileName string, fromStr string, dstStr string) error {
	input, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if line == fromStr {
			lines[i] = dstStr
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(fileName, []byte(output), 0744)
	return err
}
