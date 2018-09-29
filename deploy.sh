#!/usr/bin/env bash
cd ./api-aws/
if [ $1 = "production" ]; then
    sls deploy --stage prod --conceal
else
    sls deploy --conceal
fi