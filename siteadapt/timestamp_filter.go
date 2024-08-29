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
	return fmt.Sprintf("%d", stringx.TimeStamp(f.text)), nil
}
