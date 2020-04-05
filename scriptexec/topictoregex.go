package main

import (
	"errors"
	"regexp"
	"strings"
)

var ErrInvalidTopic = errors.New("Invalid topic")

func TopicToRegex(topic string) (*regexp.Regexp, error) {
	var regex strings.Builder
	end := false
	start := true

	regex.WriteString("^")

	for _, topicLevel := range strings.Split(topic, "/") {
		if end {
			return nil, ErrInvalidTopic
		}

		if topicLevel == "#" {
			// Multi level wildcard
			regex.WriteString("(\\.?[a-zA-Z0-9]+)*")
			end = true
		} else {
			if !start {
				regex.WriteString("\\.")
			}
			if topicLevel == "+" {
				// Single level wildcard
				regex.WriteString("[a-zA-Z0-9]+")
			} else {
				regex.WriteString(topicLevel)
			}
		}

		if start {
			start = false
		}
	}

	regex.WriteString("$")

	return regexp.Compile(regex.String())
}
