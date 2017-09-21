#!/bin/sh
LOC=$(dirname "$0")
killall main

# Delete old logs.
find $LOC/../logs -mindepth 1 -type f -mtime +10 -delete

$LOC/main \
				-log_dir=$LOC/logs/ \
				-resources=$LOC/resources \
				-stderrthreshold=info \
				-alsologtostderr=true \
				-logtostderr=false \
				-en_pic \
				-en_roomba \
				-en_io \
				-en_compass=false \
				-v=2 
	#			&
