package main

import "fmt"

func wmaAlgorithm(coins []ticker) {
	type wmaTrade struct {
		period     float64
		weight     float64
		cash       float64
		trades     int
		profitable float64
	}
	trades := []wmaTrade{}
	for period := 2.0; period < 50; period++ {
		trade := wmaTrade{
			period: period,
			cash:   0.0,
			trades: 0,
		}
		for _, c := range coins {
			wma := WMA(c.bars, period)
			bought := false
			boughtAt := 0.0
			cash := 1.0
			numTrades := 0
			profitable := 0.0
			wmai := 1
			for candle := period + 1; candle < float64(len(c.bars)); candle++ {
				if wma[wmai] > wma[wmai-1] && bought == false {
					boughtAt = c.bars[int(candle)].Close * lossPerTrade
					bought = true
					numTrades++
				}
				if bought == true && wma[wmai] < wma[wmai-1] {
					prevCash := cash
					amountBought := cash / boughtAt
					cash = c.bars[int(candle)].Close * amountBought * lossPerTrade
					bought = false
					if cash > prevCash {
						profitable++
					}
				}
				wmai++
			}
			if bought == true {
				prevCash := cash
				amountBought := cash / boughtAt
				cash = c.bars[len(c.bars)-1].Close * amountBought * lossPerTrade
				if cash > prevCash {
					profitable++
				}
			}
			trade.cash = cash - 1
			trade.trades = numTrades
			trade.profitable = profitable
			trades = append(trades, trade)
		}
	}
	profitable := 0
	totProfitable := 0
	totPeriod := 0.0
	totTrades := 0
	totCash := 0.0
	for _, t := range trades {
		if t.profitable > float64(t.trades)*profitPrecent {
			profitable++
			totPeriod += t.period
			totTrades += t.trades
			totCash += t.cash
			totProfitable += int(t.profitable)
		}
	}
	total := len(trades)
	avgPeriod := totPeriod / float64(profitable)
	avgTrades := totTrades / profitable
	avgProfitable := totProfitable / profitable

	avgCash := totCash / float64(profitable)

	fmt.Printf("Total:%d | AvgProfitable:%d | AvgCash:%f | AvgPeriod %f | AvgTrades %d \n",
		total,
		avgProfitable,
		avgCash,
		avgPeriod,
		avgTrades)

	simCash := 0.0
	simTrades := 0
	intPeriod := int(avgPeriod)
	profitableTrades := 0
	for _, c := range coins {
		wma := WMA(c.bars, float64(intPeriod))
		bought := false
		boughtAt := 0.0
		cash := 1.0
		wmai := 1
		for candle := intPeriod + 1; candle < len(c.bars); candle++ {
			if wma[wmai] > wma[wmai-1] && bought == false {
				boughtAt = c.bars[int(candle)].Close * lossPerTrade
				bought = true
				simTrades++
			}
			if bought == true && wma[wmai] < wma[wmai-1] {
				prevCash := cash
				amountBought := cash / boughtAt
				cash = c.bars[int(candle)].Close * amountBought * lossPerTrade
				bought = false
				if cash > prevCash {
					profitableTrades++
				}
			}
			wmai++
		}
		if bought == true {
			prevCash := cash
			amountBought := cash / boughtAt
			cash = c.bars[len(c.bars)-1].Close * amountBought * lossPerTrade
			if cash > prevCash {
				profitableTrades++
			}
		}
		simCash += cash - 1
	}
	fmt.Printf("Total Coins:%d | Total Profit: %f | Total Trades %d | Period:%d | Profitable:%d",
		len(coins),
		simCash,
		simTrades,
		intPeriod,
		profitableTrades,
	)
}

