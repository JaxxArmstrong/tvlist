# tvlist
Lists current state of TV shows from either a predefined list (textfile), TV shows defined at the command line or search for a specific TV show. All modes leverages the TV Maze API.

UPDATE (2021-04-18):  Fixed the project to match the new modular Go structure, added season/number to the previous episode output to make it easier to track what shows you have seen or not and recompiled the Linux binary.

![](images/tvlist1.png?raw=true)
![](images/tvlist2.png?raw=true)
![](images/tvlist3.png?raw=true)

This personal project is published as-is and offers no support beyond this README.

It's written in [Go](https://golang.org/) and is part of my effort to learn the specific syntax by a learn-by-doing approach. The code is not elegant, but it's working as intended. I will polish it when time allows for it.

I've supplied a pre-compiled *Linux* standalone binary along with the source code for you who just want to clone the repo and have it working right off the bat, without installing *Go* as a language on your workstation.

## How it works

This utility fetches properties for predefined TV show from the JSON API backend at [TV Maze](http://www.tvmaze.com/). These are then shown in a short and aligned list in your terminal so you can get a quick overview of the TV shows you're following. It can also be used to search for a specific TV show (with more properties shown) instead of listing predefined TV shows.

If you want to use the utility to iterate through a textfile that contains TV show IDs you need to first create that file on your own and populate it with the appropriate unique TV Show ID found at TV Maze. See the below image for an example.

If filename is omitted on the command line for the list mode the utility will check for a default file in the following places:

- ~/Documents/tvshows
- ~/.tvshows
- ~/tvshows

![](images/tvlist4.png?raw=true)

Note that **only** the first column is extracted by the utility and the remaining text to the right of the TV Maze ID is just for comments.

The unique TV Maze ID can easily be fetched with the *search-mode* of this utility. Note that the TV Maze API tries to guess the correct TV show if you supply it with one or more keywords. If you need to use more than one keyword, please enclose it in quotes. See *usage* below.

The ID can also be gathered by using the API directly in a webbrowser. Note that sometimes more than one show has the same name so check the JSON result of the search to make sure you have the correct show. [Example search](http://api.tvmaze.com/singlesearch/shows?q=girls)


## Usage

```
tvlist list file [<textfile>]
tvlist list id 32819,12607,8898
tvlist search 'marvels agents'
tvlist search billions
```

I intentionally kept the commands & the arguments in a longer format for ease of understanding. If you don't agree with them you can either create an alias to your preferred shorter format or edit the source code and re-compile the executable.


If you have _Go_ installed as a language and your _$GOPATH_ set, you can just pull it with:
```
go get github.com/jaxxarmstrong/tvlist
```

---

## Credits

A big thank you to [gotbletu](https://www.youtube.com/channel/UCkf4VIqu3Acnfzuk3kRIFwA) for helping out with testing and recommending features.
