# POE STASH

## Host your own instance of this tool as a website

It could be useful for you to host this tool as a website to ease usage. For
example if you have a guild with trusted members, you may want to provide them
an online tool.

Follow these steps!

### Step 1: securing (at least a bit)

You can run a server without any securities, but if you want to host this tool
publicly, it is advised to enable some safeguards. With the *--passwords*
option, you can choose a file containing a set of allowed users. These users are
allow to generate new profile or refresh existing ones.

Usually, creating a `pass.txt` file containing `login:password` lines is enough
(passwords are in plain text for now).

Example:
```
mylogin:mypassword
malachai:immortality
user:35*rfs
```

### Step 2: Redirecting logs in a file

When launching this server, redirect everything in a file, like
`/var/log/poe-stash/logs.txt`. You can use logrotate tool if needed.

### Step 3: Launch as a daemon

A simple trick to do that, is by using `screen`.
Either directly by using `screen <command>` or by using:
```
screen
<your command>
ctrl+a, then d
```

You can got the instance back by typing `screen -r`

### Launch the server: an example

This how the final command line could looks like:
```
screen ./poe-stash-server --passwords pass.txt --port 2121 > /var/log/poe-stash/poestash.log
```
