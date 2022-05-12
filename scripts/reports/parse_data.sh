#!/bin/bash

# table
#grep TPS $1 | awk -F '[ :,]+' '{print "<tr><td> "$2" </td><td> 传输 </td><td> "$4" </td><td> "$6" </td><td> "$8" </td></tr>"}'

# csv
grep TPS $1 | awk -F '[ :,]+' '{split($8,a,"m");split($10,b,"m");print "["$4","substr($15,3)","substr($11,4)","$6","a[1]","b[1]}'
