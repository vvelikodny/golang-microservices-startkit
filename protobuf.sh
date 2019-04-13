#!/bin/sh

protoc --gofast_out=query-client-service news/news.proto
protoc --gofast_out=storage-service news/news.proto