package loader

import (
	"regexp"
)

var (
	tagRe *regexp.Regexp
)

func init() {
	tagRe = regexp.MustCompile(`build:"-"`)
}

func skipField(tag string) bool {
	/*
		Possible tags combination:

		build:"-"

	*/
	return tagRe.MatchString(tag)
}
