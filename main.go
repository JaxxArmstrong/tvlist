package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"strconv"

	"github.com/Guitarbum722/align"
	"github.com/denisbrodbeck/striphtmltags"
	"github.com/docopt/docopt-go"
)

var searchCMD bool
var searchArgs string
var listCMD bool
var listModeID bool
var listModeIDArgs string
var listModeFile bool
var listModeFileArgs string
var lines []string
var url string
var list string
var mytvmaze1 = MyTVmaze{}
var mytvmaze2 = MyTVmazeDates{}
var mytvmaze3 = MyTVmazeEpisodes{}

var myFile = []string{"~/Documents/tvshows", "~/.tvshows", "~/tvshows"}

type MyTVmaze struct {
	ID           int      `json:"id"`
	URL          string   `json:"url"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Language     string   `json:"language"`
	Genres       []string `json:"genres"`
	Status       string   `json:"status"`
	Runtime      int      `json:"runtime"`
	Premiered    string   `json:"premiered"`
	OfficialSite string   `json:"officialSite"`
	Schedule     struct {
		Time string   `json:"time"`
		Days []string `json:"days"`
	} `json:"schedule"`
	Rating struct {
		Average float64 `json:"average"`
	} `json:"rating"`
	Weight     int         `json:"weight"`
	Network    interface{} `json:"network"`
	WebChannel struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Country struct {
			Name     string `json:"name"`
			Code     string `json:"code"`
			Timezone string `json:"timezone"`
		} `json:"country"`
	} `json:"webChannel"`
	Externals struct {
		Tvrage  interface{} `json:"tvrage"`
		Thetvdb int         `json:"thetvdb"`
		Imdb    string      `json:"imdb"`
	} `json:"externals"`
	Image struct {
		Medium   string `json:"medium"`
		Original string `json:"original"`
	} `json:"image"`
	Summary string `json:"summary"`
	Updated int    `json:"updated"`
	Links   struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Previousepisode struct {
			Href string `json:"href"`
		} `json:"previousepisode"`
		Nextepisode struct {
			Href string `json:"href"`
		} `json:"nextepisode"`
	} `json:"_links"`
}

type MyTVmazeDates struct {
	ID       int       `json:"id"`
	URL      string    `json:"url"`
	Name     string    `json:"name"`
	Season   int       `json:"season"`
	Number   int       `json:"number"`
	Airdate  string    `json:"airdate"`
	Airtime  string    `json:"airtime"`
	Airstamp time.Time `json:"airstamp"`
	Runtime  int       `json:"runtime"`
	Image    struct {
		Medium   string `json:"medium"`
		Original string `json:"original"`
	} `json:"image"`
	Summary interface{} `json:"summary"`
	Links   struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
}

type MyTVmazeEpisodes struct {
	ID       int       `json:"id"`
	URL      string    `json:"url"`
	Name     string    `json:"name"`
	Season   int       `json:"season"`
	Number   int       `json:"number"`
	Type     string    `json:"type"`
	Airdate  string    `json:"airdate"`
	Airtime  string    `json:"airtime"`
	Airstamp time.Time `json:"airstamp"`
	Runtime  int       `json:"runtime"`
	Image    struct {
		Medium   string `json:"medium"`
		Original string `json:"original"`
	} `json:"image"`
	Summary string `json:"summary"`
	Links   struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
}

func printOutput(showList []string) {

	var showListExploded = ""

	sort.Strings(showList)

	for _, line := range showList {
		showListExploded = showListExploded + line + "\n"
	}

	input := strings.NewReader(showListExploded)
	output := bytes.NewBufferString("")
	q := align.TextQualifier{On: true, Qualifier: "\""}
	aligner := align.NewAlign(input, output, "|", q)
	aligner.OutputSep("/")
	aligner.Align()

	fmt.Println(output)

	return
}

func doFileExist(contraband string) bool {
	if _, err := os.Stat(contraband); err == nil {
		// path/to/whatever does exists
		return true
	}
	// the path does not exist or some error occurred.
	return false
}

func getTVshows(url string) {

	tvmazeClient := http.Client{
		Timeout: time.Second * 3, // Maximum of 3 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:66.0) Gecko/20100101 Firefox/66.0")

	res, getErr := tvmazeClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	mytvmaze1 = MyTVmaze{}

	jsonErr := json.Unmarshal(body, &mytvmaze1)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return
}

func getTVshowsAirdate(url string) string {

	tvmazeClient2 := http.Client{
		Timeout: time.Second * 3, // Maximum of 3 secs
	}

	req2, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req2.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:66.0) Gecko/20100101 Firefox/66.0")

	res2, getErr := tvmazeClient2.Do(req2)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body2, readErr := ioutil.ReadAll(res2.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	mytvmaze2 = MyTVmazeDates{}

	jsonErr2 := json.Unmarshal(body2, &mytvmaze2)
	if jsonErr2 != nil {
		log.Fatal(jsonErr2)
	}

	return mytvmaze2.Airdate
}

func gatherPrevEpNum(url string) string {

	prevepnum := ""

	tvmazeClient3 := http.Client{
		Timeout: time.Second * 3, // Maximum of 3 secs
	}

	req3, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req3.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:66.0) Gecko/20100101 Firefox/66.0")

	res3, getErr := tvmazeClient3.Do(req3)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body3, readErr := ioutil.ReadAll(res3.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	mytvmaze3 = MyTVmazeEpisodes{}

	jsonErr3 := json.Unmarshal(body3, &mytvmaze3)
	if jsonErr3 != nil {
		log.Fatal(jsonErr3)
	}

	prevseason := strconv.FormatInt(int64(mytvmaze3.Season), 10)
	prevnumber := strconv.FormatInt(int64(mytvmaze3.Number), 10)

	prevepnum = "s" + prevseason + "e" + prevnumber

	return prevepnum
}

func gatherInfo(url string) (string, string, string) {
	nextep := ""
	prevep := ""
	prevepnum := ""

	getTVshows(url)

	if mytvmaze1.Links.Nextepisode.Href != "" {
		nextep = getTVshowsAirdate(mytvmaze1.Links.Nextepisode.Href)
	} else {
		nextep = "N/A"
	}

	if mytvmaze1.Links.Previousepisode.Href != "" {
		prevep = getTVshowsAirdate(mytvmaze1.Links.Previousepisode.Href)
		url = mytvmaze1.Links.Previousepisode.Href
		prevepnum = " ("+gatherPrevEpNum(url)+")"
	} else {
		prevep = "N/A"
	}

	return nextep, prevep, prevepnum
}

func main() {
	fmt.Println()

	usage := `tvlist v0.7.1

Usage:
	tvlist search (<name>)
	tvlist list file [<filename>]
	tvlist list id (<id>)
	tvlist -h | --help

Options:
	-h --help   Show this screen.
	id          ID mode.
	<id>        ID(s) of the TV shows to list.
	file        file mode.
	<filename>  File to parse for TV shows.`

	args, _ := docopt.ParseDoc(usage)

	searchCMD, _ = args.Bool("search")
	searchArgs, _ = args.String("<name>")
	listCMD, _ = args.Bool("list")
	listModeID, _ = args.Bool("id")
	listModeIDArgs, _ = args.String("<id>")
	listModeFile, _ = args.Bool("file")
	listModeFileArgs, _ = args.String("<filename>")

	nextep := ""
	prevep := ""
	prevepnum := ""

	if listCMD {
		if listModeFile {
			usr, _ := user.Current()
			dir := usr.HomeDir

			if listModeFileArgs == "" {
				for _, i := range myFile {

					// golang is annoying so we need to expand ~
					if strings.HasPrefix(i, "~/") {
						i = filepath.Join(dir, i[2:])
					}

					// Does a default file exist?
					foundYou := doFileExist(i)
					if foundYou == true {
						listModeFileArgs = i
						break
					} else {
						// No default file found so flip the user the bird.
						fmt.Println("No default file nor a file given at the command line, exiting.")
						os.Exit(1)
					}
				}
			}

			file, err := os.Open(listModeFileArgs)
			if err != nil {
				fmt.Println()
				log.Fatal(err)
			}

			fmt.Printf("Working")

			fscanner := bufio.NewScanner(file)

			for fscanner.Scan() {
				fmt.Printf(".")

				myString := fscanner.Text()

				s := strings.Split(myString, " ")[0]
				url = "http://api.tvmaze.com/shows/" + s

				nextep, prevep, prevepnum = gatherInfo(url)

				lines = append(lines, nextep+"|"+prevep+prevepnum+"|"+mytvmaze1.Name+"|"+mytvmaze1.Status)
			}

			defer file.Close()

		} else {
			s := strings.Split(listModeIDArgs, ",")

			fmt.Printf("Working")

			for i := range s {
				fmt.Printf(".")

				url = "http://api.tvmaze.com/shows/" + s[i]

				nextep, prevep, prevepnum = gatherInfo(url)

                                lines = append(lines, nextep+"|"+prevep+" ("+prevepnum+")|"+mytvmaze1.Name+"|"+mytvmaze1.Status)
			}
		}

		fmt.Printf("\n\n")

		printOutput(lines)
	}

	if searchCMD {
		url = "http://api.tvmaze.com/singlesearch/shows?q=" + searchArgs

		nextep, prevep, prevepnum = gatherInfo(url)

		// fmt.Println()
		fmt.Println("ID:", mytvmaze1.ID)
		fmt.Println("Name:", mytvmaze1.Name)
		fmt.Println("Status:", mytvmaze1.Status)
		fmt.Println("Rating:", mytvmaze1.Rating.Average)

		fmt.Printf("Genres: ")

		for i := range mytvmaze1.Genres {
			fmt.Printf(mytvmaze1.Genres[i])
			fmt.Printf(" ")
		}

		fmt.Println("\nLanguage:", mytvmaze1.Language)

		fmt.Println("Network:", mytvmaze1.WebChannel.Name)

		fmt.Println("Premiered:", mytvmaze1.Premiered)
		fmt.Println("URL:", mytvmaze1.URL)

		fmt.Println("---")

		mySummary := striphtmltags.StripTags(mytvmaze1.Summary)
		fmt.Println("Summary:", mySummary)

		fmt.Println("---")

		fmt.Println("Previous episode:", prevep)
		fmt.Println("Next episode:", nextep)
	}

	os.Exit(0)
}

