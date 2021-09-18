package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/tkido/vkihou/config"
	"github.com/tkido/vkihou/myarr"
)

func main() {
	s := convert(config.Source)
	f, err := os.Create(config.Result)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(s)
}

var (
	reEmpty    = regexp.MustCompile(`^$`)
	reComment  = regexp.MustCompile(`^#`)
	reMessage  = regexp.MustCompile(`^\*.*$`)
	reHeadline = regexp.MustCompile(`^[★■●◆◇].*$`)
)

func headLine(line string) (string, MessageType) {
	line = line[1:]
	if line == "" {
		return `☆`, Narrate
	}
	ss := strings.Split(line, ".")
	switch len(ss) {
	case 0:
		fmt.Println("here!!!")
		return `☆`, Narrate
	case 1:
		return fmt.Sprintf("☆%s", ss[0]), Say
	default:
		tipe := Think
		return fmt.Sprintf("☆%s", ss[0]), tipe
	}
}

type Message struct {
	Name string
	Type MessageType
}

type MessageType int

const (
	Narrate MessageType = iota
	Say
	Think
)

func message(lines *myarr.MyArr, msg Message) *myarr.MyArr {
	// if lines.Size() >= 4 {
	// 	log.Fatalf("too many lines!!")
	// }

	var head, spacer string
	switch msg.Type {
	case Say:
		head = "「"
		spacer = "　"
	case Think:
		head = "（"
		spacer = "　"
	default:
		head = ""
		spacer = ""
	}
	buf := myarr.NewMyArr()
	buf.Push(head + lines.Pop())
	for lines.Size() > 0 {
		buf.Push(spacer + lines.Pop())
	}
	for i := buf.Size(); i < 4; i++ {
		buf.Push("")
	}
	return buf
}

func convert(path string) string {
	lines := myarr.ReadLines(path)
	buf := myarr.NewMyArr()
	msg := Message{}

	for lines.Size() > 0 {
		first := lines.First()
		switch {
		case first == "":
			lines.Pop()
		case reComment.MatchString(first):
			lines.Pop()
		case reHeadline.MatchString(first):
			buf.Push(lines.Pop())
			msg.Type = Narrate
		case reMessage.MatchString(first):
			script, tipe := headLine(lines.Pop())
			msg.Type = tipe
			buf.Push(script)
			buf.Push("")
			buf.Push("")
			buf.Push("")
		default:
			buf.Concat(message(lines.TakeBlockNot(reEmpty), msg))
		}
	}

	return buf.Join("\n")
}
