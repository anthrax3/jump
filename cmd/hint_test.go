package cmd

import (
	"os"
	p "path"
	"strings"
	"testing"

	"github.com/gsamokovarov/jump/cli"
	s "github.com/gsamokovarov/jump/scoring"
)

func Test_hintCmd(t *testing.T) {
	wc := p.Join(td, "web-console")
	web := p.Join(td, "/client/website")

	conf := &testConfig{
		Entries: s.Entries{
			s.Entry{wc, &s.Score{Weight: 100, Age: s.Now}},
			s.Entry{web, &s.Score{Weight: 90, Age: s.Now}},
		},
	}

	output := capture(&os.Stdout, func() {
		hintCmd(cli.Args{}, conf)
	})

	lines := strings.Fields(output)

	if lines[0] != wc {
		t.Fatalf("Expected first line to be %s, got %s", wc, lines[0])
	}

	if lines[1] != web {
		t.Fatalf("Expected first line to be %s, got %s", web, lines[1])
	}
}

func Test_hintCmd_smart(t *testing.T) {
	p1 := p.Join(td, "web-console")
	p2 := p.Join(td, "/client/website")
	p3 := p.Join(td, "web")

	conf := &testConfig{
		Entries: s.Entries{
			s.Entry{p1, &s.Score{Weight: 100, Age: s.Now}},
			s.Entry{p2, &s.Score{Weight: 90, Age: s.Now}},
			s.Entry{p3, &s.Score{Weight: 80, Age: s.Now}},
		},
	}

	lines := strings.Fields(capture(&os.Stdout, func() {
		hintCmd(cli.Args{"wc", "--smart"}, conf)
	}))

	if len(lines) != 1 {
		t.Fatalf("Expected to get exactly one line, got %v", lines)
	}

	if lines[0] != p1 {
		t.Fatalf("Expected line to be %s, got %s", p1, lines[0])
	}

	// If you write more than 3 chars, maybe you need more options.
	lines = strings.Fields(capture(&os.Stdout, func() {
		hintCmd(cli.Args{"webs", "--smart"}, conf)
	}))

	if len(lines) != 3 {
		t.Fatalf("Expected to get exactly 3 lines, got %v", lines)
	}

	// If you wrote more than 7 chars, well, we tried.
	lines = strings.Fields(capture(&os.Stdout, func() {
		hintCmd(cli.Args{"web-conso", "--smart"}, conf)
	}))

	if len(lines) != 1 {
		t.Fatalf("Expected to get exactly one line, got %v", lines)
	}

	if lines[0] != p1 {
		t.Fatalf("Expected line to be %s, got %s", p1, lines[0])
	}

}