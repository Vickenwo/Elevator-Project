#!/bin/bash

echo off

#gnome-terminal --title="Elevator Server" -- ElevatorServer


gnome-terminal --geometry 65x20+0+0 --title="Simulator 1" -- ./Simulator/SimElevatorServer --port 12000
gnome-terminal --geometry 65x20+700+0 --title="Simulator 2" -- ./Simulator/SimElevatorServer --port 12001
gnome-terminal --geometry 65x20+1500+0 --title="Simulator 3" -- ./Simulator/SimElevatorServer --port 12002
#gnome-terminal --geometry 65x20+2200+0 --title="Simulator 4" -- ./Simulator/SimElevatorServer --port 12003

gnome-terminal --geometry 65x20+0+450 --title="Elevator 1" -- go run main.go -id=0 -localhostId="12000"
gnome-terminal --geometry 65x20+700+450 --title="Elevator 2" -- go run main.go -id=1 -localhostId="12001"
gnome-terminal --geometry 65x20+1500+450 --title="Elevator 3" -- go run main.go -id=2 -localhostId="12002"
#gnome-terminal --geometry 65x20+2200+450 --title="Elevator 4" -- go run main.go -id=3 -localhostId="12003"