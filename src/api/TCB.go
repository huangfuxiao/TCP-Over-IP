package api

import (
	".././ipv4"
	".././linklayer"
	".././pkg"
	".././tcp"
	//"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

type TCB struct {
	Fd         int
	State      tcp.State
	Addr       SockAddr
	Seq        int
	Ack        int
	RecvBuffer []byte
	SendBuffer []byte
	node       *pkg.Node
	u          linklayer.UDPLink
}

type SockAddr struct {
	LocalAddr  string
	LocalPort  int
	RemoteAddr string
	RemotePort int
}

func BuildTCB(fd int, node *pkg.Node, u linklayer.UDPLink) TCB {
	buf := make([]byte, 0)
	s := tcp.State{State: 1}
	add := SockAddr{"0.0.0.0", 0, "0.0.0.0", 0}
	seqn := int(rand.Uint32())
	ackn := 0
	return TCB{fd, s, add, seqn, ackn, buf, buf, node, u}
}

func (tcb *TCB) SendCtrlMsg(ctrl int) {
	taddr := tcb.Addr
	tcph := tcp.BuildTCPHeader(taddr.LocalPort, taddr.RemotePort, tcb.Seq, tcb.Ack, ctrl, 0xaaaa)
	data := tcph.Marshal()
	tcph.Checksum = tcp.Csum(data, to4byte(taddr.LocalAddr), to4byte(taddr.RemoteAddr))
	data = tcph.Marshal()
	/*
		ipp := ipv4.BuildIpPacket(data, 6, taddr.LocalAddr, taddr.RemoteAddr)
		fmt.Println(ipp)
	*/
	//Search the interface and send to the actual address and port
	//------------TO DO--------------
	tcb.Seq += 1
	v, ok := tcb.node.RouteTable[taddr.RemoteAddr]
	if ok {
		for _, link := range tcb.node.InterfaceArray {
			if strings.Compare(v.Next, link.Src) == 0 {
				if link.Status == 0 {
					return
				}

				ipPkt := ipv4.BuildIpPacket(data, 6, taddr.LocalAddr, taddr.RemoteAddr)
				//fmt.Println(ipPkt.IpHeader.TTL)
				//fmt.Println(ipPkt.IpHeader.Protocol)
				tcb.u.Send(ipPkt, link.RemoteAddr, link.RemotePort)
				return
			}
		}

	}

}

func to4byte(addr string) [4]byte {
	parts := strings.Split(addr, ".")
	b0, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Fatalf("to4byte: %s (latency works with IPv4 addresses only, but not IPv6!)\n", err)
	}
	b1, _ := strconv.Atoi(parts[1])
	b2, _ := strconv.Atoi(parts[2])
	b3, _ := strconv.Atoi(parts[3])
	return [4]byte{byte(b0), byte(b1), byte(b2), byte(b3)}
}

/*
func receiveSynAck(laddr, raddr string) {
	for {
		buf := make([]byte, 1024)
		numRead, raddr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatalf("ReadFrom: %s\n", err)
		}
		if raddr.String() != remoteAddress {
			// this is not the packet we are looking for
			continue
		}
		tcp := NewTCPHeader(buf[:numRead])
		// Closed port gets RST, open port gets SYN ACK
		if tcp.HasFlag(RST) || (tcp.HasFlag(SYN) && tcp.HasFlag(ACK)) {
			break
		}
	}
}
*/
