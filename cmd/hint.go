package cmd

import (
	"sort"

	"github.com/gsamokovarov/jump/cli"
	"github.com/gsamokovarov/jump/config"
	"github.com/gsamokovarov/jump/scoring"
)

func hintCmd(args cli.Args, conf config.Config) error {
	var hints scoring.Entries

	term := args.CommandName()
	smart := args.Has("--smart")

	entries, err := conf.ReadEntries()
	if err != nil {
		return err
	}

	if len(term) == 0 {
		// We usually keep them reversely sort to optimize the fuzzy search.
		sort.Sort(sort.Reverse(entries))

		hints = hintSmartSelect(entries, term, smart)
	} else {
		fuzzyEntries := scoring.NewFuzzyEntries(entries, term)

		hints = hintSmartSelect(fuzzyEntries.Entries, term, smart)
	}

	for _, entry := range hints {
		cli.Outf("%s\n", entry.Path)
	}

	return nil
}

func hintSmartSelect(entries scoring.Entries, term string, smart bool) scoring.Entries {
	if !smart {
		return entries
	}

	termLength := len(term)

	switch {
	case termLength == 0:
		return hintSliceEntries(entries, 5)
	case termLength < 7:
		return hintSliceEntries(entries, 1)
	case termLength >= 7 && termLength < 10:
		return hintSliceEntries(entries, 3)
	default:
		return hintSliceEntries(entries, 1)
	}
}

func hintSliceEntries(entries scoring.Entries, limit int) scoring.Entries {
	if limit < len(entries) {
		return entries[0:limit]
	}

	return entries
}

func init() {
	cli.RegisterCommand("hint", "Hints relevant paths for jumping.", hintCmd)
}