func wmaCrossAlgorithm(coins []ticker) {
	type wmaTrade struct {
		trendPeriod  float64
		signalPeriod float64
		cash         float64
		trades       int
		profitable   float64
		stopLossHits int
	}
	trades := []wmaTrade{}
	for signalPeriod := 2.0; signalPeriod < 50; signalPeriod++ {
		for trendPeriod := signalPeriod + 1; trendPeriod < 50; trendPeriod++ {
			trade := wmaTrade{
				trendPeriod:  trendPeriod,
				signalPeriod: signalPeriod,
				cash:         0.0,
				trades:       0,
			}
			for _, c := range coins {
				swma := WMA(c.bars, signalPeriod)
				twma := WMA(c.bars, trendPeriod)
				bought := false
				boughtAt := 0.0
				cash := 1.0
				numTrades := 0
				profitable := 0.0
				stopLossHits := 0
				swmai := int(trendPeriod - signalPeriod)
				twmai := 0
				for candle := trendPeriod; candle < float64(len(c.bars)); candle++ {
					if swma[swmai] > twma[twmai] && bought == false {
						boughtAt = c.bars[int(candle)].Close * lossPerTrade
						bought = true
						numTrades++
					}
					if bought == true && swma[swmai] < swma[twmai] {
						prevCash := cash
						amountBought := cash / boughtAt
						cash = c.bars[int(candle)].Close * amountBought * lossPerTrade
						bought = false
						if cash > prevCash {
							profitable++
						}
					}
					if bought == true && c.bars[int(candle)].Close < (1-stoploss)*boughtAt {
						amountBought := cash / boughtAt
						cash = c.bars[int(candle)].Close * amountBought * lossPerTrade
						bought = false
						stopLossHits++
					}
					swmai++
					twmai++
				}
				if bought == true {
					prevCash := cash
					amountBought := cash / boughtAt
					cash = c.bars[len(c.bars)-1].Close * amountBought * lossPerTrade
					if cash > prevCash {
						profitable++
					}
				}
				trade.cash = cash - 1
				trade.trades = numTrades
				trade.profitable = profitable
				trade.stopLossHits = stopLossHits
				trades = append(trades, trade)
			}
		}
	}
	profitable := 0
	totProfitable := 0
	totSPeriod := 0.0
	totTPeriod := 0.0
	totTrades := 0
	totCash := 0.0
	totStoploss := 0
	for _, t := range trades {
		if t.profitable > float64(t.trades)*profitPrecent {
			profitable++
			totSPeriod += t.signalPeriod
			totTPeriod += t.trendPeriod
			totTrades += t.trades
			totCash += t.cash
			totProfitable += int(t.profitable)
			totStoploss += t.stopLossHits
		}
	}
	total := len(trades)
	avgSPeriod := totSPeriod / float64(profitable)
	avgTPeriod := totTPeriod / float64(profitable)
	avgTrades := totTrades / profitable
	avgProfitable := totProfitable / profitable
	avgStopLossHits := totStoploss / profitable
	avgCash := totCash / float64(profitable)

	fmt.Printf("Total:%d | TotalProfitable: %d | AvgProfitable:%d | AvgCash:%f | AvgSPeriod %f | AvgTPeriod %f | AvgTrades %d | AvgStopLoss %d \n",
		total,
		profitable,
		avgProfitable,
		avgCash,
		avgSPeriod,
		avgTPeriod,
		avgTrades,
		avgStopLossHits,
	)

	simCash := 0.0
	simTrades := 0
	intSPeriod := int(avgSPeriod)
	intTPeriod := int(avgTPeriod)
	profitableTrades := 0
	stopLossHits := 0
	swmai := 1
	twmai := 1
	for _, c := range coins {
		swma := WMA(c.bars, float64(intSPeriod))
		twma := WMA(c.bars, float64(intTPeriod))
		bought := false
		boughtAt := 0.0
		cash := 1.0
		for candle := intTPeriod; candle < len(c.bars)-intTPeriod; candle++ {
			if swma[swmai] > twma[twmai] && bought == false {
				boughtAt = c.bars[int(candle)].Close
				amountBought := 1 / boughtAt
				totalFees += boughtAt * amountBought * feesPerTrade
				bought = true
				simTrades++
			}
			if bought == true && swma[swmai] < twma[twmai] {
				prevCash := cash
				amountBought := cash / boughtAt
				cash = c.bars[int(candle)].Close * amountBought
				totalFees += c.bars[int(candle)].Close * amountBought * feesPerTrade
				bought = false
				if cash > prevCash {
					profitableTrades++
				}
			}
			if bought == true && c.bars[int(candle)].Close < (1-stoploss)*boughtAt {
				amountBought := cash / boughtAt
				cash = c.bars[int(candle)].Close * amountBought
				totalFees += c.bars[int(candle)].Close * amountBought * feesPerTrade
				bought = false
				stopLossHits++
			}
		}
		/*
			if bought == true {
				prevCash := cash
				amountBought := cash / boughtAt
				cash = c.bars[len(c.bars)-1].Close * amountBought * lossPerTrade
				if cash > prevCash {
					profitableTrades++
				}
			}
		*/
		simCash += cash - 1
	}
	fmt.Printf("Total Coins:%d | Total Profit: %f | Total Trades %d | SPeriod:%d | TPeriod:%d | Profitable:%d | Stoploss:%d \n",
		len(coins),
		simCash,
		simTrades,
		intSPeriod,
		intTPeriod,
		profitableTrades,
		stopLossHits,
	)
}

