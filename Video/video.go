package video

import (
	"encoding/json"
	"net/http"
	"fmt"
	"io"
	"math/rand"
)

type ArrayError struct{}

func (a *ArrayError) Error() string {
	return "no channel found"
}

type Video struct {
	ChannelName	string
	VideoTitle	string
	Thumbnail	string
	ReleaseDate	string //Maybe Date
	ChannelId	string
}


type Playlist struct {
	Items []struct {
		Snippet struct {
			Title string `json:"title"`
			Thumbnails struct {
				Default struct {
					URL string `json:"url"`
				}
			}
		}
	}
}

type CID struct {
	Items []struct {
		Id string `json:"id"`
	}
}

func (v *Video) GetChannelID(apik string) error {
	var channelsURL = "https://youtube.googleapis.com/youtube/v3/channels"
	client := &http.Client{}

	req, err := http.NewRequest("GET", channelsURL, nil)
	if err != nil {
		fmt.Println("error creating request:", err)
		return err
	}

	query := req.URL.Query()
	query.Add("part", "snippet")
	query.Add("forUsername", v.ChannelName)
	query.Add("key", apik)
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error sending request:", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failedf with status code %d\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil{
		fmt.Println("Error reading response body: ", err)
		return err
	}

	var s CID
	if err := json.Unmarshal(body, &s); err != nil {
		fmt.Println("Can not unmarshal JSON:", err)
		return err
	}

	if len(s.Items) == 0 {
		return &ArrayError{}
	}
	v.ChannelId = s.Items[0].Id
	return nil
}

func (v *Video) GetLatestVideo(apik string) {
	maxresult := "1"
	v.GetChannelID(apik)
	p := RetrieveAPIData(apik, v.ChannelId, maxresult)
	v.Thumbnail = p.Items[0].Snippet.Thumbnails.Default.URL
	v.VideoTitle = p.Items[0].Snippet.Title

	// return p.Items[0].Snippet.Title
}

func (v *Video) GetRandomVideo(apik string) {
	maxresult := "200"
	randomInt := rand.Intn(50)
	v.GetChannelID(apik)
	p := RetrieveAPIData(apik, v.ChannelId, maxresult) //youtube json with 200 videos
	v.Thumbnail = p.Items[randomInt].Snippet.Thumbnails.Default.URL
	v.VideoTitle = p.Items[randomInt].Snippet.Title
}

func RetrieveAPIData(apik string, id string, maxresult string) Playlist{
	var playlisturl = "https://youtube.googleapis.com/youtube/v3/playlistItems"
	client := &http.Client{}
	playlistid := "UU" + id[2:]

	req, err := http.NewRequest("GET", playlisturl, nil)
	if err != nil {
		fmt.Println("error creating request:", err)
	}

	query := req.URL.Query()
	query.Add("part", "snippet")
	query.Add("maxResults", maxresult)//"1")
	query.Add("playlistId", playlistid)
	query.Add("key", apik)
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error sending request:", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failedf with status code %d\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil{
		fmt.Println("Error reading response body: ", err)
	}

	var s Playlist
	if err := json.Unmarshal(body, &s); err != nil {
		fmt.Println("Can not unmarshal JSON:", err)
	}

	return s
}
