package unicorn

import (
	"syscall"
)

type Epoller struct {
	Fd        int
	ActiveEvs []syscall.EpollEvent
}

func NewEpoller(capacity int) (*Epoller, error) {
	fd, err := syscall.EpollCreate1(syscall.EPOLL_CLOEXEC)
	if err != nil {
		return nil, err
	}
	return &Epoller{
		Fd:        fd,
		ActiveEvs: make([]syscall.EpollEvent, capacity),
	}, nil
}

func (ep *Epoller) AddEvent(ev Event) error {
	return syscall.EpollCtl(ep.Fd, syscall.EPOLL_CTL_ADD, ev.Fd, ep.stdEv2epEv(ev))
}

func (ep *Epoller) DelEvent(ev Event) error {
	return syscall.EpollCtl(ep.Fd, syscall.EPOLL_CTL_DEL, ev.Fd, nil)
}

func (ep *Epoller) ModEvent(ev Event) error {
	return syscall.EpollCtl(ep.Fd, syscall.EPOLL_CTL_MOD, ev.Fd, ep.stdEv2epEv(ev))
}

func (ep *Epoller) WaitActiveEvents(evs []Event) (n int, err error) {
	n, err = syscall.EpollWait(ep.Fd, ep.ActiveEvs, -1)
	if err != nil && !TemporaryErr(err) {
		return
	}
	for i := 0; i < n; i++ {
		evs[i] = ep.epEv2StdEv(ep.ActiveEvs[i])
	}
	return
}

func (ep *Epoller) Close() error {
	return syscall.Close(ep.Fd)
}

func (ep *Epoller) epEv2StdEv(epEv syscall.EpollEvent) (stdEv Event) {
	stdEv.Fd = int(epEv.Fd)
	if epEv.Events&syscall.EPOLLIN != 0 {
		stdEv.Flag |= EventRead
	}
	if epEv.Events&syscall.EPOLLOUT != 0 {
		stdEv.Flag |= EventWrite
	}
	return
}

func (ep *Epoller) stdEv2epEv(stdEv Event) (epEv *syscall.EpollEvent) {
	epEv = new(syscall.EpollEvent)
	if stdEv.Flag&EventRead != 0 {
		epEv.Events |= syscall.EPOLLIN
	}
	if stdEv.Flag&EventWrite != 0 {
		epEv.Events |= syscall.EPOLLOUT
	}
	epEv.Fd = int32(stdEv.Fd)
	return
}
