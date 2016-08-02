package main
 
import (
	"net/rpc"
	"log"
	"net"
	"service"
	"yar"
	"context"
)

func main() {
	srv := rpc.NewServer()
	srv.Register(&service.Article{context.DefaultContext})

	addrPort := ":8900"
	lis, err := net.Listen("tcp", addrPort)
	if err != nil {
		log.Fatalf("Listen on %v failed: %v\n", addrPort, err)
		return
	}
	defer lis.Close()
	log.Printf("rpc service running on port %s\n", addrPort)

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("accept() failed: %v", err)
		}

		go srv.ServeCodec(yar.NewServerCodec(conn))
	}
}
