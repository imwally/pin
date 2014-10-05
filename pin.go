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
	addFlag    = flag.String("a", "", "url of bookmark to add")
	delFlag    = flag.String("d", "", "url of bookmark to delete")
	privFlag   = flag.Bool("p", false, "private bookmark")
	readFlag   = flag.Bool("r", false, "read later bookmark")
	extFlag    = flag.String("e", "", "longer description of bookmark")
	tagFlag    = flag.String("tag", "", "tags for bookmark")
	longFlag   = flag.Bool("l", false, "display long format")
	titleFlag  = flag.String("title", "", "title of the bookmark")
	showFlag   = flag.Int("show", 0, "show the most recent bookmarks")

	token string
)

func Add(p pinboard.Post) {
    p.Url = *addFlag
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

func Delete(p pinboard.Post) {
    p.Url = *delFlag
    p.Encode()
    err := p.Delete()
    if err != nil {
        fmt.Println(err)
    }
}

func Show(p pinboard.Post) {
    if *tagFlag != "" {
        p.Tag = *tagFlag
    }

    p.Count = *showFlag
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

        flag.Parse()

        if flag.NFlag() < 1 {
            fmt.Fprintf(os.Stderr, "No command given.\n")
            flag.Usage()
            return
        }

        if *addFlag != "" {
            Add(p)
        }

        if *delFlag != "" {
            Delete(p)
        }

        if *showFlag > 0 {
            Show(p)
        }
    }

}
