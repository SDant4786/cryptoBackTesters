package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/preichenberger/go-coinbasepro/v2"
)

// 6/5/2023 VWWMA 5:20pm, tp = 38, sp = 18, cash = 229, sl = .01
// Stop loss at 1% profit take at 4%
var tp, sp, cash, coinsBought = 38.0, 18.0, 229.0, 0
var totalFees = 0.0

func wmaCrossLiveTest() {
	getCoinsWorthTrading()
	fmt.Println("Starting...")
	ran := false
	for {
		if runCheck(5, ran) {
			vwwmaGetBarsTest()
			ran = true
			fmt.Printf("CurCash:%f | CoinsBought:%d | TotalFees:%f\n", cash, coinsBought, totalFees)
		}
		if resetBuyFlagCheck(ran) {
			ran = false
		}
		time.Sleep(time.Second)
	}
}

func runCheck(timeframe int, ran bool) bool {
	return (((time.Now().Minute() < 5 && time.Now().Minute() >= 0) ||
		(time.Now().Minute() < 20 && time.Now().Minute() >= 15) ||
		(time.Now().Minute() < 35 && time.Now().Minute() >= 30) ||
		(time.Now().Minute() < 50 && time.Now().Minute() >= 45)) &&
		ran == false)
}
func resetBuyFlagCheck(ran bool) bool {
	return ((time.Now().Minute() > 6 && time.Now().Minute() < 15) ||
		(time.Now().Minute() > 21 && time.Now().Minute() < 30) ||
		(time.Now().Minute() > 36 && time.Now().Minute() < 45) ||
		(time.Now().Minute() > 51 && time.Now().Minute() < 60)) &&
		ran == true
}
func getBarsTest() []ticker {
	client.UpdateConfig(&coinbasepro.ClientConfig{
		BaseURL:    "",
		Key:        "",
		Passphrase: "",
		Secret:     "",
	})
	client.HTTPClient = &http.Client{
		Timeout: 15 * time.Second,
	}

	for _, c := range coinsWorthTrading {
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
		//c.wmaCheckForBuySell()
		time.Sleep(time.Second / 20)
	}
	return coinsWorthTrading
}

/*
func (t *ticker) wmaCheckForBuySell() {
	swma := WMA(t.bars, float64(sp))
	twma := WMA(t.bars, float64(tp))

	lastSignal := swma[len(swma)-1]
	lastTrend := twma[len(twma)-1]
	lastBar := t.bars[len(t.bars)-1]

	if t.bought == true {
		if lastSignal < lastTrend || lastBar.Close < (1-stoploss)*t.boughtAt {
			amountBought := 1 / t.boughtAt
			cash += lastBar.Close * amountBought * lossPerTrade
			coinsBought--
			t.bought = false
		}
	} else {
		if lastSignal > lastTrend && swma[len(swma)-2] < twma[len(twma)-2] {
			t.boughtAt = lastBar.Close * lossPerTrade
			cash -= 1
			coinsBought++
			t.bought = true
		}
	}
}*/

func vwwmaGetBarsTest() {
	client.UpdateConfig(&coinbasepro.ClientConfig{
		BaseURL:    "",
		Key:        "",
		Passphrase: "",
		Secret:     "",
	})
	client.HTTPClient = &http.Client{
		Timeout: 15 * time.Second,
	}

	for i, c := range coinsWorthTrading {
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
		c.vwwmaCheckForBuySell()
		time.Sleep(time.Second / 20)
		coinsWorthTrading[i] = c
	}
}

func (t *ticker) vwwmaCheckForBuySell() {
	swma := VWWMA(t.bars, float64(sp))
	twma := VWWMA(t.bars, float64(tp))

	lastSignal := swma[len(swma)-1]
	lastTrend := twma[len(twma)-1]
	lastBar := t.bars[len(t.bars)-1]

	if t.bought == true {
		if lastBar.Close < (1-stoploss)*t.boughtAt ||
			lastBar.Close > (1+profitTake)*t.boughtAt {
			amountBought := 1 / t.boughtAt
			totalFees += .01
			cash += lastBar.Close * amountBought
			coinsBought--
			t.bought = false
		}
	} else {
		if lastSignal > lastTrend && swma[len(swma)-2] < twma[len(twma)-2] {
			t.boughtAt = lastBar.Close
			totalFees += .01
			cash -= 1
			coinsBought++
			t.bought = true
		}
	}
}
