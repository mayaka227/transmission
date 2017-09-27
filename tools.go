package transmission

import (
	"encoding/base64"
	"io/ioutil"
)

func (c *Client) AddTorrentFile(torrentfile, torrentpath string) (*Torrent, error) {
	fileData, err := ioutil.ReadFile(torrentfile)
	if err != nil {
		return nil, err
	}
	meta := base64.StdEncoding.EncodeToString(fileData)

	addt := AddTorrentArg{
		DownloadDir: torrentpath,
		Metainfo:    meta,
	}

	torrent, err := c.AddTorrent(addt)

	return torrent, err
}

func (c *Client) RemoveTorrent(torrent *Torrent, removeData bool) error {
	ids := []int{torrent.ID}
	type arg struct {
		Ids             []int `json:"ids,string"`
		DeleteLocalData bool  `json:"delete-local-data,omitempty"`
	}

	tReq := &Request{
		Arguments: arg{
			Ids:             ids,
			DeleteLocalData: removeData,
		},
		Method: "torrent-remove",
	}
	r := &Response{}
	err := c.request(tReq, r)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) StartTorrent(torrents []*Torrent) error {
	return c.queueAction("torrent-start", torrents)
}
