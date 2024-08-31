package siteadapt

import (
	"fmt"
	"github.com/heibizi/go-media-kit/core/utils/stringx"
)

type timestampFilter struct {
	text   string
	filter *Filter
}

func (f timestampFilter) doFilter() (string, error) {
	timeStamp := stringx.TimeStamp(f.text)
	if timeStamp == 0 {
		return "", nil
	}
	return fmt.Sprintf("%d", timeStamp), nil
}
