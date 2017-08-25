package server

import (
	"encoding/binary"
	"io"

	"github.com/juju/ratelimit"
)

const (
	defaultHeadSize   = 4
	defaultPacketSize = 4 * 1024
)

//reader
type reader struct {
	r         io.Reader
	buf       []byte
	ratelimit *ratelimit.Bucket
}

func newReader(r io.Reader) *reader {
	var (
		rtl = new(ratelimit.Bucket)
	)
	if defaultRate <= 0.0 || defaultCapacity <= 0 {
		rtl = nil
	} else {
		rtl = ratelimit.NewBucketWithRate(defaultRate, defaultCapacity)
	}
	return &reader{
		r:         r,
		buf:       make([]byte, defaultPacketSize),
		ratelimit: rtl,
	}
}

func (r *reader) readPacket() (packet []byte, err error) {
	//read head, get packet length
	n, err := r.readHead()
	if err != nil {
		return
	}
	// read body
	if r.ratelimit != nil {
		r.ratelimit.Wait(int64(n))
	}
	_, err = r.readBody(n)
	if err != nil {
		return
	}
	packet = make([]byte, 0, n)
	packet = append(packet, r.buf[:n]...)
	return
}

func (r *reader) readHead() (hlen int, err error) {
	_, err = io.ReadFull(r.r, r.buf[:defaultHeadSize])
	if err != nil {
		return
	}
	n := binary.BigEndian.Uint32(r.buf[:defaultHeadSize])
	return int(n), nil
}

func (r *reader) readBody(blen int) (n int, err error) {
	if blen > defaultPacketSize {
		r.buf = make([]byte, blen)
	}
	return io.ReadFull(r.r, r.buf[:blen])
}

//===============================
//writer
type writer struct {
	w   io.Writer
	buf []byte
}

func newWriter(w io.Writer) *writer {
	return &writer{
		w:   w,
		buf: make([]byte, defaultHeadSize),
	}
}

func (w *writer) writePacket(packet []byte) (n int, err error) {
	//tcp stream length in header
	n, err = w.writeHead(len(packet))
	if err != nil {
		return 0, err
	}
	return w.writeBody(packet)
}

func (w *writer) writeHead(plen int) (n int, err error) {
	binary.BigEndian.PutUint32(w.buf[:defaultHeadSize], uint32(plen))
	return w.w.Write(w.buf[:defaultHeadSize])
}

func (w *writer) writeBody(body []byte) (n int, err error) {
	return w.w.Write(body)
}
