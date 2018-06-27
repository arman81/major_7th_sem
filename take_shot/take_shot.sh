#!/bin/bash
#author anil.khadwal@gmail.com

#command to take screen_shot using s/w scrot (previously installed)

scrot  '%Y:%m:%d:%H:%M:%S.png' -e 'mv $f ./server_copies'

scrot 'snapshot.png' -e 'mv $f ./server_copies'
