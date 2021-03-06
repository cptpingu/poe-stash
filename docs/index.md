# POE STASH

Share your stuff with others in Path of Exile!

Tired of not be able to show what you have to others? This tool is for you!

This project let you fetch your whole account and generate an html file with all
your stuff, browsable as if you were in game.

## Goal

Share a stash with a friend:
  * Easy
  * Nothing to install, just run it
  * Portable (should be compatible with everything)
  * Simple to use
  * Intuitive (very similar as the in-game interface)

## Download

Released versions are available here:

  * [Windows](https://gitreleases.dev/gh/cptpingu/poe-stash/latest/poe-stash-windows-amd64.zip)
  * [Linux](https://gitreleases.dev/gh/cptpingu/poe-stash/latest/poe-stash-linux-x86_64.tar.gz)
  * [MacOS](https://gitreleases.dev/gh/cptpingu/poe-stash/latest/poe-stash-darwin-x86_64.tar.gz)

## Getting started

Download archive, extract executable, launch it, and then go to:
[http://localhost:2121](http://localhost:2121) with your browser. That's all!

More detailed instructions for your system here:
  * [Windows tutorial](tutorial_windows.md)
  * [MacOS tutorial](tutorial_mac.md)
  * [Linux tutorial](tutorial_linux.md)

## Demo

Want to see what generated files look like?

[See this one for example](http://0217021.free.fr/poe-stash/demo/all_stash_types-standard.html)

I generated some more here: [Demo and screenshots](http://0217021.free.fr/poe-stash)

## Screenshots

Here is what the files generated by this tool looks like:

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

## Items pricing

This project also allows you to manage your item prices easily. It's not its
main feature, and will not replace more advanced dedicated tool, but it's still
useful (at least for me!).

Read some explanation here: [Pricing items](prices.md)

## Search

It's possible to search for item's name, characteristic or type. Note that this
is a very naive search engine. By default everything is consider an "AND". The
"OR" can be use with the "|" symbol.

Search for items with life on it:
```
life
```

Search for items with 3 elements:
```
fire cold lightning
```

Search for items with chaos or speed:
```
chaos | speed
```

Search for a ring, a map or a divination:
```
type:ring,map,divination
```

Search for a rarity unique, shaper or rare:
```
rarity:unique,rare,shaper
```

Search for an item level between 65 and 85:
```
ilvl:65,85
```

Search for a ring, an item named "Vix Lunaris" or everything which increase
damage:
```
type:ring | vix lunaris | increased damage
```

Search for a elder base two hand axe, with an item level greater than 84, with
fortify:
```
rarity:elder fortify | type:twoaxe | ilvl:84,100
```

## I need your help!

I'm facing some issues, I could use a hand:
  * English is not my first language, thus checking this documentation is
    correctly written would be great!
  * I don't have a unique premium stash type. So I don't know how to handle it.
    If you got one, send me the generated json, so I can start to support it.
  * Do you have a character with a lot of jewels? I would like to know how it
    behave in my tools.

Thanks you!

[![HitCount](http://hits.dwyl.io/cptpingu/poe-stash.svg)](http://hits.dwyl.io/cptpingu/poe-stash)
