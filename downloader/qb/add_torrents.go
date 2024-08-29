package qb

import (
	"bytes"
	"fmt"
	"github.com/heibizi/go-media-kit/core/utils/netx"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type AddTorrentsInfo struct {
	Urls         string   `json:"urls"`
	Torrents     []string `json:"torrents"`
	SavePath     string   `json:"savePath"`
	Category     string   `json:"category"`
	Tags         string   `json:"tags"`
	SkipChecking bool     `json:"skipChecking"`
	Paused       bool     `json:"paused"`
	RootFolder   bool     `json:"rootFolder"`
	AutoTMM      bool     `json:"autoTMM"`
}

func (c *Client) AddTorrents(req any) error {
	info := req.(AddTorrentsInfo)
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	params := map[string]string{
		"urls":               info.Urls,
		"paused":             fmt.Sprintf("%t", info.Paused),
		"skip_checking":      fmt.Sprintf("%t", info.SkipChecking),
		"autoTMM":            fmt.Sprintf("%t", info.AutoTMM),
		"sequentialDownload": "false",
		"firstLastPiecePrio": "false",
		"contentLayout":      "Original",
		"stopCondition":      "None",
		"root_folder":        fmt.Sprintf("%t", info.RootFolder),
		"tags":               info.Tags,
		"category":           info.Category,
		"savepath":           info.SavePath,
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
	resp, err := netx.NewHttpx(netx.HttpRequestConfig{
		Ctx:         c.config.Ctx,
		Url:         c.config.Host + "/api/v2/torrents/add",
		Cookie:      c.ck,
		Referer:     c.config.Host,
		ContentType: writer.FormDataContentType(),
	}).Post(nil, payload)
	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("添加种子失败: %v", err)
	}
	return nil
}
