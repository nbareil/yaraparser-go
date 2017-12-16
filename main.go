package yaraparser

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
)

// YaraPatternType represents the kind of pattern
type YaraPatternType int

const (
	RegularString YaraPatternType = iota
	Regexp
	HexFmt
)

type YaraModifier string

type YaraPattern struct {
	Type      YaraPatternType
	String    string
	Regexp    string
	Hex       string
	Modifiers []YaraModifier
}

// YaraRule is the parsed Yara rule
type YaraRule struct {
	Name    string
	Tags    []string
	Metas   map[string]string
	Strings map[string]YaraPattern
}

var eof = rune(0)

func main() {
	for _, fn := range os.Args[1:] {
		buf, err := ioutil.ReadFile(fn)
		if err != nil {
			log.Fatal(err)
		}

		parser := NewParser(bytes.NewBuffer(buf))
		parser.scan()
	}
}
