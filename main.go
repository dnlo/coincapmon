package main

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
	"time"
	"strings"
	"fmt"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/dustin/go-humanize"

)

type Currency struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Symbol string `json:"symbol"`
	Rang string `json:"rank"`
	PriceUSD string `json:"price_usd"`
	PriceBTC string`json:"price_btc"`
	USDVolume24h string `json:"24h_volume_usd"`
	MarketCapUSD string `json:"market_cap_usd"`
	AvailableSupply string `json:"available_supply"`
	TotalSupply string `json:"total_supply"`
	MaxSupply string `json:"max_supply"`
	PercentChange1h string `json:"percent_change_1h"`
	PercentChange24h string `json:"percent_change_24h"`
	PercentChange7d string `json:"percent_change_7d"`
	LastUpdated string `json:"last_updated"`
}

type Listing []Currency

func readConfig(file string) []string {
	f, err := ioutil.ReadFile(file)
	if err != nil { 
		fmt.Println("error reading config file", err)
		return nil
	}
	return strings.Split(string(f), "\n")
}

// Read json from coinmarketcap API
func data() Listing {
	listing := Listing{}
	resp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/")
	if err != nil { fmt.Println("error getting the listing from CoinMarketCap, check your internet connection: ", err) }
	data, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(data, &listing)
	return listing
}

// Determines font color based on value
// Faling is red, Rising is green
func checkStatus(v string, table *tview.Table, r int, c int) {
	if strings.HasPrefix(v, "-") {
		table.SetCell(r+1, c,
				tview.NewTableCell("  "+v).
					SetTextColor(tcell.ColorRed).
					SetAlign(tview.AlignLeft))
	} else {
		table.SetCell(r+1, c,
			tview.NewTableCell("  "+v).
				SetTextColor(tcell.ColorGreen).
				SetAlign(tview.AlignLeft))
	}
}
//Convert big numbers to ints so humanize can be used
func toInt(s string) int64 {
	sl := strings.Split(s, ".")
	s = sl[0]
	num, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Println("error converting string to int", err)
		return 0
	}
	return num
}
	
func draw(table *tview.Table, listing Listing, config []string)  {
	
	cols := []string{
		" Symbol ",
		"  Price USD ",
		"  Change 1h ", 
		"  Change 24h ", 
		"  Change 7d ", 
		"  MarketCap USD ",
		"  Volume 24h USD ",
		}
	rows := len(listing)
				
	for i := range cols {
		table.SetCell(0, i,
			tview.NewTableCell(cols[i]).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignLeft))
	}
	
	for r := 0; r < rows; r++ {
		color := tcell.ColorWhite
		if config != nil {
			for i := range config {
				if listing[r].Symbol == config[i] {
					table.SetCell(i+1, 0,
						tview.NewTableCell(" "+listing[r].Symbol).
							SetTextColor(color).
							SetAlign(tview.AlignLeft))
					table.SetCell(i+1, 1,
						tview.NewTableCell("  "+listing[r].PriceUSD).
							SetTextColor(color).
							SetAlign(tview.AlignLeft))
					checkStatus(listing[r].PercentChange1h+"%", table, i, 2)
					checkStatus(listing[r].PercentChange24h+"%", table, i, 3)
					checkStatus(listing[r].PercentChange7d+"%", table, i, 4)
					table.SetCell(i+1, 5,
						tview.NewTableCell("  "+humanize.Comma(toInt(listing[r].MarketCapUSD))).
							SetTextColor(color).
							SetAlign(tview.AlignLeft))
					table.SetCell(i+1, 6,
						tview.NewTableCell("  "+humanize.Comma(toInt(listing[r].USDVolume24h))).
							SetTextColor(color).
							SetAlign(tview.AlignLeft))
				}
			}
		} else {
			table.SetCell(r+1, 0,
						tview.NewTableCell(" "+listing[r].Symbol).
							SetTextColor(color).
							SetAlign(tview.AlignLeft))
					table.SetCell(r+1, 1,
						tview.NewTableCell("  "+listing[r].PriceUSD).
							SetTextColor(color).
							SetAlign(tview.AlignLeft))
					checkStatus(listing[r].PercentChange1h+"%", table, r, 2)
					checkStatus(listing[r].PercentChange24h+"%", table, r, 3)
					checkStatus(listing[r].PercentChange7d+"%", table, r, 4)
					table.SetCell(r+1, 5,
						tview.NewTableCell("  "+humanize.Comma(toInt(listing[r].MarketCapUSD))).
							SetTextColor(color).
							SetAlign(tview.AlignLeft))
					table.SetCell(r+1, 6,
						tview.NewTableCell("  "+humanize.Comma(toInt(listing[r].USDVolume24h))).
							SetTextColor(color).
							SetAlign(tview.AlignLeft))
		}
	}
}

func main() {
		app := tview.NewApplication()
		table := tview.NewTable().SetBorders(true).SetBordersColor(tcell.ColorBlack)
		listing := data()
		
		conf := readConfig("watch.txt")
		draw(table, listing, conf)

		go func() {
			for {
				time.Sleep(time.Second * 60)
				listing := data()
				table.Clear()
				draw(table, listing, conf)
				app.SetRoot(table, true).SetFocus(table).Draw()
			}
		}()

		table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				app.Stop()
			}
		})
		if err := app.SetRoot(table, true).SetFocus(table).Run(); err != nil {
			panic(err)
		}
}