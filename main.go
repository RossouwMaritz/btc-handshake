package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/peer"
	"github.com/btcsuite/btcd/wire"
)

var (
	direction = flag.String("direction", "inbound", "defines the connection direction (inbound/outbound)")
	address   = flag.String("address", "127.0.0.1", "defines where the address where node will run")
	port      = flag.Int("port", 18555, "defines the port where the node will be running")
)

func main() {
	// Parse input flags
	flag.Parse()
	validateInputParams()

	nodeAddress := fmt.Sprintf("%s:%d", *address, *port)
	switch *direction {
	case "inbound":
		err := startInboundPeer(nodeAddress)
		if err != nil {
			log.Fatal("failed to start inbound peer", err)
		}
	case "outbound":
		err := startOutboundPeer(nodeAddress)
		if err != nil {
			log.Fatal("failed to start outpound peer", err)
		}
	default:
		log.Fatal("invalid peer direction specified")
	}
}

func validateInputParams() {
	if *direction != "inbound" && *direction != "outbound" {
		log.Fatalf("invalid direction flag: %s", *direction)
	}

	if net.ParseIP(*address) == nil {
		log.Fatalf("invalid address flag: %s", *address)
	}

	if *port < 1 || *port > 65535 {
		log.Fatalf("invalid port flag: %d", *port)
	}
}

func startInboundPeer(nodeAddress string) error {
	peerCfg := &peer.Config{
		UserAgentName:    "inbound",
		UserAgentVersion: "1.0.0",
		ChainParams:      &chaincfg.SimNetParams,
		TrickleInterval:  time.Second * 10,
		AllowSelfConns:   true,
	}

	listener, err := net.Listen("tcp", nodeAddress)
	if err != nil {
		return err
	}

	p := peer.NewInboundPeer(peerCfg)
	fmt.Printf("starting inbound peer on %s!\n", nodeAddress)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("failed to accept connection: %v\n", err)
			continue
		}

		p.AssociateConnection(conn)
	}
}

func startOutboundPeer(nodeAddress string) error {
	ack := make(chan struct{})
	peerCfg := &peer.Config{
		UserAgentName:    "outbound",
		UserAgentVersion: "1.0.0",
		ChainParams:      &chaincfg.SimNetParams,
		Services:         0,
		TrickleInterval:  time.Second * 10,
		Listeners: peer.MessageListeners{
			OnVersion: func(p *peer.Peer, msg *wire.MsgVersion) *wire.MsgReject {
				fmt.Printf("outbound version received: %s\n", msg.UserAgent)
				return nil
			},
			OnVerAck: func(p *peer.Peer, msg *wire.MsgVerAck) {
				ack <- struct{}{}
			},
		},
		AllowSelfConns: true,
	}

	p, err := peer.NewOutboundPeer(peerCfg, nodeAddress)
	if err != nil {
		return fmt.Errorf("failed to create outbound peer %v", err)
	}

	conn, err := net.Dial("tcp", p.Addr())
	if err != nil {
		return fmt.Errorf("failed to establish tcp connection: %v", err)
	}
	p.AssociateConnection(conn)

	select {
	case <-ack:
		fmt.Println("handshake ack received!")
	case <-time.After(time.Second * 1):
		log.Fatal("handshake timeout")
	}

	p.Disconnect()
	p.WaitForDisconnect()

	return nil
}
