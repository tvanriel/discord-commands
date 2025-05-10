package commands

import "strings"

// Strip the prefix from a command.
func StripPrefix(trigger string, command string) func(string) string {
	prefixLen := len(trigger) + len(command)

	return func(message string) string {
		return strings.TrimSpace(message[prefixLen:])
	}
}

// HasCommandPrefix determines whether or not a message has the given command prefix.
func HasCommandPrefix(trigger string, command string, message string) bool {
	if message == trigger+command {
		return true
	}

	return strings.HasPrefix(message, trigger+command+" ")
}

// SplitArgs splits the arguments of a command into a string slice.
func SplitArgs(s string) []string {
	return strings.Split(s, " ")
}
