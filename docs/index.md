# POE STASH

Share your stuff with others in Path of Exile!

Tired of not be able to show what you have to others?
This tool is for you!

This project let you fetch your whole account and generate an html
file with all your stuff, browsable as if you were in game.

### Goal

Share a stash with a friend:
  * Easily
  * Nothing to install (if hosted)
  * Portable (should be compatible with everything)
  * Simple to use

## Demo

[Try it yourself!](http://0217021.free.fr/poe-stash/demo/cptpingu-standard.html)

You can find more demo here: [Demo and screenshots](http://0217021.free.fr/poe-stash)

## Screenshots

This is how this tool looks like:

**MainScreen**
![MainScreen](http://0217021.free.fr/poe-stash/screenshots/MainScreen.png)

**Stash tabs**
![Stash tabs](http://0217021.free.fr/poe-stash/screenshots/Stash%20tabs.png)

**Special stash tabs**
![Special stash tabs](http://0217021.free.fr/poe-stash/screenshots/Special%20stash%20tabs.png)

**Inventoy**
![Inventoy](http://0217021.free.fr/poe-stash/screenshots/Inventoy.png)

**Mouseover**
![Mouseover](http://0217021.free.fr/poe-stash/screenshots/Mouseover.png)

**Set a price**
![Set a price](http://0217021.free.fr/poe-stash/screenshots/Set%20a%20price.png)

**Price thumbnails**
![Price thumbnails](http://0217021.free.fr/poe-stash/screenshots/Price%20thumbnails.png)

**Price highlights**
![Price highlights](http://0217021.free.fr/poe-stash/screenshots/Price%20highlights.png)

**Shop generated**
![Shop generated](http://0217021.free.fr/poe-stash/screenshots/Shop%20generated.png)

**Server main page**
![Server main page](http://0217021.free.fr/poe-stash/screenshots/Server%20main%20page.png)

## Bonus

This project also allows you to manage your item prices easily.
Read the whole doc here: [https://github.com/cptpingu/poe-stash](https://github.com/cptpingu/poe-stash)

## Usage

### Generate a file with the CLI
```
go run cmd/cli/main.go --account <YOUR_ACCOUNT> --poesessid <YOUR_POESESSID> --league <YOUR_LEAGUE>
```
Example:
```
go run cmd/cli/main.go --account cptpingu --poesessid ef87f9320ba7428149fe562236e32 --league standard
```
A file "cptpingu-standard.html" should be created.

### Launch the server (graphical interface)

If you don't like the command line, you can launch the server locally:
```
go run cmd/server/main.go
```

Then, go to: `http://localhost:2121` with your browser. That's all!
