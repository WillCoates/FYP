package main

import (
	"testing"
)

func TestSimpleTopic(t *testing.T) {
	const topic = "foo/bar"
	const expectedRegex = "^foo\\.bar$"
	regex, err := TopicToRegex(topic)
	if err != nil {
		t.Fatal(err)
	}
	if regex.String() != expectedRegex {
		t.Error("Regex mismatch, expect", expectedRegex, "actual", regex.String())
	}
	if regex.MatchString("foo.bash") {
		t.Error("Regex matches invalid value")
	}
	if !regex.MatchString("foo.bar") {
		t.Error("Regex doesn't match valid value")
	}
}

func TestMultilevelWildcardTopic(t *testing.T) {
	const topic = "foo/#"
	const expectedRegex = "^foo(\\.?[a-zA-Z0-9]+)*$"
	regex, err := TopicToRegex(topic)
	if err != nil {
		t.Fatal(err)
	}
	if regex.String() != expectedRegex {
		t.Error("Regex mismatch, expect", expectedRegex, "actual", regex.String())
	}
	if regex.MatchString("fash") {
		t.Error("Regex matches invalid value")
	}
	if !regex.MatchString("foo") {
		t.Error("Regex doesn't match root topic")
	}
	if !regex.MatchString("foo.bar") {
		t.Error("Regex doesn't match single depth subtopic")
	}
	if !regex.MatchString("foo.bar.bash") {
		t.Error("Regex doesn't match multi depth subtopic")
	}
}
