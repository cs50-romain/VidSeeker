package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"test/youtubecli/YT/Video"
	"test/youtubecli/YT/db"
)

var Reset = "\033[0m"
var Yellow = "\033[33m"
var Red = "\033[31m"
var White = "\033[97m"

var apikey = "AIzaSyCBXw_TDGnsxIJMJiuT6itozH6oYGwd-GI"
var youtubers = []string{}

func displayThumbnail(url string) {
	cmdCurl := exec.Command("curl", "-s", url)//s.Items[0].Snippet.Thumbnails.Default.URL)
	cmdImgCat := exec.Command("imgcat")

	cmdImgCat.Stdin, _ = cmdCurl.StdoutPipe()
	cmdImgCat.Stdout = os.Stdout

	if err := cmdCurl.Start(); err != nil {
		fmt.Println(err)
	}

	if err := cmdImgCat.Start(); err != nil {
		fmt.Println(err)
	}

	if err := cmdCurl.Wait(); err != nil {
		fmt.Println(err)
	}
	
	if err := cmdImgCat.Wait(); err != nil {
		fmt.Println("error:",err)
	}
}

func parseCmd(input string) (string, error) {
	input = strings.TrimSuffix(input, "\n")

	args := strings.Split(input, " ")

	return args[0], nil 
}

func GetRandomVideos(youtubers []string, apikey string) {
	for _, youtuber := range youtubers{
		v := new(video.Video)
		v.ChannelName = youtuber
		v.GetRandomVideo(apikey)
		fmt.Println(Red + youtuber + Yellow + ":\t" + Red + v.VideoTitle + Reset)
		displayThumbnail(v.Thumbnail)
	}
}

func GetMostRecentVideos(youtubers []string, apikey string) {
	for _, youtuber := range youtubers{
		v := new(video.Video)
		v.ChannelName = youtuber
		title := v.GetLatestVideo(apikey)
		fmt.Println(Red + youtuber + Yellow + ":\t" + Red + title + Reset)
		displayThumbnail(v.Thumbnail)
	}
}

func RemoveYoutuber(youtubers []string, youtuber string) []string {
	youtuber = strings.TrimSuffix(youtuber, "\n")
	for idx, val := range youtubers {
		if val == youtuber {
			youtubers = append(youtubers[:idx], youtubers[idx+1:]...)
		}
	}

	db.RemoveYoutuber(youtuber)

	return youtubers
}

// THE WHOLE VIDEO PACKAGE WILL NEED REHTINKING
func AddYoutuber(input string, apik string){
	input = strings.TrimSuffix(input, "\n")
	v := new(video.Video)
	v.ChannelName = input
	err := v.GetChannelID(apik)
	if err != nil {
		fmt.Println(err)
		return
	}
	playlistid := "UU" + v.ChannelId[2:]

	db.AddYoutuber(v.ChannelName, v.ChannelId, playlistid)
}

func ListOptions() {
	fmt.Println(Yellow + "Options:")
	fmt.Println("\t1: Get most recent video from chosen youtuber")
	fmt.Println("\t2: Get most recent video from each favorited youtuber")
	fmt.Println("\t3: Add youtuber to favorites")
	fmt.Println("\t4: Remove youtuber from favorites")
	fmt.Println("\t5: Get random video from chosen youtuber")
	fmt.Println("\t6: List Options" + Reset)
}

func parseOptions(option string) {
	if option == "exit" {
		fmt.Println(White + "[!] Quitting..." + Reset)
		os.Exit(0)
	} else if option == "1" {
		fmt.Print(Yellow + "Youtuber: " + Reset)
		yreader := bufio.NewReader(os.Stdin)
		uinput, _ := yreader.ReadString('\n')
		var arryoutuber = []string{uinput}
		GetMostRecentVideos(arryoutuber, apikey)
	} else if option == "2" {
		GetMostRecentVideos(youtubers, apikey)
	}else if option == "3" {
		fmt.Print(Yellow + "Youtuber: " + Reset)
		yreader := bufio.NewReader(os.Stdin)
		uinput,_ := yreader.ReadString('\n')
			//NEEDS RETHINKING - I don't want to pass the apikey everywhere
		AddYoutuber(uinput, apikey)
		youtubers = append(youtubers, uinput)
	} else if option == "4" {
		fmt.Print(Yellow + "Youtuber: " + Reset)
		yreader := bufio.NewReader(os.Stdin)
		rminput,_ := yreader.ReadString('\n')
		youtubers = RemoveYoutuber(youtubers, rminput)
	} else if option == "5" {
		var arryoutuber = []string{}
		for {
			fmt.Print(Yellow + "Youtuber (q to stop): " + Reset)
			yreader := bufio.NewReader(os.Stdin)
			uinput, _ := yreader.ReadString('\n')
			if uinput == "q\n" {
				break
			}else {
				arryoutuber = append(arryoutuber, uinput)
			}
		}
		GetRandomVideos(arryoutuber, apikey)
	} else if option == "6" {
		ListOptions()
	}
}

func main() {
	dberr := db.Connect()
	if dberr != nil {
		fmt.Println(dberr)
	}
	defer db.Close()

	// Initialize youtuber array from youtubers in database
	// Index starts at 4 because id database table starts at 4
	idx := 4 
	for {
		name, err := db.YoutuberById(idx)
		if err != nil {
			fmt.Println(err)
			break
		}
		youtubers = append(youtubers, name)
		idx++
	}

	if len(youtubers) != 0 {
		fmt.Println(Yellow + "Current favorited youtubers" + Reset)
		for _, youtuber := range youtubers {
			fmt.Print(Red + youtuber + Yellow + " | " + Reset)
		}
	} else {
		fmt.Println(Yellow + "No favorite youtuber, press 3 to enter one" + Reset)
	}
	fmt.Println()

	// List options and START SHELL
	ListOptions()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		option, err := parseCmd(input)
		parseOptions(option)
	} // END OF SHELL INPUT
}
