#!/bin/sh
LOC=$(dirname "$0")
killall main

# Delete old logs.
#find $LOC/../logs -mindepth 1 -type f -mtime +1 -delete
rm -f $LOC/../logs/*

$LOC/main \
				-baud=115200 \
				-log_dir=$LOC/../logs/ \
				-resources=$LOC/../resources \
				-stderrthreshold=info \
				-alsologtostderr=true \
				-logtostderr=false \
				-en_pic \
				-en_compass \
				-v=2 
	#			&
