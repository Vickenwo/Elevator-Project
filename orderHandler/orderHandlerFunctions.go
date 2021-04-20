package orderHandler

import (
	"project-group11/config"
	"project-group11/hardware"
)

const (
	weightDirection int = 3
	weightDistance      = 1
	weightState         = 1
	weightOrders        = 1
)

func addOrder(id int, newOrder config.ElevatorOrder, elevators *[config.NumberOfElevators]config.ElevatorState, informCh config.InformChannels) {
	if id == newOrder.DesignatedElevatorId {
		hardware.SetButtonLamp(newOrder.Type, newOrder.Floor, true)
	}
	if newOrder.DesignatedElevatorId != -1 {
		elevators[newOrder.DesignatedElevatorId].LocalQueue[newOrder.Floor][newOrder.Type] = true

		if newOrder.DesignatedElevatorId != id && newOrder.Type != hardware.BT_Cab {
			hardware.SetButtonLamp(newOrder.Type, newOrder.Floor, true)
		}
	} else {
	}
	informCh.OrderUpdate <- true
}

func DeleteOrder(id int, executedOrder config.ElevatorOrder, elevators *[config.NumberOfElevators]config.ElevatorState) {
	elevators[executedOrder.DesignatedElevatorId].LocalQueue[executedOrder.Floor][executedOrder.Type] = false
	if executedOrder.DesignatedElevatorId == id {
		hardware.SetButtonLamp(executedOrder.Type, executedOrder.Floor, false)
	} else if executedOrder.Type != hardware.BT_Cab {
		hardware.SetButtonLamp(executedOrder.Type, executedOrder.Floor, false)
	}

}

// calculating the costs for each elevator for each order, based on weighting
func costFunction(id int, newOrder config.ElevatorOrder, elevators *[config.NumberOfElevators]config.ElevatorState, availableElevators *[config.NumberOfElevators]bool) (lowestCostId int) {
	if newOrder.Type == hardware.BT_Cab {
		return id
	}
	lowestCost := 1000
	lowestCostId = id

	
	for currentId, elevator := range elevators {
		// if elevator not available, go to next elevator
		if !availableElevators[currentId] || elevator.State == config.Init {
			continue
		}

		// if order already exists on the floor
		if elevator.LocalQueue[newOrder.Floor][newOrder.Type] {
			return -1
		}

		ordersAbove := CheckIfOrdersAbove(elevator)
		ordersBelow := CheckIfOrdersBelow(elevator)
		cost := 0
		distance := newOrder.Floor - elevator.Floor
		
		
		if distance == 0 {
			if elevator.State != config.Moving {
				return currentId
			} else {
				//Elevator moved away from order
				cost += weightDirection
			}
		} else if distance > 0 {
			//Elevator is above order must move up
			if newOrder.Type == hardware.BT_HallDown && ordersAbove {
				cost += weightDirection
			}
			//Check if orders between
			for floor := elevator.Floor + 1; floor < newOrder.Floor; floor++ {
				if elevator.LocalQueue[floor][hardware.BT_HallUp] || elevator.LocalQueue[floor][hardware.BT_Cab] {
					cost += weightOrders
				}
			}
			//Check if right direction, if not add weight
			if elevator.Dir == hardware.MD_Down {
				cost += weightDirection
			} else if elevator.State == config.Moving {
				cost -= weightState
			}

		} else if distance < 0 {
			//Elevator is below order, must move down
			if newOrder.Type == hardware.BT_HallUp && ordersBelow {
				cost += weightDirection
			}
			//Check if orders between
			for floor := newOrder.Floor + 1; floor < elevator.Floor; floor++ {
				if elevator.LocalQueue[floor][hardware.BT_HallDown] || elevator.LocalQueue[floor][hardware.BT_Cab] {
					cost += weightOrders
				}
			}

			distance = -distance
			//Check if right direction, if not add weight
			if elevator.Dir == hardware.MD_Up {
				cost += weightDirection
			} else if elevator.State == config.Moving {
				cost -= weightState
			}
		}
		if elevator.State == config.Executing {
			cost += weightState
			if elevator.Obstructed {
				cost += 5 * weightState
			}

		}

		cost += distance * weightDistance
		//Check if this is the new lowest cost
		if cost < lowestCost {
			lowestCost = cost
			lowestCostId = currentId
		}

	}
	return lowestCostId
}

func CheckIfOrdersAbove(elevatorState config.ElevatorState) (ordersAbove bool) {
	for floor := elevatorState.Floor + 1; floor < config.NumberOfFloors; floor++ {
		for button := hardware.BT_HallUp; button <= hardware.BT_Cab; button++ {
			
			if floor > elevatorState.Floor && elevatorState.LocalQueue[floor][button] {

				return true
			}
		}
	}
	return false
}

func CheckIfOrdersBelow(elevatorState config.ElevatorState) (ordersBelow bool) {
	for floor := 0; floor < elevatorState.Floor; floor++ {
		for button := hardware.BT_HallUp; button <= hardware.BT_Cab; button++ {
			if floor < elevatorState.Floor && elevatorState.LocalQueue[floor][button] {
				return true
			}
		}
	}
	return false
}

func CheckIfOrdersHere(elevatorState config.ElevatorState) (ordersHere bool) {
	for button := hardware.BT_HallUp; button <= hardware.BT_Cab; button++ {
		if elevatorState.LocalQueue[elevatorState.Floor][button] {
			return true
		}
	}
	return false
}

func distributeOrder(id int, order config.ElevatorOrder, elevators *[config.NumberOfElevators]config.ElevatorState, availableElevators *[config.NumberOfElevators]bool, informCh config.InformChannels) {
	lowestCostId := costFunction(id, order, elevators, availableElevators)

	if lowestCostId == id {
		order.DesignatedElevatorId = id
		addOrder(id, order, elevators, informCh)

	} else if lowestCostId != -1 {
		order.DesignatedElevatorId = lowestCostId
		addOrder(id, order, elevators, informCh)
	}
}
