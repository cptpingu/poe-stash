# POE STASH

## Tutorial for MAC

This tutorial's goal is to quickly be able to use this tool. For advanced
options just check
[the technical documentation](https://github.com/cptpingu/poe-stash).

Download the
[latest release](https://github.com/cptpingu/poe-stash/releases/latest) and
choose the MacOS archive. Then just extract everything by double clicking on
the archive.

You can either use the CLI or the graphical interface.

### How to use the CLI (command line interface)

To use the command line, open a terminal, go in the directory where you
extracted the archives and type:
```
./poe-stash-cli --interactive
```
Then follow the instructions.

Or, give all options at once:
```
./poe-stash-cli --account <YOUR_ACCOUNT> --poesessid <YOUR_POESESSID> --league <YOUR_LEAGUE>
```
Replace:
  * YOUR_ACCOUNT: your account name (not your character name)
  * YOUR_POESESSID: the poesessid [(how to get it)](poesessid.md)
  * YOUR_LEAGUE: your league name (standard, delve, legion, ...)

It will fetch everything and generate an html file you can share with others.

### How to launch the server (graphical web interface)

If you do not like the command line, there is a web graphical interface.

Double click on: `poe-stash-server`.
If you encounter the message `Your security preferences allow installation of
only apps from the App Store and identified developers`, then close this message
and right click on this binary. Select `Open With`, then `Terminal`. A window
will appear, do not close it! It's the local server.

You can then go to [http://localhost:2121](http://localhost:2121). On this page,
follow the instruction. When the profile will be generated, you will have a
download button in the top right corner.
