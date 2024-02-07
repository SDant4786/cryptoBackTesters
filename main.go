package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/preichenberger/go-coinbasepro/v2"
)

var client = coinbasepro.NewClient()

var timeFrameMin = 60
var lossPerTrade = .98
var profitPrecent = .01
var stoploss = .02
var profitTake = .09
var feesPerTrade = .01
var coinsWorthTrading []ticker

type Cloud struct {
	TenkanSen      float64
	KijunSen       float64
	SenkouSpanA    float64
	SenkouSpanB    float64
	ChikouSpan     float64
	ConversionLine float64
	BaseLine       float64
}

var dontBuy = []string{
	"RAI-USD",
	"PAX-USD",
	"DAI-USD",
	"QUICK-USD",
	"BTC-USD",
	"WBTC-USD",
	"ETH-USD",
	"MKR-USD",
	"LINK-USD",
	"YFI-USD",
	"YFII-USD",
	"IOTX-USD",
}

type ticker struct {
	product  coinbasepro.Product
	bars     []coinbasepro.HistoricRate
	bought   bool
	boughtAt float64
}

func (t *ticker) asdf() {
	t.bought = true
}
func main() {
	coinsWorthTrading = getBars()
	//calcAvgMarketChanges(coinsWorthTrading)
	//wmaAlgorithm(coinsWorthTrading)
	//wmaCrossLiveTest()

	//fmt.Println("WMA")
	//wmaCrossAlgorithm(coinsWorthTrading)
	//fmt.Println("VWWMA")
	vwwmaCrossAlgorithm(coinsWorthTrading)

}

func calcAvgMarketChanges(coins []ticker) {
	totMarketChange := 0.0
	for _, c := range coins {
		start := c.bars[0].Close
		end := c.bars[len(c.bars)-1].Close
		if start > end {
			increase := end - start
			totMarketChange += (increase / start) * 100
		} else {
			decrease := start - end
			totMarketChange += (decrease / start) * 100
		}
	}
	fmt.Printf("Market:%f \n", totMarketChange/float64(len(coins)))
}
func getBars() []ticker {
	client.UpdateConfig(&coinbasepro.ClientConfig{
		BaseURL:    "",
		Key:        "",
		Passphrase: "",
		Secret:     "",
	})
	client.HTTPClient = &http.Client{
		Timeout: 15 * time.Second,
	}

	//Get Coins
	products, _ := client.GetProducts()
	allCoins := []ticker{}
	time.Sleep(time.Millisecond * 100)
	for _, p := range products {
		//If not bought, pull candles
		if strings.Contains(p.ID, "-USD") &&
			!strings.Contains(p.ID, "USDT") &&
			!strings.Contains(p.ID, "USDC") &&
			!strings.Contains(p.ID, "UST") {
			buy := true
			for _, db := range dontBuy {
				if db == p.ID {
					buy = false
					break
				}
			}
			if buy == true {
				allCoins = append(allCoins, ticker{
					product: p,
					bars:    nil,
				})
			}
		}
	}
	coinsWorthTrading := []ticker{}
	//ct := []ticker{}
	for _, c := range allCoins {
		candles, err := client.GetHistoricRates(c.product.ID, coinbasepro.GetHistoricRatesParams{
			Start:       time.Now().AddDate(0, 0, -(timeFrameMin*300/60)/24),
			End:         time.Now(),
			Granularity: timeFrameMin * 60,
		})
		if err != nil {
			log.Println(err)
		}
		for i, j := 0, len(candles)-1; i < j; i, j = i+1, j-1 {
			candles[i], candles[j] = candles[j], candles[i]
		}

		if len(candles) < 50 {
			continue
		}
		c.bars = candles
		coinsWorthTrading = append(coinsWorthTrading, c)
		time.Sleep(time.Second / 20)
	}
	return coinsWorthTrading
}
func getCoinsWorthTrading() {
	coinsWorthTrading = []ticker{}
	//Get Coins
	products, _ := client.GetProducts()
	time.Sleep(time.Millisecond * 100)

	for _, p := range products {
		//If not bought, pull candles
		if strings.Contains(p.ID, "-USD") &&
			!strings.Contains(p.ID, "USDT") &&
			!strings.Contains(p.ID, "USDC") &&
			!strings.Contains(p.ID, "UST") {
			buy := true
			for _, db := range dontBuy {
				if db == p.ID {
					buy = false
					break
				}
			}
			if buy == true {
				coinsWorthTrading = append(coinsWorthTrading, ticker{
					product:  p,
					bars:     []coinbasepro.HistoricRate{},
					bought:   false,
					boughtAt: 0.0,
				})
			}
		}
	}
}
