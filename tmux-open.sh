#!/bin/sh

session="frr"

tmux kill-session -t $session
tmux new-session -d -s $session

window=1
tmux new-window -t $session:$window -n 'main'
tmux send-keys -t $session:$window "docker exec -it --privileged main $1" C-m

window=2
tmux new-window -t $session:$window -n 'start'
tmux send-keys -t $session:$window "docker exec -it --privileged start $1" C-m

window=3
tmux new-window -t $session:$window -n 'top'
tmux send-keys -t $session:$window "docker exec -it --privileged top $1" C-m

window=4
tmux new-window -t $session:$window -n 'end'
tmux send-keys -t $session:$window "docker exec -it --privileged end $1" C-m

window=5
tmux new-window -t $session:$window -n 'bottom'
tmux send-keys -t $session:$window "docker exec -it --privileged bottom $1" C-m

window=6
tmux new-window -t $session:$window -n 'main-peer'
tmux send-keys -t $session:$window "docker exec -it --privileged main-peer $1" C-m

tmux kill-window -t $session:0
tmux select-window -t $session:1

tmux attach-session -t $session
