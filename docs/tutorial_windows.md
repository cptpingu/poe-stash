# POE STASH

## Tutorial for Windows

This tutorial's goal is to quickly be able to use this tool. For advanced
options just check
[the technical documentation](https://github.com/cptpingu/poe-stash).

Download the
[latest release](https://github.com/cptpingu/poe-stash/releases/latest) and
choose the Windows archive. Then just extract everything by double clicking on
the archive (or right click on it, then choose `Extract here`).

You can either use the CLI or the graphical interface.

### How to launch the server (graphical web interface)

Double click on: `poe-stash-server.exe`.

If you encounter an error message (from Windows, or from your antivirus) do not
worry it's a false alert. This application is perfectly safe. It launches a web
server locally and use a port. Opening a port is sometime flags as suspicious by
antivirus.

A black terminal window will appear, do not close it! It's the local server.

You can then go to [http://localhost:2121](http://localhost:2121). On this page,
follow the instruction. When the profile will be generated, you will have a
download button in the top right corner.

### How to use the CLI (command line interface)

If you prefer the command line, there is a CLI available.

To use the command line, open a terminal, go in the directory where you
extracted the archive and type:
```
poe-stash-cli.exe --interactive
```
Then follow the instructions.

Or, give all options at once:
```
poe-stash-cli.exe --account <YOUR_ACCOUNT> --poesessid <YOUR_POESESSID> --league <YOUR_LEAGUE>
```
Replace:
  * YOUR_ACCOUNT: your account name (not your character name)
  * YOUR_POESESSID: the poesessid [(how to get it)](poesessid.md)
  * YOUR_LEAGUE: your league name (standard, delve, legion, ...)

It will fetch everything and generate an html file you can share with others.

