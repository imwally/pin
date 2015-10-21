package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/imwally/pinboard"
)

var (
	privFlag   = flag.Bool("private", false, "private bookmark")
	readFlag   = flag.Bool("readlater", false, "read later bookmark")
	extFlag    = flag.String("text", "", "longer description of bookmark")
	tagFlag    = flag.String("tag", "", "tags for bookmark")
	longFlag   = flag.Bool("l", false, "display long format")
	titleFlag  = flag.String("title", "", "title of the bookmark")

	token string
)

// Add checks flag values and encodes the GET URL for adding a bookmark.
func Add(p pinboard.Post) {

    p.Description = *titleFlag

    if *privFlag {
        p.Shared = "no"
    }

    if *readFlag {
        p.Toread = "yes"
    }

    if *extFlag != "" {
        p.Extended = *extFlag
    }

    if *tagFlag != "" {
        p.Tags = *tagFlag
    }

    p.Encode()
    err := p.Add()
    if err != nil {
        fmt.Println(err)
    }
}

// Delete will delete the URL value of the -d flag.
func Delete(p pinboard.Post) {
    p.Encode()
    err := p.Delete()
    if err != nil {
        fmt.Println(err)
    }
}

// Show will list the most recent bookmarks. The -show flag indicates how many
// bookmarks to show with a max of up to 100 bookmarks.
func Show(p pinboard.Post) {
    if *tagFlag != "" {
        p.Tag = *tagFlag
    }

    p.Encode()

    recent := p.ShowRecent()

    if *longFlag {
        for _, v := range recent.Posts {
            var shared string
            if v.Shared == "no" {
                shared = "*"
            }
            fmt.Println(shared + v.Description)
            fmt.Println(v.Href)
            fmt.Println(v.Tags, "\n")
        }
    } else {
        for _, v := range recent.Posts {
            fmt.Println(v.Href)
        }
    }
}

// TokenIsSet will check to make sure an authentication token is set before
// making an API calls.
func TokenIsSet() bool {
    if token == "" {
        return false
    }
    return true
}

func init() {
	u, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}

	content, err := ioutil.ReadFile(u.HomeDir + "/.pinboard")
	if err != nil {
        fmt.Println("No authorization token found. Please add your authorization token to ~/.pinboard")
	} 

	token = string(content)
}

func main() {
	
	if TokenIsSet() {
		var p pinboard.Post
		p.Token = token

		cmd := os.Args[1]
		
		flag.Parse()

		if cmd == "add" {
			p.URL = os.Args[2]
			Add(p)
		}

		if cmd == "rm"  {
			p.URL = os.Args[2]
			Delete(p)
		}

		if cmd == "ls" {
			Show(p)
		}
	}

}
