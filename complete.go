package main

import (
	"path/filepath"
	"runtime"
	"strings"
	"os"
	"unicode"
)

func fileComplete(query, ctx string, fields []string) ([]string, string) {
	var completions []string
	var extra string

	prefixes := []string{"./", "../", "/", "~/"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(query, prefix) {
			completions, extra, _ = matchPath(strings.Replace(query, "~", curuser.HomeDir, 1), query)
		}
	}

	if len(completions) == 0 && len(fields) > 1 {
		completions, extra, _ = matchPath("./" + query, query)
	}

	return completions, extra
}

func binaryComplete(query, ctx string, fields []string) ([]string, string) {
	var completions []string

	prefixes := []string{"./", "../", "/", "~/"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(query, prefix) {
			fileCompletions, prefix := fileComplete(query, ctx, fields)
			if len(fileCompletions) != 0 {
				for _, f := range fileCompletions {
					name := strings.Replace(query + f, "~", curuser.HomeDir, 1)
					if err := findExecutable(name, false, true); err != nil {
						continue
					}
					completions = append(completions, f)
				}
			}
			return completions, prefix
		}
	}

	// filter out executables, but in path
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		// print dir to stderr for debugging
		// search for an executable which matches our query string
		if matches, err := filepath.Glob(filepath.Join(dir, query + "*")); err == nil {
			// get basename from matches
			for _, match := range matches {
				// check if we have execute permissions for our match
				err := findExecutable(match, true, false)
				if err != nil {
					continue
				}
				// get basename from match
				name := filepath.Base(match)
				// add basename to completions
				completions = append(completions, name)
			}
		}
	}

	// add lua registered commands to completions
	for cmdName := range commands {
		if strings.HasPrefix(cmdName, query) {
			completions = append(completions, cmdName)
		}
	}

	completions = removeDupes(completions)

	return completions, query
}

func matchPath(path, pref string) ([]string, string, error) {
	var entries []string
	matches, err := filepath.Glob(desensitize(path) + "*")
	if err == nil {
		args := []string{
			"\"", "\\\"",
			"'", "\\'",
			"`", "\\`",
			" ", "\\ ",
			"(", "\\(",
			")", "\\)",
			"[", "\\[",
			"]", "\\]",
		}

		r := strings.NewReplacer(args...)
		for _, match := range matches {
			name := filepath.Base(match)
			matchFull, _ := filepath.Abs(match)
			if info, err := os.Stat(matchFull); err == nil && info.IsDir() {
				name = name + string(os.PathSeparator)
			}
			name = r.Replace(name)
			entries = append(entries, name)
		}
	}

	return entries, filepath.Base(pref), err
}

func desensitize(text string) string {
	if runtime.GOOS == "windows" {
		return text
	}

	p := strings.Builder{}

	for _, r := range text {
		if unicode.IsLetter(r) {
			p.WriteString("[" + string(unicode.ToLower(r)) + string(unicode.ToUpper(r)) + "]")
		} else {
			p.WriteString(string(r))
		}
	}

	return p.String()
}
