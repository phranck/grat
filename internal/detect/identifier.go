package detect

import "strings"

// identifierCharacters are what a detected name may consist of.
//
// Letters, digits, underscore, hyphen and full stop, and nothing else. The set
// is deliberately smaller than what a filesystem or a Swift manifest allows,
// because the question here is not what a name may be called but what may be
// written into a command that later runs through /bin/sh.
const identifierCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-."

// safeIdentifier reports whether a name read out of a repository may be put into
// a command, and names the first character that says otherwise.
//
// Detection builds command lines from what a project's files say: an executable
// target out of Package.swift, a directory name below cmd, the module file of a
// Python application. Those files come with the repository, so their contents
// are the author's rather than the reader's, and a name carrying a semicolon or
// a backtick becomes a second command the moment the line is run.
//
// Quoting would be the other answer and is the wrong one here. The value ends up
// in a grat.config that a person reads and edits, so it has to look like the
// name it is; a name that cannot be written plainly is a name grat reports
// rather than one it escapes.
func safeIdentifier(name string) (offending rune, ok bool) {
	if name == "" {
		return 0, false
	}
	for _, character := range name {
		if !strings.ContainsRune(identifierCharacters, character) {
			return character, false
		}
	}
	return 0, true
}

// rejectedName is a name a detector read and will not put into a command,
// carried back to the caller so the reason names the name and the character
// rather than saying nothing was found.
type rejectedName struct {
	name      string
	offending rune
}

// unresolvedIdentifier is what a detector reports instead of a command, naming
// the file and the character so the reader can see what to change.
func unresolvedIdentifier(marker string, what string, name string, offending rune) Unresolved {
	return Unresolved{
		Marker: marker,
		Reason: what + " " + quoteForReason(name) + " carries " + describeRune(offending) +
			", which cannot go into a command, so grat proposes none rather than writing that line",
	}
}

// quoteForReason shows the name as it was read, so the reader can find it.
func quoteForReason(name string) string {
	return "\"" + name + "\""
}

// describeRune names a character the reader has to look for, since several of
// the ones that matter here are invisible in a message.
func describeRune(character rune) string {
	switch character {
	case ' ':
		return "a space"
	case '\t':
		return "a tab"
	case '\n':
		return "a line break"
	case '\r':
		return "a carriage return"
	}
	return "the character " + string(character)
}
