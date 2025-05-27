#!/usr/bin/env bash
# Author: Andy Gorman

set -e

if tmux has-session -t=gorman-zone 2> /dev/null; then
	tmux attach -t gorman-zone
	exit
fi

tmux new-session -d -s gorman-zone -n frontend -x "$(tput cols)" -y "$(tput lines)" -c "frontend/"

tmux split-window -t gorman-zone:frontend -h -c "#{pane_current_path}" -p 40
tmux split-window -t gorman-zone:frontend.right -v -c "#{pane_current_path}"

tmux send-keys -t gorman-zone:frontend.left "vim" Enter
tmux send-keys -t gorman-zone:frontend.top-right "git s" Enter

tmux new-window -d -t gorman-zone -n api -c "api/"

tmux split-window -t gorman-zone:api -h -p 40 -c "api/"
tmux split-window -t gorman-zone:api.right -v -c "api/"

tmux send-keys -t gorman-zone:api.left "vim" Enter
tmux send-keys -t gorman-zone:api.top-right "git s" Enter



tmux attach -t gorman-zone:frontend.left
