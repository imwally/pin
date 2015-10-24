# pin

A simple command line pinboard client.

You can add, delete, and list your bookmarks right in your terminal. Incredible.

### Setup your authentication token

You must first create a `.pinboard` configuration file in your home directory
that contains your authentication token. This token can be found on the password
tab of the settings [page](https://pinboard.in/settings/password).

### Add a bookmark

The minimum requirements to add a bookmark are a URL and a title. 

`$ pin add http://www.sweetwebsite.com -title "One Sweet Site"`

Of course you can tag your bookmark as well. Use the `-tag` flag with space
dilimited terms.

`$ pin add http://www.sweetwebsite.com -title "One Sweet Site" -tag "sweet site
cool"`

Need even more context for your bookmark? Use the `-text` flag.

`$ pin add http://www.sweetwebsite.com -title "One Sweet Site" -tag "sweet site
cool" -text "I think this is one sweet site so I'm bookmarking it."`

You can also specify the private `-private` or read later `-readlater` flags.

### Delete a bookmark

The only requirement to delete a bookmark is the URL.

`$ pin rm http://www.sweetwebsite.com`

### Show your bookmarks

You can list the most recent bookmarks.

`$ pin ls`

If you want more information use the long format `-l` flag.

`$ pin ls -l`
