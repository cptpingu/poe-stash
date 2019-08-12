## POE STASH

Share your stuff with others!
Scrap stash from the official **Path of Exile** website
(http://pathofexile.com) and generate a sharable file.
Main website here: https://cptpingu.github.io/poe-stash

### Goal

Share a stash with a friend:
  * Easily
  * Nothing to install (if hosted)
  * Portable (should be compatible with everything)
  * Simple to use

### Demo

I put some files online, to show what the generated files look like:
http://0217021.free.fr/poe-stash/

### Project history

This project has been created because I wanted to show what I got to a
friend. I was astonished there was no way to share stash and didn't
find any tool to do that. I needed a tool which was simple to use and
which generate a sharable portable file. That's why I choose a single
html file (css and javascript are embeded).

The initial tool was a CLI which fetch the official API, and generated
this html file. Issue is, it's not easy for everyone to use the
command line. So, I created a simple http server to serve this file
and remember the poe sessid. By hosting this project, there is nothing
to install! A simple browser is enough.

Then, I wanted to handle item selling. I didn't like existing tool
(even if there are more advanced). I wanted something like the
premium tab stash in game. I handled that using browser local storage
(keep in mind that file are not necessarily hosted online, so it has
to work as a single html file open on a device).

I'm heavily using this tool and satisfied with that. I often read
online that people wanted to share what they got, but have no way to
do that. So I decided to clean this project, and share it online.

### Generate a file with the CLI

```
go install cmd/cli
```

You can use `cli` or `go run cmd/cli/main.go`. I will use the later in
my example.

```
go run cmd/cli/main.go --account <YOUR_ACCOUNT> --poesessid <YOUR_POESESSID> --league <YOUR_LEAGUE>
```

Account is your account name (not your character name)
Poesessid is the token id used after log in on the official website (in storage > cookie)
League, the league name (standard, legion, ...)

Example:
```
go run cmd/cli/main.go --account cptpingu --poesessid ef87f9320ba7428149fe562236e32 --league standard
```
A file "cptpingu-standard.html" should be created.

Type --help for a description of all other existing options.

### Launch the server (graphical interface)

If you don't like the command line, you can launch the server locally:
```
go run cmd/server/main.go
```

Then, go to: `http://localhost:2121` with your browser. That's all!

## Passwords

Nothing is needed to view files, but a login/pass is required for
generation. You need to create a `pass.txt` file containing
`login:password` lines.
Example:
```
mylogin:mypassword
malachai:immortality
user:35*rfs
```

### Note about the POE sessid

Never ever share your POE sessid! This token is used to identify you
on the official website. This tool need it to fetch all information
about your account, but you should not share it. If a website ask your
poesessid (and it's not launched by you locally), you should not trust
it!

### Using the generated file

Viewing items should be pretty intuitive as it works exactly like in
the game. To price an item, left click on it. To remove a price, just
set an empty price. When all prices are set, just click on "generate
items shop", it will copy in your clipboard what you need to paste on
the trade forum.

Note that every prices are store locally in your browser, and the
storage is associated with the exact url name. It means all prices
will be lost if you:
  * Rename the file
  * Use another browser
  * Clear the browser cache
  * Reinstall your browser
  * Open the file on another device

To avoid that, there are an `import shop` and an `export shop` buttons,
which allows you to save and reload the prices you set.

### Dependencies

  * gin-tonic: https://github.com/gin-gonic/gin (`go get -u github.com/gin-gonic/gin`)
  * tippy: https://atomiks.github.io/tippyjs/

### Not supported yet

  * Map stash tab (known issue: https://www.pathofexile.com/forum/view-thread/1733474#p13674912)
  * Quad stash tab
  * Essence stash tab
  * Unique stash tab
  * Divination stash tab
  * Remove only stash tab
  * Too much jewels on the character will overflow

### TODO
  * More rigourous item description generation
  * Shop id for the link after shop generation
  * Search bar? (lot of works)
