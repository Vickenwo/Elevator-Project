echo off

::start "Elevator server" cmd /k "ElevatorServer"
::mode con: cols=40 lines=100

start "Simulator 1" cmd /k ".\Simulator-v2-master\Simulator-v2-master\SimElevatorServer.exe --port 12345"
::start "Simulator 1" cmd /k "C:\Users\mmats\OneDrive\Skrivebord\Sanntid\project-group11\Simulator-v2-master\Simulator-v2-master\SimElevatorServer.exe --port 12345"
start "Simulator 2" cmd /k "C:\Users\mmats\OneDrive\Skrivebord\Sanntid\project-group11\Simulator-v2-master\Simulator-v2-master\SimElevatorServer.exe --port 12346"
start "Simulator 3" cmd /k "C:\Users\mmats\OneDrive\Skrivebord\Sanntid\project-group11\Simulator-v2-master\Simulator-v2-master\SimElevatorServer.exe --port 12347"

start "Elevator 1" cmd /k "go run main.go -id=^"0^" -localhostId=^"12345^""
start "Elevator 2" cmd /k "go run main.go -id=^"1^" -localhostId=^"12346^""
start "Elevator 3" cmd /k "go run main.go -id=^"2^" -localhostId=^"12347^""

