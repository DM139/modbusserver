package payload




type requestloop struct {
	ch                chan interface{}        // command channel
	idx               int                     // loop index
}

func (rl *requestloop) loopRun() {
	var err error
	//defer func() {
	//
	//}
	for v := range rl.ch {
		switch v := v.(type) {
		case error:
			err = v
		//case *stdConn:
		//	err = el.loopAccept(v)
		//case *tcpIn:
		//	err = el.loopRead(v)
		//case *udpIn:
		//	err = el.loopReadUDP(v.c)
		//case *stderr:
		//	err = el.loopError(v.c, v.err)
		//case wakeReq:
		//	err = el.loopWake(v.c)
		//case func() error:
		//	err = v()
		}
		if err != nil {
			el.svr.logger.Printf("event-loop:%d exits with error:%v\n", el.idx, err)
			break
		}
	}
}
