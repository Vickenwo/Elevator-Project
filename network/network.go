package network

import (
	"fmt"
	"project-group11/config"
	"strconv"
)

func CheckPeers(id int, elevators *[config.NumberOfElevators]config.ElevatorState, availableElevators *[config.NumberOfElevators]bool, netCh config.NetworkChannels, informCh config.InformChannels) {

	for {
		select {
		case p := <-netCh.PeerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

			for currentPeer := range p.Peers {
				elevatorPeer, _ := strconv.Atoi(p.Peers[currentPeer])
				availableElevators[elevatorPeer] = true
			}

			for currentPeer := range p.New {
				newId, _ := strconv.Atoi(p.New[currentPeer])
				availableElevators[newId] = true
				informCh.OrderRecoverCab <- newId
			}

			for currentPeer := range p.Lost {
				lostId, _ := strconv.Atoi(p.Lost[currentPeer])
				availableElevators[lostId] = false
				informCh.OrderRedistribute <- lostId
			}

			informCh.SendState <- elevators[id]
		}
	}
}

func SendMessage(id int, netCh config.NetworkChannels, informCh config.InformChannels) {
	for {
		select {
		case state := <-informCh.SendState:
			for i := 0; i < 10; i++ {
				netCh.SendStateNet <- state 
			}
		case newOrder := <-informCh.SendOrder:
			for i := 0; i < 10; i++ {
				netCh.SendOrderNet <- newOrder
			}
		}
	}
}
