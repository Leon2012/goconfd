#!/bin/bash
set -e
rm -rf ./nohup.out
nohup ./build/monitor --config ./build/monitor.toml &
nohup ./build/agent --config ./build/agent.toml &
nohup ./build/dashboard --config ./build/dashboard.toml &