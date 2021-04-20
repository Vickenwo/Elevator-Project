package config

import (
	"project-group11/hardware"
	"project-group11/network/peers"
	"time"
)

//Constants

const (
	NumberOfFloors    int = 4
	NumberOfElevators     = 3
	NumberOfButtons       = 3

	TimerDoorDuration          time.Duration = 3 * time.Second
	TimerMechanicCheckDuration time.Duration = 1 * time.Second

	TimeoutMechanicalError time.Duration = 5 * time.Second
)

type States int

const (
	Init          States = 0
	Idle                 = 1
	Moving               = 2
	Executing            = 3
	SystemFailure        = 4
)

//Types

type ElevatorState struct {
	Id         int
	Floor      int
	Dir        hardware.MotorDirection
	State      States
	LocalQueue [NumberOfFloors][NumberOfButtons]bool
	Obstructed bool
}

type ElevatorOrder struct {
	Floor                int
	Type                 hardware.ButtonType
	DesignatedElevatorId int
	Executed             bool
}


//Channels

type HardwareChannels struct {
	DrvButtons chan hardware.ButtonEvent
	DrvFloors  chan int
	DrvObstr   chan bool
}

type NetworkChannels struct {
	SendOrderNet       chan ElevatorOrder
	ReceiveOrderNet    chan ElevatorOrder
	PeerUpdateCh       chan peers.PeerUpdate
	PeerTxEnable       chan bool
	SendStateNet       chan ElevatorState
	ReceiveStateNet    chan ElevatorState
}

type InformChannels struct {
	LostCon           chan ElevatorState
	SendState         chan ElevatorState
	SendOrder         chan ElevatorOrder
	OrderUpdate       chan bool
	OrderRedistribute chan int
	OrderRecoverCab   chan int
	MechanicalError   chan int
}

type TimerChannels struct {
	TimerDoorStart    chan bool
	TimerDoorStop     chan bool
	TimerDoorFinished chan bool
}