func vwwmaCrossAlgorithm(coins []ticker) {
	type wmaTrade struct {
		trendPeriod  float64
		signalPeriod float64
		fees         float64
		cash         float64
		trades       int
		profitable   float64
		stopLossHits int
	}
	trades := []wmaTrade{}
	for signalPeriod := 2.0; signalPeriod < 100; signalPeriod++ {
		for trendPeriod := signalPeriod + 1; trendPeriod < 100; trendPeriod++ {
			trade := wmaTrade{
				trendPeriod:  trendPeriod,
				signalPeriod: signalPeriod,
				fees:         0.0,
				cash:         0.0,
				trades:       0,
			}
			for _, c := range coins {
				swma := VWWMA(c.bars, signalPeriod)
				twma := VWWMA(c.bars, trendPeriod)
				bought := false
				boughtAt := 0.0
				cash := 1.0
				numTrades := 0
				profitable := 0.0
				stopLossHits := 0
				swmai := int(trendPeriod-signalPeriod) + 1
				twmai := 0 + 1
				for candle := trendPeriod + 1; candle < float64(len(c.bars)); candle++ {
					if bought == false {
						if swma[swmai] > twma[twmai] &&
							swma[swmai-1] < twma[twmai-1] {
							boughtAt = c.bars[int(candle)].Close
							bought = true
						}
					} else {
						if c.bars[int(candle)].Close < (1-stoploss)*boughtAt {
							amountBought := cash / boughtAt
							cash = c.bars[int(candle)].Close * amountBought
							bought = false
							trade.fees += .02
							stopLossHits++
							numTrades++
						} else if c.bars[int(candle)].Close > (1+profitTake)*boughtAt {
							amountBought := cash / boughtAt
							cash = c.bars[int(candle)].Close * amountBought
							bought = false
							trade.fees += .02
							profitable++
							numTrades++
						}
					}
					swmai++
					twmai++
				}
				trade.cash = cash - 1
				trade.trades = numTrades
				trade.profitable = profitable
				trade.stopLossHits = stopLossHits
				trades = append(trades, trade)
			}
		}
	}
	profitable := 0
	totProfitable := 0
	totSPeriod := 0.0
	totTPeriod := 0.0
	totTrades := 0
	totCash := 0.0
	totStoploss := 0
	for _, t := range trades {
		//if t.profitable > float64(t.trades)*profitPrecent {
		if t.cash-t.fees > 0.0 {
			profitable++
			totSPeriod += t.signalPeriod
			totTPeriod += t.trendPeriod
			totTrades += t.trades
			totCash += t.cash
			totProfitable += int(t.profitable)
			totStoploss += t.stopLossHits
		}
	}
	total := len(trades)
	avgSPeriod := totSPeriod / float64(profitable)
	avgTPeriod := totTPeriod / float64(profitable)
	avgTrades := totTrades / profitable
	avgProfitable := totProfitable / profitable
	avgStopLossHits := totStoploss / profitable
	avgCash := totCash / float64(profitable)

	fmt.Printf("Total:%d | TotalProfitable: %d | AvgProfitable:%d | AvgCash:%f | AvgSPeriod %f | AvgTPeriod %f | AvgTrades %d | AvgStopLoss %d \n",
		total,
		profitable,
		avgProfitable,
		avgCash,
		avgSPeriod,
		avgTPeriod,
		avgTrades,
		avgStopLossHits,
	)

	simCash := 0.0
	simTrades := 0
	/*
		intSPeriod := 18
		intTPeriod := 38
	*/
	intSPeriod := int(avgSPeriod)
	intTPeriod := int(avgTPeriod)
	profitableTrades := 0
	stopLossHits := 0
	fees := 0.0

	swmai := intTPeriod - intSPeriod + 1
	twmai := 1
	for _, c := range coins {
		swma := VWWMA(c.bars, float64(intSPeriod))
		twma := VWWMA(c.bars, float64(intTPeriod))
		bought := false
		boughtAt := 0.0
		cash := 1.0
		for candle := intTPeriod + 1; candle < len(c.bars)-intTPeriod; candle++ {
			if bought == false {
				if swma[swmai] > twma[twmai] &&
					swma[swmai-1] < twma[twmai-1] {
					boughtAt = c.bars[int(candle)].Close
					bought = true
				}
			} else {
				if c.bars[int(candle)].Close < (1-stoploss)*boughtAt {
					amountBought := cash / boughtAt
					cash = c.bars[int(candle)].Close * amountBought
					bought = false
					fees += .02
					stopLossHits++
					simTrades++
				} else if c.bars[int(candle)].Close > (1+profitTake)*boughtAt {
					amountBought := cash / boughtAt
					cash = c.bars[int(candle)].Close * amountBought
					bought = false
					fees += .02
					profitableTrades++
					simTrades++
				}
			}
		}
		simCash += cash - 1
	}
	fmt.Printf("Total Coins:%d | Total Profit: %f | Total Trades %d | SPeriod:%d | TPeriod:%d | Profitable:%d | Stoploss:%d | Fees:%f\n",
		len(coins),
		simCash,
		simTrades,
		intSPeriod,
		intTPeriod,
		profitableTrades,
		stopLossHits,
		fees,
	)
}
