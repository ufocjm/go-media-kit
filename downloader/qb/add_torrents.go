package qb

import (
	"bytes"
	"fmt"
	"github.com/heibizi/go-media-kit/core/utils/netx"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type AddTorrentsInfo struct {
	Urls               []string `json:"urls"`
	Torrents           []string `json:"torrents"`
	SavePath           string   `json:"savePath"`
	Cookie             string   `json:"cookie"`
	Category           string   `json:"category"`
	Tags               []string `json:"tags"`
	SkipChecking       bool     `json:"skipChecking"`
	Paused             bool     `json:"paused"`
	RootFolder         bool     `json:"rootFolder"`
	Rename             string   `json:"rename"`
	UpLimit            int      `json:"upLimit"` // bytes/second
	DlLimit            int      `json:"dlLimit"` // bytes/second
	RatioLimit         float64  `json:"ratioLimit"`
	SeedingTimeLimit   int      `json:"seedingTimeLimit"` // minutes
	AutoTMM            bool     `json:"autoTMM"`
	SequentialDownload bool     `json:"sequentialDownload"`
	FirstLastPiecePrio bool     `json:"firstLastPiecePrio"`
}

func (c *Client) AddTorrents(req any) error {
	info := req.(AddTorrentsInfo)
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	params := map[string]string{
		"urls":               strings.Join(info.Urls, "\n"),
		"savepath":           info.SavePath,
		"cookie":             info.Cookie,
		"category":           info.Category,
		"tags":               strings.Join(info.Tags, ","),
		"skip_checking":      fmt.Sprintf("%t", info.SkipChecking),
		"paused":             fmt.Sprintf("%t", info.Paused),
		"root_folder":        fmt.Sprintf("%t", info.RootFolder),
		"rename":             info.Rename,
		"upLimit":            strconv.Itoa(info.UpLimit),
		"dlLimit":            strconv.Itoa(info.DlLimit),
		"ratioLimit":         fmt.Sprintf("%f", info.RatioLimit),
		"seedingTimeLimit":   strconv.Itoa(info.SeedingTimeLimit),
		"autoTMM":            fmt.Sprintf("%t", info.AutoTMM),
		"sequentialDownload": fmt.Sprintf("%t", info.SequentialDownload),
		"firstLastPiecePrio": fmt.Sprintf("%t", info.FirstLastPiecePrio),
	}
	for k, v := range params {
		err := writer.WriteField(k, v)
		if err != nil {
			return fmt.Errorf("write field: %v", err)
		}
	}
	if len(info.Torrents) > 0 {
		for _, torrent := range info.Torrents {
			file, err := os.Open(torrent)
			if err != nil {
				return fmt.Errorf("open torrent file: %v", err)
			}
			torrentReader, err := writer.CreateFormFile("torrents", filepath.Base(torrent))
			if err != nil {
				return fmt.Errorf("create form file: %v", err)
			}
			_, err = io.Copy(torrentReader, file)
			if err != nil {
				return fmt.Errorf("copy torrent file: %v", err)
			}
		}
	}
	resp, err := netx.NewHttpx(netx.HttpRequestParams{
		Ctx:         c.config.Ctx,
		Method:      http.MethodPost,
		Url:         c.config.Host + "/api/v2/torrents/add",
		Cookie:      c.ck,
		Referer:     c.config.Host,
		ContentType: writer.FormDataContentType(),
		Body:        payload,
	}).Request()
	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("添加种子失败: %v", err)
	}
	return nil
}
