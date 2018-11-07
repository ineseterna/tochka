# Tochka Free Market

Tochka Free Market is free Dark Net Maketplace (DNM) software for buildign trading communities.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

* golang
* git
* Postgres 
* tor / torsocks

### Installing

To get Tochka running:

```
# 1. Get Tochka source code
torsocks go get -insecure qxklmrhx7qkzais6.onion/Tochka/tochka-free-market
# 2. Build Tochka from source
cd $GOPATH/src/qxklmrhx7qkzais6.onion/Tochka/tochka-free-market
go build
# 3. Sync DB models and supplementary data
su postgres
    createdb go_t
    psql go_t < dumps/cities.sql
    psql go_t < dumps/countries.sql 
exit
/tochka-free-market sync
# 4. Edit settings
cp settings.json.example settings.json
vim settings.json
# 5. Run HTTP server
./tochka-free-market

```

Go to http://localhost:8081/ and register a new user. Add admin privelegies to new account:

```
./tochka-free-market user <username> grant admin
```

## License
 
The MIT License (MIT)

Copyright (c) 2015 Chris Kibble

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
