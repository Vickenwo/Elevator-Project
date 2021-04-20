package fsm

import (
	"project-group11/config"
	"project-group11/hardware"
	"project-group11/orderHandler"
)

func determineDirection(elevatorState *config.ElevatorState) (newDirection hardware.MotorDirection) {
	ordersAbove, ordersBelow, ordersHere := orderHandler.CheckIfOrders(elevatorState)
	switch elevatorState.State {
	case config.Idle:
		if ordersHere {
			return hardware.MD_Stop 
		} else if ordersAbove && elevatorState.Dir != hardware.MD_Down {
			return hardware.MD_Up 
		} else if ordersBelow && elevatorState.Dir != hardware.MD_Up {
			return hardware.MD_Down 
		}

	case config.Moving:
		if !ordersAbove && !ordersBelow {
			return hardware.MD_Stop 
		} else if (ordersAbove && elevatorState.Dir == hardware.MD_Up) || (ordersBelow && elevatorState.Dir == hardware.MD_Down) {
			return elevatorState.Dir
		}

	case config.Executing:
		if ordersHere || (!ordersAbove && !ordersBelow) {
			return hardware.MD_Stop 
		} else if ordersAbove && !ordersBelow {
			return hardware.MD_Up
		} else if ordersBelow && !ordersAbove {
			return hardware.MD_Down 
		} else if elevatorState.Dir == hardware.MD_Stop {
			return hardware.MD_Down
		}
	}
	
	return elevatorState.Dir 
}

func CheckObstruction(id int, elevators *[config.NumberOfElevators]config.ElevatorState, informCh config.InformChannels, hardwareCh config.HardwareChannels, timerCh config.TimerChannels) {
	for {
		select {
		case obstruction := <-hardwareCh.DrvObstr:
			elevators[id].Obstructed = obstruction
			if elevators[id].State == config.Executing {
				if obstruction {
					timerCh.TimerDoorStop <- true
				} else {
					timerCh.TimerDoorStart <- true
				}
			}
			informCh.SendState <- elevators[id]
		}
	}
}

func goToExecuting(id int, elevators *[config.NumberOfElevators]config.ElevatorState, timerCh config.TimerChannels, informCh config.InformChannels) {
	hardware.SetMotorDirection(hardware.MD_Stop)
	hardware.SetDoorOpenLamp(true)

	if !elevators[id].Obstructed {
		timerCh.TimerDoorStart <- true
	}

	elevators[id].State = config.Executing
	informCh.SendState <- elevators[id]
}
