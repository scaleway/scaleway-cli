package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/creack/goselect"
)

type fder interface {
	Fd() uintptr
}

// ErrNotFder .
var ErrNotFder = fmt.Errorf("Not a fder")

type readFder interface {
	io.Reader
	fder
}

type writeFder interface {
	io.Reader
	fder
}

// SelectReader .
type selectReader struct {
	readFder
	ready chan struct{}
	stop  chan struct{}
}

func (r *selectReader) Read(buf []byte) (int, error) {
	select {
	case <-r.stop:
		return 0, io.EOF
	case <-r.ready:
		return r.readFder.Read(buf)
	}
}

// Select .
type Select struct {
	readers      []*selectReader
	writers      []writeFder
	controlPipeR *os.File
	controlPipeW *os.File
}

// NewSelect .
func NewSelect() (*Select, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	s := &Select{
		readers:      []*selectReader{},
		writers:      []writeFder{},
		controlPipeR: r,
		controlPipeW: w,
	}
	return s, nil
}

func (s *Select) run() {
	rFDSet := &goselect.FDSet{}
	var max uintptr
	var fd uintptr
	for {
		// TODO: this can be cached and changed only upon controlePipe event.
		rFDSet.Zero()
		fd = s.controlPipeR.Fd()
		max = fd
		rFDSet.Set(fd)
		for _, r := range s.readers {
			fd = r.Fd()
			if max < fd {
				max = fd
			}
			rFDSet.Set(fd)
		}
		println("-------> preselect")
		if err := goselect.Select(int(max)+1, rFDSet, nil, nil, -1); err != nil {
			log.Fatal(err)
		}
		for i := uintptr(0); i < syscall.FD_SETSIZE; i++ {
			if rFDSet.IsSet(i) {
				println(i, "is ready")
			}
		}
		println("<-------- postselect")
		for _, r := range s.readers {
			if rFDSet.IsSet(r.Fd()) {
				func(r *selectReader) { r.ready <- struct{}{} }(r)
			}
		}
		if rFDSet.IsSet(s.controlPipeR.Fd()) {
			buf := make([]byte, 1024)
			n, err := s.controlPipeR.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}
			fmt.Printf("----> control: %s\n", buf[:n])
		}
	}
}

// NewReader .
func (s *Select) NewReader(r io.Reader) (io.Reader, error) {
	rr, ok := r.(readFder)
	if !ok {
		return nil, ErrNotFder
	}
	ret := &selectReader{
		readFder: rr,
		ready:    make(chan struct{}),
		stop:     make(chan struct{}),
	}
	s.readers = append(s.readers, ret)
	if _, err := s.controlPipeW.Write([]byte("newreader")); err != nil {
		return nil, err
	}
	return ret, nil
}

func test6() error {
	s, err := NewSelect()
	if err != nil {
		return err
	}
	_ = s

	fmt.Printf("%d\n", os.Getpid())
	const MAX = 4
	rs := []io.Reader{}
	ws := []io.Writer{}
	for i := 0; i < MAX; i++ {
		r, w, _ := os.Pipe()
		println(i, "--->", r.Fd(), w.Fd())
		rs = append(rs, r)
		ws = append(ws, w)
	}

	c := make(chan struct{})
	wg := sync.WaitGroup{}
	for i := 0; i < MAX; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(2 * time.Second)
			<-c
			wg.Done()
		}()
	}
	time.Sleep(4 * time.Second)
	close(c)
	wg.Wait()
	for i := 0; i < MAX; i++ {
		time.Sleep(5 * time.Second)
		wg.Add(1)
		go func(r io.Reader, i int) {
			rr, _ := s.NewReader(r)
			r = rr
			buf := make([]byte, 1024)
			for {
				n, err := r.Read(buf)
				if err != nil {
					log.Printf("[%d] error read: %s\n", i, err)
					break
				}
				_ = n
				fmt.Printf("[%d] %s\n", i, buf[:n])
			}
			wg.Done()
		}(rs[i], i)
	}
	for i := 0; i < 1; i++ {
		for i := 0; i < MAX; i++ {
			fmt.Fprintf(ws[i], "{%d} hello", i)
		}
		time.Sleep(5 * time.Second)
	}
	wg.Wait()
	return nil
}

