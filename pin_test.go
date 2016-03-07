package main

import (
	"fmt"
	"testing"
	"time"
)

func spinner() {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func TestGetPageTitle(t *testing.T) {
	urls := []string{
		"https://divan.github.io/posts/go_concurrency_visualize/",
		"https://github.com/imwally/pin/blob/master/pin.go",
		"https://amazon.com",
		//"https://laskdjflajlhgalkghalgkl.com",
		"http://www.cs.umd.edu/~waa/414-F11/IntroToCrypto.pdf",
	}

	for _, url := range urls {
		go spinner()
		fmt.Printf("Trying %s: \n", url)
		title, err := PageTitle(url)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("\r[\u2713] %s\n", title)
	}
}
