[![serverless](http://public.serverless.com/badges/v3.svg)](http://www.serverless.com)
[![Build Status](https://travis-ci.org/RadioCheckerApp/api.svg?branch=master)](https://travis-ci.org/RadioCheckerApp/api)

# 🔌 RadioChecker.com API Services
The RC API Services provide and maintain API endpoints used by
RadioChecker.com apps, i.e. the web frontend. Built upon the
[AWS Lambda](https://aws.amazon.com/lambda) stack, the API leverages a
serverless architecture developed with the
[Serverless Framework](https://serverless.com).

## Endpoints
- `GET /meta`
- `GET /stations`
- `GET /stations/{station}/tracks?date=2018-02-12&filter=top`
- `GET /stations/{station}/tracks?week=2018-02-12&filter=all`
- `GET /tracks/search?date=2018-02-12&q=Dani+California`
- `GET /tracks/search?week=2018-02-12&q=The+Adventures+Of+Rain+Dance+Maggie`

- `POST /tracks`