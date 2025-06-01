# basic-ota-server

- Listens on :5901 and uses /wire/otas as a filepath for the OTAs by default.

- Your file structure should look like:

(in the / directory, on your machine)

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

- `latest` should be a file only containing the version of your latest OTA - like `0.5.0.2`
- Before you start copying in a new OTA to dev and oskr, touch a file called /wire/otas/dnar. After your OTAs are done copying, delete it.
- To build: `go build main.go` in the root of the cloned repo
