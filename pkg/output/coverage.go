package output

import (
	"fmt"
	"io"
	"log"

	"mvdan.cc/sh/v3/syntax"
)

const HIGHLIGHT_COLOR = "\033[32m"
const RESET_COLOR = "\033[0m"

func OutputCoverage(w io.Writer, file io.Reader, covered []syntax.Pos, lens []uint) {
	p := uint(0)
	for i, pos := range covered {
		if p > pos.Offset() {
			continue
		}
		len := pos.Offset() - p
		prebuf := make([]byte, len)
		n, err := file.Read(prebuf)
		if err != nil && err != io.EOF {
			log.Fatalf("failed to read from file: %v", err)
		}
		fmt.Fprintf(w, "%s", prebuf[:n])
		buf := make([]byte, lens[i])
		n, err = file.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatalf("failed to read from file: %v", err)
		}
		fmt.Fprintf(w, "%s%s%s", HIGHLIGHT_COLOR, buf[:n], RESET_COLOR)
		p = pos.Offset() + lens[i]
	}
	for {
		buf := make([]byte, 1024)
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatalf("failed to read from file: %v", err)
		}
		if err == io.EOF {
			break
		}
		fmt.Fprintf(w, "%s", buf[:n])
	}
}
