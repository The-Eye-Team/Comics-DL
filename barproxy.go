package main

import "github.com/vbauerster/mpb"

type BarProxy struct {
	T int
	B *mpb.Bar
}

func (b *BarProxy) AddToTotal(by int) {
	b.T += by
	b.B.SetTotal(int64(b.T), false)
}

func (b *BarProxy) Increment(by int) {
	b.B.IncrBy(by)
}

func (b *BarProxy) FinishNow() {
	b.B.SetTotal(int64(b.T), true)
}
