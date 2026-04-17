package main

import "github.com/sergi/go-diff/diffmatchpatch"

// WordDiffSegment is a chunk of one side of a paired diff line.
// Kind is 'e' (equal), '+' (inserted on the new side), or '-' (removed from old).
type WordDiffSegment struct {
	Text string
	Kind byte
}

// computeWordDiff returns the per-line segments for an old/new pair of lines.
// The old-side slice only contains 'e' and '-' segments; the new-side slice
// only contains 'e' and '+' segments. Semantic cleanup collapses trivial diffs.
func computeWordDiff(oldText, newText string) (oldSegs, newSegs []WordDiffSegment) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(oldText, newText, true)
	diffs = dmp.DiffCleanupSemantic(diffs)

	for _, d := range diffs {
		switch d.Type {
		case diffmatchpatch.DiffEqual:
			oldSegs = append(oldSegs, WordDiffSegment{Text: d.Text, Kind: 'e'})
			newSegs = append(newSegs, WordDiffSegment{Text: d.Text, Kind: 'e'})
		case diffmatchpatch.DiffDelete:
			oldSegs = append(oldSegs, WordDiffSegment{Text: d.Text, Kind: '-'})
		case diffmatchpatch.DiffInsert:
			newSegs = append(newSegs, WordDiffSegment{Text: d.Text, Kind: '+'})
		}
	}
	return oldSegs, newSegs
}
