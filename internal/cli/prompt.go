package cli

import (
	"errors"
	"io"
	"strings"
)

func readPromptLine(input io.Reader) (string, error) {
	var value strings.Builder
	buffer := make([]byte, 1)
	for {
		count, err := input.Read(buffer)
		if count > 0 {
			switch buffer[0] {
			case '\n':
				return strings.TrimSuffix(value.String(), "\r"), nil
			case '\r':
				continue
			default:
				value.WriteByte(buffer[0])
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) && value.Len() > 0 {
				return value.String(), nil
			}
			return "", err
		}
	}
}
