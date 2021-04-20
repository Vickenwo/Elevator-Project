package main

import (
	"flag"
	"project-group11/config"
	"project-group11/fsm"
	"project-group11/hardware"
	"project-group11/orderHandler"
	"project-group11/util"

	"project-group11/network"
	"project-group11/network/bcast"
	"project-group11/network/peers"
	"strconv"
)


func main() {

	var (
		localhost          string
		id                 int
		elevators          [config.NumberOfElevators]config.ElevatorState
		availableElevators [config.NumberOfElevators]bool
	)

	flag.IntVar(&id, "id", 0, "id of this peer")
	flag.StringVar(&localhost, "localhostId", "15657", "localhostId of this peer")
	flag.Parse()


	hardware.Init("localhost:"+localhost, config.NumberOfFloors)


	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network

	hardwareChannels := config.HardwareChannels{
		DrvButtons: make(chan hardware.ButtonEvent),
		DrvFloors:  make(chan int),
		DrvObstr:   make(chan bool),
	}

	networkChannels := config.NetworkChannels{
		PeerUpdateCh:       make(chan peers.PeerUpdate),
		PeerTxEnable:       make(chan bool),
		SendStateNet:       make(chan config.ElevatorState),
		ReceiveStateNet:    make(chan config.ElevatorState),
		SendOrderNet:       make(chan config.ElevatorOrder),
		ReceiveOrderNet:    make(chan config.ElevatorOrder),
	}

	informChannels := config.InformChannels{
		// Order channels
		LostCon:           make(chan config.ElevatorState),
		SendState:         make(chan config.ElevatorState),
		SendOrder:         make(chan config.ElevatorOrder),
		OrderUpdate:       make(chan bool),
		OrderRedistribute: make(chan int),
		OrderRecoverCab:   make(chan int),
		MechanicalError:   make(chan int),
	}

	timerChannels := config.TimerChannels{
		TimerDoorStart:    make(chan bool),
		TimerDoorStop:     make(chan bool),
		TimerDoorFinished: make(chan bool),
	}

	// Hardware

	go hardware.PollButtons(hardwareChannels.DrvButtons)
	go hardware.PollFloorSensor(hardwareChannels.DrvFloors)
	go hardware.PollObstructionSwitch(hardwareChannels.DrvObstr)

	// Network

	go peers.Transmitter(15647, strconv.Itoa(id), networkChannels.PeerTxEnable) 
	go peers.Receiver(15647, networkChannels.PeerUpdateCh)

	go bcast.Transmitter(15777, networkChannels.SendStateNet)
	go bcast.Receiver(15777, networkChannels.ReceiveStateNet)

	go bcast.Transmitter(15877, networkChannels.SendOrderNet)
	go bcast.Receiver(15877, networkChannels.ReceiveOrderNet)

	go network.CheckPeers(id, &elevators, &availableElevators, networkChannels, informChannels)
	go network.SendMessage(id, networkChannels, informChannels)

	// Orders

	go orderHandler.OrderHandler(id, &elevators, &availableElevators, hardwareChannels, networkChannels, informChannels)

	// Timer

	go util.Timer(timerChannels.TimerDoorStart, timerChannels.TimerDoorStop, timerChannels.TimerDoorFinished, config.TimerDoorDuration)

	// Fsm

	go fsm.CheckObstruction(id, &elevators, informChannels, hardwareChannels, timerChannels)
	fsm.Fsm(id, hardwareChannels, networkChannels, &elevators, &availableElevators, informChannels, timerChannels)
}
