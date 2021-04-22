Elevator Project
================


Summary
-------
In this project we have created software for controlling `3` elevators working in parallel across `4` floors, written in Golang. Furthermore, we use UDP broadcasting in order to to make the elevators communicate over network. In this way we can use information from all the elevators in order to calculate the best suited elevator to handle each order. 


How to run the program
-----------------
In order to run the program, you can simply run [runLinux.sh](runLinux.sh) in the terminal. This will simulate 3 elevators with id 0,1 and 2, respectively. These elevators will connect to ports `12000 + id`. 

In addition, the [elevator simulator](Simulator/SimElevatorServer) must be downloaded. Remember to `chmod +x SimElevatorServer` in order to give yourself permission to run downloaded files.

Modules
----------------------
### Config
The config module contains the configurational settings for the system, meaning that it defines constants and types that will be used in the rest of the modules. This is the cornerstone of the project.

### FSM
The Finite State Machine module concerns the different states in the system, `Init`, `Idle`, `Moving`, `Executing` and `SystemFailure`. [fsmFunctions.go](fsm/fsmFunctions.go) contains functions and goroutines that deals with the states, such as determining the direction and checking for obstruction.

### Hardware
Contains delivered code related to hardware in the system. 

### Network
This module contains delivered code with functionality for connecting, broadcasting and receiving messages over the network. In addition, we made functions for monitoring peers and sending messages over the network based on informchannels in [network.go](network/network.go). The network communication is based on UDP, with peer to peer communication.


### Order Handler
The order handler module operates based on actions polled from different channels. We have implemented functions for adding and deleting orders, checking for orders under, above or at the current floor and distributing orders. We have also implemented a cost function that calculates a cost for each elevator based on weights related to different conditions. In this way we can find the most suited elevator for each order. 


### Util
The util module contains a timer template that we use in our program. 

Fault tolerance
---------------------
In this project we have implemented fault tolerance for several scenarios:
- **Network loss**: If one elevator loses connection to the network, it will work as a single elevator. This means that it will complete all of the orders in the order list at the time it disconnected. The other elevators will also complete these orders, redistributing the lost elevator’s orders. In addition, the cab orders at the time will be saved in the order list until the disconnected elevator eventually reconnects. 
- **Power loss**: If power loss occurs, the hall orders of the lost elevator will be redistributed to the remaining available elevators in the system. The system will now work as a system with `n-1` elevators.

