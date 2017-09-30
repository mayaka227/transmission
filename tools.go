package transmission

import (
	"github.com/swatkat/gotrntmetainfoparser"
	"encoding/base64"
	"io/ioutil"
	"fmt"
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

func (c *Client) AddTrackers(torrentfile string) (*Torrent, error) {
	name, trackers := getInfoFromTorrentFile(torrentfile)
	torrents, err := c.GetTorrents()
	if err2 != nil {
		return nil, err
	}
	for _, torrent := range torrents {
		if torrent.Name == name {
			for _, x := range trackers {
				y := []string{x}
				sArg := SetTorrentArg{
					TrackerAdd: y,
				}
				torrent.Set(sArg)
			}
			return torrent, nil
		}
	}
	return nil,fmt.Errorf("Torrent not exists")
}

func getInfoFromTorrentFile(filename string) (string, []string) {
	var a gotrntmetainfoparser.MetaInfo
	a.ReadTorrentMetaInfoFile(filename)
	var tracker []string
	for _, x := range a.AnnounceList {
		for _, y := range x {
			tracker = append(tracker, y)
		}
	}
	return a.Info.Name, tracker
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
