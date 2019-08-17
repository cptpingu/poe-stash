## POE STASH

Share your stuff with others!

Scrap stash from the official **Path of Exile** website (http://pathofexile.com)
and generate a sharable file.

Main website here: https://cptpingu.github.io/poe-stash

### Demo

I put some files online, to show what the generated files look like:
http://0217021.free.fr/poe-stash/

### Project history

This project has been created because I wanted to show what I got to a friend. I
was astonished there was no way to share stash and didn't find any tool to do
that. I needed a tool which was simple to use and which generate a sharable
portable file. That's why I choose a single html file (css and javascript are
embeded).

The initial tool was a CLI which fetch the official API, and generated this html
file. Issue is, it's not easy for everyone to use the command line. So, I
created a simple http server to serve this file and remember the poe sessid. By
hosting this project, there is nothing to install! A simple browser is enough.

Then, I wanted to handle item selling. I didn't like existing tool (even if
there are more advanced). I wanted something like the premium tab stash in game.
I handled that using browser local storage (keep in mind that file are not
necessarily hosted online, so it has to work as a single html file open on a
device).

I'm heavily using this tool and satisfied with that. I often read online that
people wanted to share what they got, but have no way to do that. So I decided
to clean this project, and share it online.

### Getting started

For a quick start, read the tutorials:
  * [Windows tutorial](docs/tutorial_windows.md)
  * [MacOS tutorial](docs/tutorial_mac.md)
  * [Linux tutorial](docs/tutorial_linux.md)

### Using the generated file

Viewing items should be pretty intuitive as it works exactly like in the game.
To price an item, left click on it. To remove a price, just set an empty price.
When all prices are set, just click on "generate items shop", it will copy in
your clipboard what you need to paste on the trade forum.

More information on pricing [here](docs/prices.md)

### Hosting this tool online

[Show this explanation](/docs/hosting.md)

### Building and releasing (technical stuff!)

To release new version, I'm using two scripts which generate the archives.

For MacOS and Linux:
```
./gen_bin.sh
```

For Windows:
```
win_gen_bin.bat
```

Then I attached the generated archives to a new release.

### Working and debugging (more technical stuff!)

If you want to work on this project or just change or improve something, you
need to know some little debug things which will help you!

To avoid fetching many urls (making many http request is quite slow if you have
a lot of stash tabs!), just enable the `--cache` option. It will store all the
resulting json reponses from the server in files. The next time you will call
it, it will use the local version (much much faster). It's particularly useful
when working on the template for html, js and css generation.

Example:
```
go run cmd/cli/main.go --account cptpingu --poesessid 87f93201234f1234f --output data/cptpingu.html --cache && open data/cptpingu.html
```

### Dependencies

  * gin-tonic: https://github.com/gin-gonic/gin (`go get -u github.com/gin-gonic/gin`)
  * tippy: https://atomiks.github.io/tippyjs/

### Not supported yet

  * Map stash tab (known issue: https://www.pathofexile.com/forum/view-thread/1733474#p13674912)
  * Unique stash tab
  * Too much jewels on the character will overflow

### FAQ

Check the FAQ: [Here](/docs/faq.md)

### Contact

To contact me, either send a mail at cptpingu@gmail.com, or open an issue on
this github.
