package main

import (
	"github.com/tatsushid/go-fastping"
	"log"
	"net"
	"time"
)

//Pingの応答を格納する型
type response struct {
	//帰って来たIPアドレス
	addr *net.IPAddr
	//RTT
	rtt time.Duration
	//成功ならtrue
	success bool
}

//Pingを打つ型
type Pinger struct {
	pinger *fastping.Pinger

	resChan   chan *response
	isPinging bool
	res       *response
}

//IPアドレスを受け取ってそれで初期化します
func (p *Pinger) init(ipaddr string) (*Pinger, error) {
	p.pinger = fastping.NewPinger()

	addr, err := net.ResolveIPAddr("ip4:icmp", ipaddr)
	if err != nil {
		log.Fatal(err)
		return p, err
	}

	p.pinger.AddIPAddr(addr)

	p.pinger.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		log.Println("ping received.")
		p.res = &response{addr: addr, rtt: rtt, success: true}
	}
	p.pinger.OnIdle = func() {
		if p.res == nil { //p.resがnilなのにIdleになったらパケットロス
			p.resChan <- &response{addr: nil, rtt: 0, success: false}
		} else {
			p.resChan <- p.res
		}
	}

	p.pinger.MaxRTT = time.Second

	p.resChan = make(chan *response)

	p.isPinging = false

	return p, nil
}

//Pingが始まっていなかったら開始します
func (p *Pinger) run() {
	if !p.isPinging {
		log.Println("pinger run")

		p.isPinging = true
		p.pinger.RunLoop()

	loop:
		for {
			select {
			case <-p.pinger.Done():
				log.Println("pinger done")
				if err := p.pinger.Err(); err != nil { //pinger.Stop()以外でDone()になったらエラー
					log.Fatal(err)
				}
				break loop
			}
		}

		log.Println("ping finished.")
		close(p.resChan)
		p.isPinging = false
	} else {
		log.Println("ping is already running.")
	}
}

//Pingが始まっていたら止めます
func (p *Pinger) stop() {
	if p.isPinging {
		log.Println("pinger stop")

		p.pinger.Stop()
		p.isPinging = false
	} else {
		log.Println("pinger is not running.")
	}
}

//新しいPinger構造体を返す
func NewPinger(ipaddr string) (*Pinger, error) {
	p := new(Pinger)
	return p.init(ipaddr)
}
