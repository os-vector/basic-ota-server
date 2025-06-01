# basic-ota-server

- Listens on :5901 and uses /wire/otas as a filepath for the OTAs by default.

- Your file structure should look like:

(in / on your machine)

```
/wire
  /otas
    latest
    /full
      /dev
         #.#.#.#.ota
	 ...
      /oskr
         #.#.#.#.ota
         ...
```

- `latest` should be a file only containing the version of your latest OTA - like `3.0.1.1`
- before you start copying in a new OTA to dev and oskr, touch a file called /wire/otas/dnar. after your OTAs are done, delete it
- To build: `go build main.go` in the root of the cloned repo
