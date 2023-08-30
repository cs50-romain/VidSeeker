package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"test/youtubecli/YT/Video"
	"test/youtubecli/YT/db"
)

var Reset = "\033[0m"
var Yellow = "\033[33m"
var Red = "\033[31m"
var White = "\033[97m"

const apikey = "AIzaSyCBXw_TDGnsxIJMJiuT6itozH6oYGwd-GI"
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

func GetRandomVideos(youtuber string, ch chan<- video.Video, wg *sync.WaitGroup) {
	v := new(video.Video)
	v.ChannelName = youtuber
	v.GetRandomVideo(apikey)
	defer wg.Done()

	ch <- *v
}

func GetMostRecentVideos(youtuber string, ch chan<- video.Video, wg *sync.WaitGroup) {
	v := new(video.Video)
	v.ChannelName = youtuber
	v.GetLatestVideo(apikey)
	defer wg.Done()
	//fmt.Println(Red + youtuber + Yellow + ":\t" + Red + title + Reset)
	//displayThumbnail(v.Thumbnail)
	ch <- *v
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
	fmt.Println("\t6: Get random video from each favorited youtuber")
	fmt.Println("\t7: List Options" + Reset)
}

func parseOptions(option string) {
	ch := make(chan video.Video)
	//video := new(video.Video)
	var wg sync.WaitGroup

	if option == "exit" {
		fmt.Println(White + "[!] Quitting..." + Reset)
		os.Exit(0)
	} else if option == "1" {
		fmt.Print(Yellow + "Youtuber: " + Reset)
		yreader := bufio.NewReader(os.Stdin)
		uinput, _ := yreader.ReadString('\n')
		wg.Add(1)
		go GetMostRecentVideos(uinput, ch, &wg)

		go func() {
			wg.Wait()
			close(ch)
		}()

		for video := range ch {
			fmt.Println(Red + video.ChannelName + Yellow + ":\t" + Red + video.VideoTitle + Reset)
			displayThumbnail(video.Thumbnail)
		}

	} else if option == "2" {
		for _, youtuber := range youtubers {
			wg.Add(1)
			go GetMostRecentVideos(youtuber, ch, &wg) // Hangs the shell until user presses enter
		}

		go func() {
			wg.Wait()
			close(ch)
		}()

		for video := range ch {
			fmt.Println(Red + video.ChannelName + Yellow + ":\t" + Red + video.VideoTitle + Reset)
			displayThumbnail(video.Thumbnail)
		}
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
		fmt.Print(Yellow + "Youtuber: " + Reset)
		yreader := bufio.NewReader(os.Stdin)
		uinput, _ := yreader.ReadString('\n')
		
		wg.Add(1)
		go GetRandomVideos(uinput, ch, &wg) // Hangs the shell until user presses enter

		go func() {
			wg.Wait()
			close(ch)
		}()

		for video := range ch {
			fmt.Println(Red + video.ChannelName + Yellow + ":\t" + Red + video.VideoTitle + Reset)
			displayThumbnail(video.Thumbnail)
		}

		// OLD BUT I LIKED AND COULD STILL BE USED FOR ANOTHER OPTION
		/*var arryoutuber = []string{}
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
		GetRandomVideos(arryoutuber)
		*/
	} else if option == "6" {
		for _, youtuber := range youtubers {
			wg.Add(1)
			go GetRandomVideos(youtuber, ch, &wg) // Hangs the shell until user presses enter
		}

		go func() {
			wg.Wait()
			close(ch)
		}()

		for video := range ch {
			fmt.Println(Red + video.ChannelName + Yellow + ":\t" + Red + video.VideoTitle + Reset)
			displayThumbnail(video.Thumbnail)
		}
	} else if option == "7" {
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
