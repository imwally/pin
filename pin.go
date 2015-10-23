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
	options   = flag.NewFlagSet("", flag.ExitOnError)
	privFlag  = options.Bool("private", false, "private bookmark")
	readFlag  = options.Bool("readlater", false, "read later bookmark")
	extFlag   = options.String("text", "", "longer description of bookmark")
	tagFlag   = options.String("tag", "", "tags for bookmark")
	longFlag  = options.Bool("l", false, "display long format")
	titleFlag = options.String("title", "", "title of the bookmark")

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
	
	if !TokenIsSet() {
		return
	}
	
	var p pinboard.Post
	p.Token = token
	
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "No command is given.\n")
	}
	
	cmd := flag.Arg(0)
	
	if cmd == "add" {
		args := flag.Args()[2:]
		p.URL = flag.Args()[1]
		options.Parse(args)
		Add(p)
	}

	if cmd == "ls" {
		args := flag.Args()[1:]
		options.Parse(args)
		Show(p)
	}

	if cmd == "rm" {
		p.URL = flag.Args()[1]
		Delete(p)
	}
}
