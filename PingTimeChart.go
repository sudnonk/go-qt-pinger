package main

import (
	"fmt"
	"github.com/therecipe/qt/charts"
	"github.com/therecipe/qt/core"
	"log"
	"math"
)

//PingのRTTを示すチャート
type PingTimeChart struct {
	core.QObject

	_ func() `constructor:"init"`

	//Qt側のチャート本体
	chart *charts.QChart
	//Qt側のチャートにのせる折れ線グラフ
	series *charts.QLineSeries
	//Pingを打つ型
	pinger *Pinger
	//QtとGo言語のやり取りをする型
	bridge *PingTimeChartBridge
	//Pingを打っている間はtrue
	isPinging bool

	//パケットロスした個数
	lossCount int
	//送信した全個数
	totalCount int
}

//QtとGo言語の間でデータをやり取りする型
type PingTimeChartBridge struct {
	core.QObject

	//新な点が追加された時にそれを送信する
	_ func(x int, y float64) `signal:"addPoint"`
	//パケットロス率が変わった時にそれを送信する
	_ func(data string) `signal:"updateLossRate"`
	//Pingを打ち始める時にそれを受信する
	_ func(ipaddr string) `slot:"startPing"`
	//Pingを止めるときにそれを受信する
	_ func() `slot:"stopPing"`
}

//初期化します
func (p *PingTimeChart) init() error {
	p.bridge = NewPingTimeChartBridge(nil)

	p.bridge.ConnectStartPing(func(ipaddr string) {
		log.Println("slot:startPing called")
		go p.startPing(ipaddr)
	})

	p.bridge.ConnectStopPing(func() {
		log.Println("slot:stopPing called")
		p.stopPing()
	})

	p.isPinging = false

	return nil
}

//Pingを打っていなかったら始める、いたら何もしない
func (p *PingTimeChart) startPing(ipaddr string) {
	if !p.isPinging {
		log.Println("startPing")

		p.reset()
		var err error
		p.pinger, err = NewPinger(ipaddr)
		if err != nil {
			log.Fatal(err)
		}

		p.isPinging = true
		go p.pinger.run()

		for r := range p.pinger.resChan {
			p.totalCount += 1
			if r.success {
				p.bridge.AddPoint(p.totalCount, math.Round(r.rtt.Seconds()*1e+5)*1e-2)
			} else {
				p.bridge.AddPoint(p.totalCount, 0)
				p.lossCount += 1
			}
			p.bridge.UpdateLossRate(fmt.Sprintf("%d/%d  %f %%", p.lossCount, p.totalCount, p.lossRate()))
		}
	} else {
		log.Println("already pinging.")
	}
}

//Pingを打っていたら止める、いなかったら何もしない
func (p *PingTimeChart) stopPing() {
	if p.isPinging {
		log.Println("stopPing")

		p.pinger.stop()
		p.reset()
	} else {
		log.Println("ping not started.")
	}
}

//パケットロス率を計算して返す
func (p *PingTimeChart) lossRate() float64 {
	return math.Round((float64(p.lossCount) / float64(p.totalCount)) * 100)
}

//この構造体を初期値に戻す
func (p *PingTimeChart) reset() {
	p.pinger = nil
	p.lossCount = 0
	p.totalCount = 0
	p.isPinging = false
}
