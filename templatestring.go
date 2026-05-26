package templatestring

import (
	"fmt"
	"regexp"
	"strings"
)

type tokenType int

const (
	literal tokenType = iota
	token   tokenType = iota
)

type segment struct {
	t tokenType
	v string
}
type templateString struct {
	segments []segment
}

var parser = regexp.MustCompile(`\$[(\{]\s*(\S+?)\s*[)\}]`)

func NewTemplateString(template string) *templateString {
	rv := templateString{}

	if len(template) == 0 {
		return &rv
	}

	rv.segments = make([]segment, 0, 1)

	if strings.Contains(template, "$") {
		matches := parser.FindAllStringSubmatchIndex(template, -1)
		i := 0
		for _, match := range matches {
			if i != match[0] {
				rv.segments = append(rv.segments, segment{
					t: literal,
					v: template[i:match[0]],
				})
			}

			rv.segments = append(rv.segments, segment{
				t: token,
				v: template[match[2]:match[3]],
			})

			i = match[1]
		}
		if i != len(template) {
			rv.segments = append(rv.segments, segment{
				t: literal,
				v: template[i:],
			})
		}
	} else {
		rv.segments = append(rv.segments, segment{
			t: literal,
			v: template,
		})
	}
	return &rv
}

func (t *templateString) Render(plugins ...Plugin) (string, error) {
	if len(t.segments) == 0 {
		return "", nil
	}

	if len(t.segments) == 1 && t.segments[0].t == literal {
		return t.segments[0].v, nil
	}

	rvlen := 0
	rv := make([]string, 0, len(t.segments))
	for _, s := range t.segments {
		if s.t == literal {
			rv = append(rv, s.v)
			rvlen += len(s.v)
		} else if s.t == token {
			var err error
			rvi := s.v
			isProcessed := false
			nProcessed := 0
			for _, plugin := range plugins {
				rvi, isProcessed, err = plugin.ProcessToken(rvi)
				if err != nil {
					return "", err
				}
				if isProcessed {
					nProcessed++
				}
			}

			if nProcessed == 0 {
				return "", fmt.Errorf("no plugins to process token [%s] specified", s.v)
			}

			rv = append(rv, rvi)
			rvlen += len(rvi)
		} else {
			// coverage-ignore
			panic("never")
		}
	}

	b := make([]byte, 0, rvlen)
	for _, rvi := range rv {
		b = append(b, rvi...)
	}
	return string(b), nil
}
