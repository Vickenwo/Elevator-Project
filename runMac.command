
osascript -e 'tell app "Terminal"
    do script "rdmd ~/Documents/TTK4145sanntid/Simulator-v2/src/sim_server.d --port 12000"
end tell'

osascript -e 'tell app "Terminal"
    do script "go run Documents/TTK4145sanntid/project-group11/main.go 12000"
end tell'

sleep .5

osascript -e 'tell app "Terminal"
    do script "rdmd ~/Documents/TTK4145sanntid/Simulator-v2/src/sim_server.d --port 12001"
end tell'

osascript -e 'tell app "Terminal"
    do script "go run Documents/TTK4145sanntid/project-group11/main.go 12001"
end tell'

sleep .5

osascript -e 'tell app "Terminal"
    do script "rdmd ~/Documents/TTK4145sanntid/Simulator-v2/src/sim_server.d --port 12002"
end tell'

osascript -e 'tell app "Terminal"
    do script "go run Documents/TTK4145sanntid/project-group11/main.go 12002"
end tell'
