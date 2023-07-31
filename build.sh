#!/usr/bin/env bash
RUN_NAME="topN"

go build -o output/bin/${RUN_NAME}
cd output
scp -i ssh -i /Users/wuhaitao/tecent_cloud/ten_key.pem -r * root@43.138.108.159:/home/whtis/topn


