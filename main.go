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
	Title     string
	VideoURL  string
}

type User struct {
	Username  string
	Password  string
}

var youtubers = []string{}

const apikey = "AIzaSyCBXw_TDGnsxIJMJiuT6itozH6oYGwd-GI"
var videos = []Video{}
var userid int

// ONE YOUTUBER'S RANDOM VIDEO
func ViewOption6(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost{
		t, err := template.ParseFiles("./static/index.tmpl")
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		youtuber := r.FormValue("random-youtuber-name")
		if youtuber == "" {
			http.Redirect(w, r, "/index", http.StatusSeeOther)
			return
		}

		video, err := GetRandomVideo(youtuber)

		singleVideo := []Video{}

		if err != nil {
			http.Redirect(w, r, "/index", http.StatusSeeOther)
			return
		} else {
			singleVideo = append(singleVideo, Video{
				Youtuber: video.ChannelName,
				Thumbnail: video.Thumbnail,
				Title: video.VideoTitle,
				VideoURL: video.VideoURL,
			})
		}

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
		fmt.Println("Received option 6")
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// ALL RANDOM VIDEOS
func ViewOption5(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		t, err := template.ParseFiles("./static/index.tmpl")
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, youtuber := range youtubers{
			video, err := GetRandomVideo(youtuber)
			if err != nil {
				break
			}
			videos = append(videos, Video{
				Youtuber: video.ChannelName,
				Thumbnail: video.Thumbnail,
				Title: video.VideoTitle,
				VideoURL: video.VideoURL,
			})
		}

		data := struct {
			Video []Video
		}{
			Video: videos,
		}

		err = t.Execute(w, data)

		videos = []Video{}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Received option5")
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// REMOVE YOUTUBER
func ViewOption4(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		t, err := template.ParseFiles("./static/index.tmpl")
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		video := Video{
			Youtuber: r.FormValue("remove-youtuber-name"),
		}

		// REMOVE FROM DATABASE
		db.RemoveYoutuber(video.Youtuber, userid)
		InitArray()

		for _, youtuber := range youtubers{
			video, err := GetMostRecentVideo(youtuber)
			if err != nil {
				break
			}
			videos = append(videos, Video{
				Youtuber: video.ChannelName,
				Thumbnail: video.Thumbnail,
				Title: video.VideoTitle,
				VideoURL: video.VideoURL,
			})
		}

		data := struct {
			Video []Video
		}{
			Video: videos,
		}

		err = t.Execute(w, data)

		videos = []Video{}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println(videos)
		fmt.Println("Received Option 4")
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// ADD YOUTUBER
func ViewOption3(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
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

		InitArray()

		for _, youtuber := range youtubers{
			video, err := GetMostRecentVideo(youtuber)
			if err != nil {
				break
			}
			videos = append(videos, Video{
				Youtuber: video.ChannelName,
				Thumbnail: video.Thumbnail,
				Title: video.VideoTitle,
				VideoURL: video.VideoURL,
			})
		}

		data := struct {
			Video []Video
		}{
			Video: videos,
		}

		err = t.Execute(w, data)

		videos = []Video{}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/option1", http.StatusSeeOther)
		fmt.Println("Received Option 3")
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// ONE YOUTUBER"S VIDEO
func ViewOption2(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		t, err := template.ParseFiles("./static/index.tmpl")
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		youtuber := r.FormValue("youtuber-name")
		if youtuber == "" {
			http.Redirect(w, r, "/index", http.StatusSeeOther)
			return
		}
		
		singleVideo := []Video{}
		video, err := GetMostRecentVideo(youtuber)
		if err != nil {
			http.Redirect(w, r, "/index", http.StatusSeeOther)
		} else {
			singleVideo = append(singleVideo, Video{
				Youtuber: video.ChannelName,
				Thumbnail: video.Thumbnail,
				Title: video.VideoTitle,
				VideoURL: video.VideoURL,
			})
		}
		
		data := struct {
			Video []Video
		}{
			Video: singleVideo,
		}

		err = t.Execute(w, data)

		videos = []Video{}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Received option 2")
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// ALL MOST RECENT VIDEOS
func ViewOption1(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		t, err := template.ParseFiles("./static/index.tmpl")
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, youtuber := range youtubers{
			video, err := GetMostRecentVideo(youtuber)
			if err != nil {
				break
			}
			videos = append(videos, Video{
				Youtuber: video.ChannelName,
				Thumbnail: video.Thumbnail,
				Title: video.VideoTitle,
				VideoURL: video.VideoURL,
			})
		}

		data := struct {
			Video []Video
		}{
			Video: videos,
		}

		err = t.Execute(w, data)

		videos = []Video{}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Received option1")
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func viewSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, "./static/login.html")
	}

	fmt.Println("VIEWING SIGNUP NOW")

	user := User {
		Username: r.FormValue("new-username"),
		Password: r.FormValue("new-password"),
	}

	db.AddUser(user.Username, user.Password)

	if loginCheck(user) {
		fmt.Println("Successfull")
		InitArray()
		fmt.Println("init: ", youtubers)
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	} else {
		fmt.Println("Invalid credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func viewLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		youtubers = []string{}
		videos = []Video{}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
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
	dbuser, dbpass, user_id, err := db.UserById(user.Username, user.Password)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if user.Username == dbuser && user.Password == dbpass {
		userid = user_id
		return true
	}

	return false
}

func GetRandomVideo(youtuber string) (yvideo.Video, error) {
	v := new(yvideo.Video)
	v.ChannelName = youtuber
	err := v.GetRandomVideo(apikey)
	if err != nil {
		return *v, nil
	}

	return *v, nil
}

func GetMostRecentVideo(youtuber string) (yvideo.Video, error) {
	v := new(yvideo.Video)
	v.ChannelName = youtuber
	err := v.GetLatestVideo(apikey)
	if err != nil {
		return *v, err
	}

	return *v, nil
}

func InitArray(){
	if len(youtubers) > 0 {
		youtubers = []string{}
	}
	idx := 4 
	for idx < 100{
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

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", viewLogin)
	http.HandleFunc("/signup", viewSignup)
	http.HandleFunc("/index", viewIndex)
	http.HandleFunc("/logout", viewLogout)
	http.HandleFunc("/option1", ViewOption1)
	http.HandleFunc("/option2", ViewOption2)
	http.HandleFunc("/option3", ViewOption3)
	http.HandleFunc("/option4", ViewOption4)
	http.HandleFunc("/option5", ViewOption5)
	http.HandleFunc("/option6", ViewOption6)
	http.ListenAndServe("192.168.4.87:8080", nil)
	fmt.Println("Listening...")
}
