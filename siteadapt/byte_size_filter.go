package siteadapt

import (
	"fmt"
	"github.com/heibizi/go-media-kit/core/utils/stringx"
)

type byteSizeFilter struct {
	text   string
	filter *Filter
}

func (f byteSizeFilter) doFilter() (string, error) {
	return fmt.Sprintf("%d", stringx.ByteSize(f.text)), nil
}
