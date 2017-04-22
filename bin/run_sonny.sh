#!/bin/sh
LOC=$(dirname "$0")
killall main

$LOC/main \
				-baud=19200 \
				-log_dir=$LOC/../logs/ \
				-resources=$LOC/../resources \
				-stderrthreshold=info \
				-alsologtostderr=true \
				-logtostderr=false \
				-v=2 \
				-en_pic
	#			&