func test4() error {
	rFDSet := &goselect.FDSet{}
	buf := make([]byte, 1024)
	for {
		// TODO: this can be cached and changed only upon controlePipe event.
		rFDSet.Zero()
		rFDSet.Set(os.Stdin.Fd())
		if err := goselect.Select(int(os.Stdin.Fd())+1, rFDSet, nil, nil, -1); err != nil {
			return err
		}
		for i := uintptr(0); i < syscall.FD_SETSIZE; i++ {
			if rFDSet.IsSet(i) {
				println(i, "is ready")
			}
		}
		if rFDSet.IsSet(os.Stdin.Fd()) {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			fmt.Printf("---->: %s\n", buf[:n])
		}
	}
	return nil
}

type Client struct {
	f     *os.File
	queue [][]byte
}

func (c *Client) Push(msg []byte) {
	c.queue = append(c.queue, msg)
	println(len(c.queue))
}

func (c *Client) Pop() []byte {
	if len(c.queue) == 0 {
		return nil
	}
	tmp := c.queue[0]
	c.queue = c.queue[1:]
	return tmp
}

func testServer() error {
	l, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		return err
	}
	ll, ok := l.(*net.TCPListener)
	if !ok {
		return fmt.Errorf("wrong listener type")
	}

	llf, err := ll.File()
	if err != nil {
		return err
	}

	rFDSet := &goselect.FDSet{}
	wFDSet := &goselect.FDSet{}

	buf := make([]byte, 1024)
	var max, fd uintptr

	clients := []*Client{}

	for {
		// TODO: this can be cached and changed only upon controlePipe event.
		rFDSet.Zero()
		wFDSet.Zero()
		fd = llf.Fd()
		max = fd
		rFDSet.Set(fd)

		for _, c := range clients {
			fd = c.f.Fd()
			rFDSet.Set(fd)
			if len(c.queue) > 0 {
				wFDSet.Set(fd)
			}
			if max < fd {
				max = fd
			}
		}

		if err := goselect.Select(int(max)+1, rFDSet, wFDSet, nil, -1); err != nil {
			return err
		}
		println("-->")

		// Watch for new clients
		if rFDSet.IsSet(llf.Fd()) {
			c, err := ll.AcceptTCP()
			if err != nil {
				return err
			}
			f, err := c.File()
			if err != nil {
				return err
			}
			fd = f.Fd()
			println("New client:", fd)
			clients = append(clients, &Client{f: f})
		}

		// Watch for client activity
		for _, c := range clients {
			fd = c.f.Fd()
			if rFDSet.IsSet(fd) {
				n, err := c.f.Read(buf)
				if err != nil {
					return err
				}
				fmt.Printf("%s", buf[:n])
				for _, cc := range clients {
					if c != cc {
						cc.Push(buf[:n])
					}
				}
			}
		}

		// Send message to clients
		for _, c := range clients {
			fd = c.f.Fd()
			if wFDSet.IsSet(fd) {
				msg := c.Pop()
				fmt.Printf("%d write ready: %v\n", fd, msg)
				if msg != nil {
					_, err := c.f.Write(msg)
					if err != nil {
						return err
					}
				}
			}
		}
	}
}

func main() {
	if err := testServer(); err != nil {
		log.Fatal(err)
	}
}

func test5() error {
	wg := sync.WaitGroup{}
	fmt.Printf("%d\n", os.Getpid())
	for i := 0; i < 200; i++ {
		r, w, _ := os.Pipe()
		go func(w io.Writer) {
			for {
				w.Write([]byte("hello"))
				time.Sleep(time.Second)
			}
		}(w)
		wg.Add(1)
		go func(r io.Reader) {
			buf := make([]byte, 1024)
			for {
				_, err := r.Read(buf)
				if err != nil {
					log.Fatal(err)
				}
			}
		}(r)
	}
	wg.Wait()
	return nil
}
