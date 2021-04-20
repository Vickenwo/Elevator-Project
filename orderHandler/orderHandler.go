package orderHandler

import (
	"fmt"
	"project-group11/config"
	"project-group11/hardware"
)

func OrderHandler(id int, elevators *[config.NumberOfElevators]config.ElevatorState, availableElevators *[config.NumberOfElevators]bool, hardwareCh config.HardwareChannels, netCh config.NetworkChannels, informCh config.InformChannels) {
	for {
		select {
		case buttonPressed := <-hardwareCh.DrvButtons:
			newOrder := config.ElevatorOrder{buttonPressed.Floor, buttonPressed.Button, -1, false}
			// send on net
			informCh.SendOrder <- newOrder

			if buttonPressed.Button == hardware.BT_Cab {
				newOrder.DesignatedElevatorId = id
				addOrder(id, newOrder, elevators, informCh)
			} else {
				distributeOrder(id, newOrder, elevators, availableElevators, informCh)
			}
		case receivedOrder := <-netCh.ReceiveOrderNet:

			if !receivedOrder.Executed {

				if receivedOrder.Type != hardware.BT_Cab {
					distributeOrder(id, receivedOrder, elevators, availableElevators, informCh)
				} else if receivedOrder.DesignatedElevatorId == id {
					// recover own cab orders
					addOrder(id, receivedOrder, elevators, informCh)

				}
			} else {
				// order is executed
				DeleteOrder(id, receivedOrder, elevators)
			}
	
		case receivedState := <-netCh.ReceiveStateNet:
			if receivedState.Id != id && receivedState.State != config.Init {
				if receivedState.State == config.SystemFailure && elevators[receivedState.Id].State != config.SystemFailure {
					availableElevators[receivedState.Id] = false
				} else if receivedState.State != config.SystemFailure && elevators[receivedState.Id].State == config.SystemFailure {
					availableElevators[receivedState.Id] = true
				}

				elevators[receivedState.Id].Floor = receivedState.Floor
				elevators[receivedState.Id].Dir = receivedState.Dir
				elevators[receivedState.Id].State = receivedState.State
				elevators[receivedState.Id].LocalQueue = receivedState.LocalQueue
				elevators[receivedState.Id].Obstructed = receivedState.Obstructed

			} 
		case redistributeId := <-informCh.OrderRedistribute:
			for floor := 0; floor < config.NumberOfFloors; floor++ {
				for button := hardware.BT_HallUp; button < hardware.BT_Cab; button++ {
					if elevators[redistributeId].LocalQueue[floor][button] {
						redistributedOrder := config.ElevatorOrder{floor, button, id, false}
						//remove redistributed order from origin
						DeleteOrder(redistributeId, redistributedOrder, elevators)
						distributeOrder(id, redistributedOrder, elevators, availableElevators, informCh)
						informCh.SendOrder <- redistributedOrder
					}
				}
			}
		case recoverId := <-informCh.OrderRecoverCab:
			for floor := 0; floor < config.NumberOfFloors; floor++ {

				if elevators[recoverId].LocalQueue[floor][hardware.BT_Cab] {
					recoverOrder := config.ElevatorOrder{floor, hardware.BT_Cab, recoverId, false}
					informCh.SendOrder <- recoverOrder
				}
			}
		}
	}
}

