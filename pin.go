package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"regexp"

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

var usage = `Usage: pin
  pin rm  URL
  pin add URL [-title title] [OPTIONS]
  pin ls [OPTIONS]

Options:
  -tag        space delimited tags 
  -private    mark bookmark as private
  -readlater  mark bookmark as read later
  -text       longer description of bookmark
  -l          long format for ls
`

// Number of bookmarks to display.
const COUNT int = 50

// PageTitle returns the title from an HTML page.
func PageTitle(url string) (title string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile("<title>(.*?)</title>")

	return string(re.FindSubmatch(body)[1]), nil
}

// Add checks flag values and encodes the GET URL for adding a bookmark.
func Add(p pinboard.Post) {

	// Check stdin first.
	read := bufio.NewReader(os.Stdin)
	line, _, err := read.ReadLine()
	if err != nil {
		fmt.Println(err)
	}

	// Use first argument if stdin is empty.
	if len(flag.Args()[1]) > 0 {
		p.URL = flag.Args()[1]
	} else {
		p.URL = string(line)
	}

	// Parse flags after the URL.
	args := flag.Args()[1:]
	options.Parse(args)

	title, err := PageTitle(p.URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pin: couldn't get title: %s", err)
	} else {
		p.Description = title
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
	err = p.Add()
	if err != nil {
		fmt.Println(err)
	}
}

// Delete will delete the URL specified.
func Delete(p pinboard.Post) {

	// Check stdin first.
	read := bufio.NewReader(os.Stdin)
	line, _, err := read.ReadLine()
	if err != nil {
		fmt.Println(err)
	}

	// Use first argument if stdin is empty.
	if len(line) > 0 {
		p.URL = string(line)
	} else {
		p.URL = flag.Args()[1]
	}

	p.Encode()
	err = p.Delete()
	if err != nil {
		fmt.Println(err)
	}
}

// Show will list the most recent bookmarks. The -tag and -readlater
// flags can be used to filter results.
func Show(p pinboard.Post) {

	args := flag.Args()[1:]
	options.Parse(args)

	if *tagFlag != "" {
		p.Tag = *tagFlag
	}
	if *readFlag {
		p.Toread = "yes"
	}

	p.Count = COUNT
	p.Encode()

	recent := p.ShowRecent()

	if *longFlag {
		for _, v := range recent.Posts {
			var shared, unread string
			if v.Shared == "no" {
				shared = "[*]"
			}
			if v.Toread == "yes" {
				unread = "[#]"
			}
			fmt.Println(unread + shared + v.Description)
			fmt.Println(v.Href)
			if v.Extended != "" {
				fmt.Println(v.Extended)
			}
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
		fmt.Printf("%s", usage)
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
		fmt.Fprintf(os.Stderr, "pin: no command is given.\n")
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
		fmt.Fprintf(os.Stderr, "pin: command %s not found.\n", cmdName)
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
