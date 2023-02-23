#! /bin/bash

queuename=album_queue
hostname=golang_api_elasticmq

part1="http://${hostname}:9324/?Action=CreateQueue&QueueName=${queuename}"
part2="&Attribute.1.Name=VisibilityTimeout&Attribute.1.Value=40"
part3="&Version=2012-11-05&Expires=2023-10-18T22%3A52%3A43PST&AUTHPARAMS"

curl "${part1}${part2}${part3}"
