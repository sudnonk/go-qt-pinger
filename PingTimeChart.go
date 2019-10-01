package main

import (
	"fmt"
	"github.com/therecipe/qt/charts"
	"github.com/therecipe/qt/core"
	"log"
	"math"
)

type PingTimeChart struct {
	core.QObject

	_ func() `constructor:"init"`

	chart     *charts.QChart
	series    *charts.QLineSeries
	pinger    *Pinger
	bridge    *PingTimeChartBridge
	isPinging bool

	lossCount  int
	totalCount int
	avg        float32
}

type PingTimeChartBridge struct {
	core.QObject

	_ func(x int, y float64) `signal:"addPoint"`
	_ func(data string)      `signal:"updateLossRate"`
	_ func(ipaddr string)    `slot:"startPing"`
	_ func()                 `slot:"stopPing"`
}

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

func (p *PingTimeChart) stopPing() {
	if p.isPinging {
		log.Println("stopPing")

		p.pinger.stop()
		p.reset()
	} else {
		log.Println("ping not started.")
	}
}

func (p *PingTimeChart) lossRate() float64 {
	return math.Round((float64(p.lossCount) / float64(p.totalCount)) * 100)
}

func (p *PingTimeChart) reset() {
	p.pinger = nil
	p.lossCount = 0
	p.totalCount = 0
	p.avg = 0
	p.isPinging = false
}
