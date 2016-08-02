// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yar

import (
	"errors"
	"io"
	"net/rpc"
	"sync"
)

var errMissingParams = errors.New("jsonrpc: request body missing params")

type YarServerCodec struct {
	r io.Reader
	w io.Writer
	c io.Closer

	// temporary work space
	req    serverRequest
	packer Packager

	mutex   sync.Mutex // protects seq, pending
	seq     uint64
	pending map[uint64]uint32
}

// NewServerCodec returns a serverCodec that communicates with the ClientCodec
// on the other end of the given conn.
func NewServerCodec(conn io.ReadWriteCloser) rpc.ServerCodec {
	return &YarServerCodec{
		r:       conn,
		w:       conn,
		c:       conn,
		pending: make(map[uint64]uint32),
	}
}

func (c *YarServerCodec) ReadRequestHeader(r *rpc.Request) error {
	c.req.Reset()
	packer, err := readPack(c.r, &c.req)
	if err != nil {
		return err
	}

	c.packer = packer
	r.ServiceMethod = c.req.Method
	c.mutex.Lock()
	c.seq++
	c.pending[c.seq] = uint32(c.req.Id)
	c.req.Id = 0
	r.Seq = c.seq
	c.mutex.Unlock()

	return nil
}

func (c *YarServerCodec) ReadRequestBody(x interface{}) error {
	if x == nil {
		return nil
	}
	if c.req.Params == nil {
		return errMissingParams
	}

	err := c.packer.Unmarshal(*c.req.Params, &x)
	return err
}

func (c *YarServerCodec) WriteResponse(r *rpc.Response, x interface{}) error {
	c.mutex.Lock()
	id, ok := c.pending[r.Seq]
	if !ok {
		c.mutex.Unlock()
		return errors.New("invalid sequence number in response")
	}

	delete(c.pending, r.Seq)
	c.mutex.Unlock()

	resp := serverResponse{
		Id:     int64(id),
		Error:  "",
		Result: nil,
		Output: "",
		Status: 0,
	}

	// FIXME: 如果某个字段是空的map，php端会得到class，是否可以将map转成array？
	if r.Error == "" {
		resp.Result = &x
	} else {
		resp.Error = r.Error
	}

	Id := (uint32)(resp.Id)
	err := writePack(c.w, c.packer, Id, &resp)
	return err
}

func (s *YarServerCodec) Close() error {
	//s.c.Close()     // TODO: 为何这里注释掉？
	return nil
}

