package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// move to go utils
type InputReader struct {
	reader *bufio.Reader
}

func (self *InputReader) Read(prompt string) (string, error) {
	if nil == self.reader {
		self.reader = bufio.NewReader(os.Stdin)
	}
	fmt.Printf("%v", prompt)
	value, err := self.reader.ReadString('\n')
	value = strings.Replace(value, "\n", "", -1)
	return value, err
}

//.end
