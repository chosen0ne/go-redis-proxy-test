/**
 *
 * @author  chosen0ne(louzhenlin86@126.com)
 * @date    2017-11-17 18:40:26
 */

package main

import (
	redisproxy "github.com/chosen0ne/go-redis-proxy"
	mstore "github.com/chosen0ne/go-redis-proxy-test/store"
	"runtime"
)

var (
	store *mstore.MemStore
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	proxy := redisproxy.NewRedisProxy()
	proxy.AddCommandHandler(&SetCommandHandler{redisproxy.DefaultCommandHandler{"SET", redisproxy.ARITY_EVEN}})
	proxy.AddCommandHandler(&GetCommandHandler{redisproxy.DefaultCommandHandler{"GET", 1}})
	proxy.AddCommandHandler(&MgetCommandHandler{redisproxy.DefaultCommandHandler{"MGET", redisproxy.ARITY_DFT}})

	proxy.Start("go-redis-proxy.conf")
}

type SetCommandHandler struct {
	redisproxy.DefaultCommandHandler
}

func (h *SetCommandHandler) Handle(cmd *redisproxy.Command) (redisproxy.Response, error) {
	for i := 0; i < cmd.ArgCount(); i += 2 {
		k := cmd.Arg(i)
		v := cmd.Arg(i + 1)

		store.Set(k, v)
	}

	return redisproxy.OkResponse, nil
}

type GetCommandHandler struct {
	redisproxy.DefaultCommandHandler
}

func (h *GetCommandHandler) Handle(cmd *redisproxy.Command) (redisproxy.Response, error) {
	v, ok := store.Get(cmd.Arg(0))
	if ok {
		return redisproxy.NewBulkString(v), nil
	} else {
		return redisproxy.NullResponse, nil
	}
}

type MgetCommandHandler struct {
	redisproxy.DefaultCommandHandler
}

func (h *MgetCommandHandler) Handle(cmd *redisproxy.Command) (redisproxy.Response, error) {
	array := redisproxy.NewArray()

	for i := 0; i < cmd.ArgCount(); i++ {
		if v, ok := store.Get(cmd.Arg(i)); ok {
			array.Add(redisproxy.NewBulkString(v))
		} else {
			array.Add(redisproxy.NullResponse)
		}
	}

	return array, nil
}

func init() {
	store = mstore.NewMemStore()
}
