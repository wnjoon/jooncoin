# jooncoin
Practice making coin based blockchain using golang learned by NomadCoder

## Setting Environment

- go.mod : Pretty simillar to package.json in Javascript
	- go mod init github.com/wnjoon/jooncoin

## Dependencies

1. [gorilla/mux](https://github.com/gorilla/mux)
- Package gorilla/mux implements a request router and dispatcher for matching incoming requests to their respective handler.
- <u>To Use Pattern of input parameter from http request (2 line)</u>
```
r := mux.NewRouter()
r.HandleFunc("/products/{key}", ProductHandler)
r.HandleFunc("/articles/{category}/", ArticlesCategoryHandler)
r.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)
```
- go get -u github.com/gorilla/mux

2. [boltdb/bolt](https://github.com/boltdb/bolt)
- An embedded key/value database for Go.
- This repository has been archived by the owner. It is now read-only.
	- <u>Never be modified since it is completey developed == Stable </u>
- <i>In this program, structure would be "Hash/Block{data, hash, prevHash ....}"</i>
- In bolt, <u>Tables are called "Buckets"</u>
- go get github.com/boltdb/bolt/...

3. [evnix/boltdbweb](https://github.com/evnix/boltdbweb)
- A simple web based boltdb GUI Admin panel.
- Available to view and <u>change</u> data in database
- **Usage** : boltdbweb --db-name=<DBfilename>[required] --port=<port>[optional] --static-path=<static-path>[optional]
- go get github.com/evnix/boltdbweb

4. [br0xen/boltbrowser](https://github.com/br0xen/boltbrowser)
- A CLI Browser for BoltDB Files
- Only can view data in database
- **Usage** : boltbrowser <filename>
- go get github.com/br0xen/boltbrowser