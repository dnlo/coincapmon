Cryptocurrency monitor
======================

Monitor runs in the terminal, stays on and updates every 30s.
You can specify which currencies you want to watch (default is all) by including their symbols in a file
named "watch.txt". The app looks for the file in the current directory.

Example of watch.txt:
```
BTC
ETH
IOTA
NEO
```

You can scroll with arrows, pageUp, pageDown, home and end or you can use Vi keybindings.

Many thanks to people behind these repositories:
- (github.com/gdamore/tcell)
- (github.com/rivo/tview)
- (github.com/dustin/go-humanize)

Data is from coinmarketcap.com API