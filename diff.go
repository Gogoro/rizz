package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	godiff "github.com/sourcegraph/go-diff/diff"
)

type FileDiff struct {
	Path     string
	OldPath  string
	IsNew    bool
	IsDelete bool
	IsBinary bool
	Added    int
	Removed  int
	Hunks    []Hunk
	Hash     string // sha256 of diff content, used to invalidate "viewed" mark when file changes
}

type Hunk struct {
	Header string
	Lines  []Line
}

type Line struct {
	Kind byte // ' ', '+', '-'
	Text string
}

func ParseDiff(raw []byte) ([]FileDiff, error) {
	parsed, err := godiff.ParseMultiFileDiff(raw)
	if err != nil {
		return nil, err
	}

	files := make([]FileDiff, 0, len(parsed))
	for _, fd := range parsed {
		file := FileDiff{
			Path:    cleanPath(fd.NewName),
			OldPath: cleanPath(fd.OrigName),
		}
		if fd.OrigName == "/dev/null" {
			file.IsNew = true
		}
		if fd.NewName == "/dev/null" {
			file.IsDelete = true
			file.Path = cleanPath(fd.OrigName)
		}

		var hashBuf bytes.Buffer
		hashBuf.WriteString(file.Path)
		for _, h := range fd.Hunks {
			hashBuf.Write(h.Body)
			hunk := Hunk{
				Header: fmt.Sprintf("@@ -%d,%d +%d,%d @@ %s",
					h.OrigStartLine, h.OrigLines, h.NewStartLine, h.NewLines, strings.TrimSpace(string(h.Section))),
			}
			for _, raw := range strings.Split(string(h.Body), "\n") {
				if raw == "" {
					continue
				}
				line := Line{Kind: raw[0]}
				if len(raw) > 1 {
					line.Text = raw[1:]
				}
				hunk.Lines = append(hunk.Lines, line)
				switch line.Kind {
				case '+':
					file.Added++
				case '-':
					file.Removed++
				}
			}
			file.Hunks = append(file.Hunks, hunk)
		}

		sum := sha256.Sum256(hashBuf.Bytes())
		file.Hash = hex.EncodeToString(sum[:])

		files = append(files, file)
	}
	return files, nil
}

func cleanPath(p string) string {
	if strings.HasPrefix(p, "a/") || strings.HasPrefix(p, "b/") {
		return p[2:]
	}
	return p
}
