# jsonconsul

A tool to make config files from Key Values in Consul.

## Description

`jsonconsul` allows the creation of json files from the values that
are located in Consul. This allows applications that are configured
via JSON to be able to read the values.

## Usage

There are five ways to run `jsonconsul`. There are:

 - Output to STDOUT
 - Output to file
 - Output to file but timestamp the output file and symlink to name.
 - Poll and output to file after a duration.
 - Poll and output to file after a duration but with timestamped output file.

### Output to STDOUT
```sh
jsonconsul -prefix="foo"
```

### Output to file
```sh
jsonconsul -prefix="foo" -config=foo.json
```

### Output to file with timestamp
```sh
jsonconsul -prefix="foo" -config=foo.json -timestamp
```

This generates a file called `foo.json.<unixtimestamp>`. `foo.json`
will then be a symbolic link to `foo.json.<unixtimestamp>`.


### Poll and output to file
```sh
jsonconsul -prefix="foo" -config=foo.json -poll
```

This polls consul every minute for changes and outputs those values to
json. If an alternate frequency is preferred then include the
`-poll_frequency` flag.

### Poll and output to file with timestamp
```sh
jsonconsul -prefix="foo" -config=foo.json -timestamp -poll
```
