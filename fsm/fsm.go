package fsm

import (
	"project-group11/config"
	"project-group11/hardware"
	"project-group11/orderHandler"
	"time"
)

type State string

func Fsm(id int, hardwareCh config.HardwareChannels, netCh config.NetworkChannels, elevators *[config.NumberOfElevators]config.ElevatorState, availableElevators *[config.NumberOfElevators]bool, informCh config.InformChannels, timerCh config.TimerChannels) {
	elevators[id].Id = id
	var newDirection hardware.MotorDirection
	var timestampMechanicalErrorCheck time.Time

	newDirection = hardware.MD_Down

	drv_stop := make(chan bool)

	go hardware.PollStopButton(drv_stop)


	hardware.SetMotorDirection(newDirection)

	for floor := 0; floor < config.NumberOfFloors; floor++ {

		hardware.SetButtonLamp(hardware.BT_HallUp, floor, false)
		hardware.SetButtonLamp(hardware.BT_HallDown, floor, false)
		hardware.SetButtonLamp(hardware.BT_Cab, floor, false)

		hardware.SetDoorOpenLamp(false)
	}

	for {
		switch elevators[id].State {
		case config.Init:
			elevators[id].Floor = <-hardwareCh.DrvFloors
			hardware.SetFloorIndicator(elevators[id].Floor)

			newDirection = hardware.MD_Stop
			elevators[id].Dir = newDirection
			hardware.SetMotorDirection(newDirection)
			elevators[id].State = config.Idle

			informCh.SendState <- elevators[id]

		case config.Idle:
			select {
			case <-informCh.OrderUpdate:
				newDirection = determineDirection(&elevators[id])
				for button := hardware.BT_HallUp; button <= hardware.BT_Cab; button++ {
					if elevators[id].LocalQueue[elevators[id].Floor][button] {
						goToExecuting(id, elevators, timerCh, informCh)
						break
					}
				}
				informCh.SendState <- elevators[id]

			default:
				if newDirection != hardware.MD_Stop {
					hardware.SetMotorDirection(newDirection)
					if elevators[id].Dir != newDirection {
						elevators[id].Dir = newDirection
					}
					elevators[id].State = config.Moving
					timestampMechanicalErrorCheck = time.Now()

					informCh.SendState <- elevators[id]
				}
			}

		case config.Moving:
			select {
			case <-informCh.OrderUpdate:
				newDirection = determineDirection(&elevators[id])
				informCh.SendState <- elevators[id]

			case elevators[id].Floor = <-hardwareCh.DrvFloors:
				timestampMechanicalErrorCheck = time.Now()

				hardware.SetFloorIndicator(elevators[id].Floor)

				ordersAbove := orderHandler.CheckIfOrdersAbove(elevators[id])
				ordersBelow := orderHandler.CheckIfOrdersBelow(elevators[id])

				orderCab := elevators[id].LocalQueue[elevators[id].Floor][hardware.BT_Cab]
				orderGoingUp := elevators[id].LocalQueue[elevators[id].Floor][hardware.BT_HallUp] && (elevators[id].Dir == hardware.MD_Up || !ordersBelow || elevators[id].Floor == 0)
				orderGoingDown := elevators[id].LocalQueue[elevators[id].Floor][hardware.BT_HallDown] && (elevators[id].Dir == hardware.MD_Down || !ordersAbove || elevators[id].Floor == config.NumberOfFloors-1)

				if orderCab || orderGoingUp || orderGoingDown {
					goToExecuting(id, elevators, timerCh, informCh)
					break
				}
				
			default:
				if newDirection == hardware.MD_Stop {
					elevators[id].Dir = newDirection
					hardware.SetMotorDirection(newDirection)
					elevators[id].State = config.Idle
					informCh.SendState <- elevators[id]
				}

				if time.Since(timestampMechanicalErrorCheck) > config.TimeoutMechanicalError {
					elevators[id].State = config.SystemFailure
					availableElevators[id] = false

					informCh.SendState <- elevators[id]
					informCh.OrderRedistribute <- id
				}
			}

		case config.Executing:

			//if order here
			for button := hardware.BT_HallUp; button <= hardware.BT_Cab; button++ {
				if elevators[id].LocalQueue[elevators[id].Floor][button] {
					executedOrder := config.ElevatorOrder{elevators[id].Floor, button, id, true}
					informCh.SendOrder <- executedOrder
					orderHandler.DeleteOrder(id, executedOrder, elevators)
				}
			}


			select {

			case <-informCh.OrderUpdate:

				for button := hardware.BT_HallUp; button <= hardware.BT_Cab; button++ {
					if elevators[id].LocalQueue[elevators[id].Floor][button] {
						if !elevators[id].Obstructed {
							timerCh.TimerDoorStart <- true
						}
						break
					}
				}
				informCh.SendState <- elevators[id]
			//implement for obstruction?
			case <-timerCh.TimerDoorFinished:

				hardware.SetDoorOpenLamp(false)

				newDirection = determineDirection(&elevators[id])
				if newDirection != hardware.MD_Stop {
					hardware.SetMotorDirection(newDirection)
					if elevators[id].Dir != newDirection {
						elevators[id].Dir = newDirection
					}
					elevators[id].State = config.Moving
					timestampMechanicalErrorCheck = time.Now()
				} else {
					elevators[id].Dir = newDirection
					elevators[id].State = config.Idle
				}
				informCh.SendState <- elevators[id]
			}
			
		case config.SystemFailure:
			select {
			case <- informCh.OrderUpdate:
			case elevators[id].Floor = <-hardwareCh.DrvFloors:
				timestampMechanicalErrorCheck = time.Now()
				availableElevators[id] = true

				hardware.SetFloorIndicator(elevators[id].Floor)
				for button := hardware.BT_HallUp; button <= hardware.BT_Cab; button++ {
					if elevators[id].LocalQueue[elevators[id].Floor][button] {
						goToExecuting(id, elevators, timerCh, informCh)
						break
					}
				}

				newDirection = determineDirection(&elevators[id])
				if newDirection != hardware.MD_Stop {
					hardware.SetMotorDirection(newDirection)
					if elevators[id].Dir != newDirection {
						elevators[id].Dir = newDirection
					}
					elevators[id].State = config.Moving
					timestampMechanicalErrorCheck = time.Now()
				} else {
					elevators[id].Dir = newDirection
					elevators[id].State = config.Idle
				}

				informCh.SendState <- elevators[id]
			}
		}
	}
}