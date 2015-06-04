# jsonconsul

`jsonconsul` allows the creation of json files from the values that
are located in Consul. It also contains the ability to import KV keys
into Consul.

## Usage

### Import

A json config can be used to update key values in Consul.

Here is an example file named `example.json`
```
{"foo":{"bar":"test","blah":"Test","do":"TEST","loud":{"asd":{"bah":"test"}}}}
```

Now if we want to import this into the root prefix then we'd do the following:

```sh
jsonconsul import -json-file example.json
```

To import into an alternate prefix the following needs to be done:

```sh
jsonconsul import -prefix='vsco/buzz' -json-file example.json
```

### Export

There are five ways to run `jsonconsul export`. There are:

 - Output to STDOUT
 - Output to file
 - Output to file but timestamp the output file and symlink to name.
 - Poll and output to file after a duration.
 - Poll and output to file after a duration but with timestamped output file.

#### Output to STDOUT
```sh
jsonconsul export -prefix="foo"
```

#### Output to file
```sh
jsonconsul export -prefix="foo" -config=foo.json
```

#### Output to file with timestamp
```sh
jsonconsul export -prefix="foo" -config=foo.json -timestamp
```

This generates a file called `foo.json.<unixtimestamp>`. `foo.json`
will then be a symbolic link to `foo.json.<unixtimestamp>`.


#### Poll and output to file
```sh
jsonconsul watch -prefix="foo" -config=foo.json -poll
```

This polls consul every minute for changes and outputs those values to
json. If an alternate frequency is preferred then include the
`-poll_frequency` flag.

#### Poll and output to file with timestamp
```sh
jsonconsul watch -prefix="foo" -config=foo.json -timestamp -poll
```
