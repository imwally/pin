package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"

	"github.com/imwally/pinboard"
	"golang.org/x/net/html"
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

var usage = `Usage: pin
  pin rm  URL
  pin add URL [-title title] [OPTIONS]
  pin ls [-l] [-tag tags]

Options:
  -tag        space delimited tags 
  -private    mark bookmark as private
  -readlater  mark bookmark as read later
  -text       longer description of bookmark
  -l          long format for ls
`

// Number of bookmarks to display.
const COUNT int = 50

func HTMLTitle(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "title" {
		for _, a := range n.Attr {
			return a.Val
		}
	}

	return ""
}

// Add checks flag values and encodes the GET URL for adding a bookmark.
func Add(p pinboard.Post) {

	// Make sure a URL is specified. add, the sub command is the
	// first argument. The second argument is the URL being added.
	if flag.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "Not enough arguments.\n")
		return
	}

	// Parse flags after the URL.
	args := flag.Args()[2:]
	p.URL = flag.Args()[1]
	options.Parse(args)

	if *titleFlag != "" {
		p.Description = *titleFlag
	} else {
		// Grab page title
		resp, err := http.Get(p.URL)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()

		doc, err := html.Parse(resp.Body)
		if err != nil {
			fmt.Println(err)
		}

		p.Description = HTMLTitle(doc)
	}

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

// Delete will delete the URL specified.
func Delete(p pinboard.Post) {

	// Make sure a URL is specified. The URL being removed is the
	// second argument to the pin program, rm being the first.
	if flag.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "Not enough arguments.\n")
		return
	}

	p.URL = flag.Args()[1]
	p.Encode()

	err := p.Delete()
	if err != nil {
		fmt.Println(err)
	}
}

// Show will list the most recent bookmarks. The -tag flag can be used
// to filter results.
func Show(p pinboard.Post) {

	args := flag.Args()[1:]
	options.Parse(args)

	if *tagFlag != "" {
		p.Tag = *tagFlag
	}

	p.Count = COUNT
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

// runCmd takes a command string, initialises a new pinboard post and
// runs the command.
func runCmd(cmd string) {
	var p pinboard.Post
	p.Token = token

	if cmd == "help" {
		fmt.Println(usage)
	}

	if cmd == "ls" {
		Show(p)
	}

	if cmd == "add" {
		Add(p)
	}

	if cmd == "rm" {
		Delete(p)
	}
}

// start takes a slice of commands, parses flag arguments and runs the
// command if it's found.
func start(cmds []string) {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "No command is given.\n")
		return
	}

	cmdName := flag.Arg(0)

	var found bool
	for _, cmd := range cmds {
		if cmdName == cmd {
			runCmd(cmd)
			return
		}
	}

	if !found {
		fmt.Fprintf(os.Stderr, "Command %s not found.\n", cmdName)
		return
	}
}

// TokenIsSet will check to make sure an authentication token is set before
// making any API calls.
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

	cmds := []string{"help", "add", "rm", "ls"}

	start(cmds)
}
