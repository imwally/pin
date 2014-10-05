# pin

A simple command line pinboard client.

You can add, delete, and list your bookmarks right in your terminal. Incredible.

### Add a bookmark

The minimum requirements for a bookmark are a URL and a title. 

`$ pin -a http://www.sweetwebsite.com -title "One Sweet Site"`

Of course you can tag your bookmark as well. Use the `-tag` flag with space
dilimited terms.

`$ pin -a http://www.sweetwebsite.com -title "One Sweet Site" -tag "sweet site
cool"`

Need even more context for your bookmark? Use the extended `-e` flag.

`$ pin -a http://www.sweetwebsite.com -title "One Sweet Site" -tag "sweet site
cool" -e "I think this is one sweet site so I'm bookmarking it."`

You can also specify the private `-p` or toread `-r` flags.

### Delete a bookmark

The only requirement to delete a bookmark is the URL.

`$ pin -d http://www.sweetwebsite.com`

### Show your bookmarks

You can list a maximum of 100 bookmarks. The `-show` flag requires any integer
from 1 - 100.

`$ pin -show 10`

If you want more information use the long format `-l` flag.

`$ pin -show 10 -l`
