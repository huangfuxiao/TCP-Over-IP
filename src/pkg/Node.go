package pkg

import (
	"fmt"
)

type Node struct {
	LocalAddr      string
	Port           int
	InterfaceArray []*Interface
	RouteTable     map[string]Entry
}

type Interface struct {
	Status     int
	Src        string
	Dest       string
	RemotePort int
	RemoteAddr string
}

type Entry struct {
	Dest         string
	Next         string
	Cost         int
	Time_to_live int
}

func (n *Node) PrintInterfaces() {
	fmt.Println("id\tdst\t\tsrc\t\tenabled")
	i := 0
	for _, link := range n.InterfaceArray {
		fmt.Printf("%d\t%s\t%s\t%d\n", i, link.Dest, link.Src, link.Status)
		i++
	}
}

func (n *Node) PrintRoutes() {
	//fmt.Println("Here is the Rounting Table: \n")
	fmt.Println("dst\t\tsrc\t\tcost")
	for k, v := range n.RouteTable {
		if len(k) < 0 {
			fmt.Println("The rounting table currently has no routes!\n")
		}
		fmt.Printf("%s\t%s\t%d\n", v.Dest, v.Next, v.Cost)
	}
}

func (n *Node) InterfacesDown() {
	for _, link := range n.InterfaceArray {
		link.Status = 0
	}
}

func (n *Node) InterfacesUp() {
	for _, link := range n.InterfaceArray {
		link.Status = 1
	}
}

func (n *Node) PrepareAndSendPacket() {
	fmt.Println("Do nothing for prepareAndSendPacket for now")
}

func (n *Node) SetMTU() {
	fmt.Println("Do nothing for setMTU for now")

}
