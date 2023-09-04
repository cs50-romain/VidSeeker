package main

import (
	"fmt"
	//"sync"
	"html/template"
	"net/http"

	"test/youtubecli/YT/Video"
	"test/youtubecli/YT/db"
)

type Video struct {
	Youtuber  string
	Thumbnail string
}

type User struct {
	Username	string
	Password	string
}

var youtubers = []string{}

var Reset = "\033[0m"
var Yellow = "\033[33m"
var Red = "\033[31m"
var White = "\033[97m"

const apikey = "AIzaSyCBXw_TDGnsxIJMJiuT6itozH6oYGwd-GI"
var videos = []Video{}
var userid int

// REMOVE YOUTUBER
func ViewOption4(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/index.tmpl")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	video := Video{
		Youtuber: r.FormValue("remove-youtuber-name"),
	}

	for idx, value := range videos {
		if video.Youtuber == value.Youtuber {
			fmt.Println("Removing: ", video.Youtuber)
			videos = append(videos[:idx], videos[idx+1:]...)
			break
		}
	}

	// REMOVE FROM DATABASE
	db.RemoveYoutuber(video.Youtuber, userid)

	data := struct {
		Video []Video
	}{
		Video: videos,
	}

	err = t.Execute(w, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(videos)
	fmt.Println("Received Option 4")
}

// ADD YOUTUBER
func ViewOption3(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/index.tmpl")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	youtuber := r.FormValue("new-youtuber-name")
	fmt.Println(youtuber)

	// ADD TO DATABASE
	v := new(yvideo.Video)
	v.ChannelName = (youtuber)
	err = v.GetChannelID(apikey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(v)
	playlistid := "UU" + v.ChannelId[2:]
	db.AddYoutuber(v.ChannelName, v.ChannelId, playlistid, userid)

	youtubers = append(youtubers, youtuber)
	// WRITE DOWN CODE TO UPDATE VIDEOS ARRAY TO SHOW MOST RECENT VIDEO OF ALL YOUTUBERS, INCLUDING ONE JSUT ADDED
	//videos = append(videos, video)

	data := struct {
		Video []Video
	}{
		Video: videos,
	}

	err = t.Execute(w, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/option1", http.StatusSeeOther)
	fmt.Println("Received Option 3")
}

// ONE YOUTUBER"S VIDEO
func ViewOption2(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/index.tmpl")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	youtuber := r.FormValue("youtuber-name")
	video := GetMostRecentVideo(youtuber)

	singleVideo := []Video{}
	singleVideo = append(singleVideo, Video{
		Youtuber: video.ChannelName,
		Thumbnail: video.Thumbnail,
	})

	data := struct {
		Video []Video
	}{
		Video: singleVideo,
	}

	err = t.Execute(w, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Received option 2")
}

// ALL MOST RECENT VIDEOS
func ViewOption1(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/index.tmpl")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tempvideos = []Video{}

	for _, youtuber := range youtubers{
		video := GetMostRecentVideo(youtuber)
		tempvideos = append(tempvideos, Video{
			Youtuber: video.ChannelName,
			Thumbnail: video.Thumbnail,
		})
	}

	data := struct {
		Video []Video
	}{
		Video: tempvideos,
	}

	err = t.Execute(w, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Received option1")
}

func viewLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, "./static/login.html")
	}

	user := User {
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	if loginCheck(user) {
		InitArray()
		fmt.Println("init: ", youtubers)
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	} else {
		fmt.Println("Invalid credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func viewIndex(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w, r, "./static/index.html")
	t, err := template.ParseFiles("./static/index.tmpl")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Video []Video
	}{
		Video: videos,
	}

	err = t.Execute(w, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func loginCheck(user User) bool{
	idx := 0
	for {
		dbuser, dbpass, user_id, err := db.UserById(idx)
		if err != nil {
			fmt.Println(err)
		}
		if user.Username == dbuser && user.Password == dbpass {
			userid = user_id
			return true
		}
		idx++
	}

	return false
}

func GetRandomVideos(youtuber string) yvideo.Video {
	v := new(yvideo.Video)
	v.ChannelName = youtuber
	v.GetRandomVideo(apikey)

	return *v
}

func GetMostRecentVideo(youtuber string) yvideo.Video {
	v := new(yvideo.Video)
	v.ChannelName = youtuber
	v.GetLatestVideo(apikey)

	return *v
}

func InitArray(){
	idx := 4 
	for idx < 15{
		name, err := db.YoutuberById(idx, userid)
		if err != nil {
			fmt.Println(err)
		} else {
			youtubers = append(youtubers, name)
		}
		idx++
	}
}

func main() {
	dberr := db.Connect()
	if dberr != nil {
		fmt.Println(dberr)
	}
	defer db.Close()

	http.HandleFunc("/", viewLogin)
	http.HandleFunc("/index", viewIndex)
	http.HandleFunc("/option1", ViewOption1)
	http.HandleFunc("/option2", ViewOption2)
	http.HandleFunc("/option3", ViewOption3)
	http.HandleFunc("/option4", ViewOption4)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Listening...")
}
