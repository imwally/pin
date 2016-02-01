package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"strings"

	"github.com/imwally/pinboard"
)

var (
	options   = flag.NewFlagSet("", flag.ExitOnError)
	privFlag  = options.Bool("private", false, "private bookmark")
	readFlag  = options.Bool("readlater", false, "read later bookmark")
	longFlag  = options.Bool("l", false, "display long format")
	extFlag   = options.String("text", "", "longer description of bookmark")
	tagFlag   = options.String("tag", "", "space delimited tags for bookmark")
	titleFlag = options.String("title", "", "title of the bookmark")

	token string
)

var usage = `Usage: pin
  pin rm  URL
  pin add URL [OPTION]
  pin ls [OPTION]

Options:
  -title      title of bookmark being added
  -tag        space delimited tags 
  -private    mark bookmark as private
  -readlater  mark bookmark as read later
  -text       longer description of bookmark
  -l          long format for ls
`

// COUNT is the number of bookmarks to display.
const COUNT int = 50

// Piped is a helper function to check for piped input. It will return
// input, true if data was piped.
func Piped() (string, bool) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pin: %s", err)
		return "", false
	}

	isPipe := (fi.Mode() & os.ModeNamedPipe) == os.ModeNamedPipe
	if isPipe {
		read := bufio.NewReader(os.Stdin)
		line, _, err := read.ReadLine()
		if err != nil {
			fmt.Fprintf(os.Stderr, "pin: %s", err)
			return "", false
		}
		return string(line), true
	}

	return "", false
}

// PageTitle attempts to parse an HTML document for the <title> tag
// using the regexp package.
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

	// A regular expression that searches for any characters or
	// new lines within the bounds of <title> and </title>.
	re := regexp.MustCompile("<title>(?s)(.*?)(?s)</title>")
	t := string(re.FindSubmatch(body)[1])

	// If no title is found, return an error.
	if len(t) < 1 {
		return "", errors.New("pin: couldn't get page title")
	}

	// Trim new lines and white spaces from title.
	t = strings.TrimSpace(t)

	return html.UnescapeString(t), nil
}

// Add checks flag values and encodes the GET URL for adding a bookmark.
func Add(p *pinboard.Post) {
	var args []string

	// Check if URL is piped in or first argument. Optional tags
	// should follow the URL.
	if url, ok := Piped(); ok {
		p.URL = url
		args = flag.Args()[1:]
	} else {
		p.URL = flag.Args()[1]
		args = flag.Args()[2:]
	}

	// Parse flags after the URL.
	options.Parse(args)

	if *titleFlag != "" {
		p.Description = *titleFlag
	} else {
		// Use page title if title flag is not supplied.
		title, err := PageTitle(p.URL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pin: %s: %s\n", err, p.URL)
			return
		}

		p.Description = title
	}

	if *privFlag {
		p.Shared = "no"
	}

	if *readFlag {
		p.Toread = "yes"
	}

	p.Extended = *extFlag
	p.Tags = *tagFlag

	err := p.Add()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pin: %s\n", err)
	}
}

// Delete will delete the URL specified.
func Delete(p *pinboard.Post) {
	// Check if URL is piped in or first argument.
	if url, ok := Piped(); ok {
		p.URL = url
	} else {
		p.URL = flag.Args()[1]
	}

	err := p.Delete()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pin: %s\n", err)
	}
}

// Show will list the most recent bookmarks. The -tag and -readlater
// flags can be used to filter results.
func Show(p *pinboard.Post) {
	args := flag.Args()[1:]
	options.Parse(args)

	if *tagFlag != "" {
		p.Tag = *tagFlag
	}
	if *readFlag {
		p.Toread = "yes"
	}

	p.Count = COUNT

	recent := p.ShowRecent()
	for _, v := range recent.Posts {
		if *longFlag {
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
		} else {
			fmt.Println(v.Href)
		}
	}
}

// Help prints pin's usage text.
func Help(p *pinboard.Post) {
	fmt.Printf("%s", usage)
}

// start takes a slice of commands, parses flag arguments and runs the
// command if it's found.
func Start(cmds map[string]func(p *pinboard.Post)) {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "pin: no command is given.\n")
		return
	}

	cmdName := flag.Arg(0)

	cmd, ok := cmds[cmdName]
	if !ok {
		fmt.Fprintf(os.Stderr, "pin: command %s not found.\n", cmdName)
		return
	}

	// Initialise a new Pinboard post and token.
	p := new(pinboard.Post)
	p.Token = token

	cmd(p)
}

// TokenIsSet will check to make sure an authentication token is set before
// making any API calls.
func TokenIsSet() bool {
	return token != ""
}

func init() {
	u, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pin: %s", err)
	}

	content, err := ioutil.ReadFile(u.HomeDir + "/.pinboard")
	if err != nil {
		fmt.Fprintf(os.Stderr, "pin: No authorization token found. Please add your authorization token to ~/.pinboard\n")
	}

	token = string(content)
}

func main() {
	if !TokenIsSet() {
		return
	}

	cmds := map[string]func(*pinboard.Post){
		"help": Help,
		"add":  Add,
		"rm":   Delete,
		"ls":   Show,
	}

	Start(cmds)
}
