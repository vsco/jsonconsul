# jsonconsul

[![Build Status](https://travis-ci.org/vsco/jsonconsul.svg)](https://travis-ci.org/vsco/jsonconsul)

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
jsonconsul import example.json
```

To import into an alternate prefix the following needs to be done:

```sh
jsonconsul import -prefix='vsco/buzz' example.json
```

### Export

There are five ways to run `jsonconsul export`. There are:

 - Output to STDOUT
 - Output to file
 - Output to file but timestamp the output file and symlink to name.
 - Poll and output to file after a duration.
 - Poll and output to file after a duration but with timestamped output file.

Options:

 - `-json-values` Convert the values from Consul and treat them as JSON values.

#### Output to STDOUT
```sh
jsonconsul export -prefix="foo"
```

If we don't want to include the prefix in the outputed json:
```sh
jsonconsul export -include-prefix=false -prefix="foo"
```

#### Output to file
```sh
jsonconsul export -prefix="foo" foo.json
```

#### Output to file with timestamp
```sh
jsonconsul export -prefix="foo" -timestamp foo.json
```

This generates a file called `foo.json.<unixtimestamp>`. `foo.json`
will then be a symbolic link to `foo.json.<unixtimestamp>`.


#### Watch and output to file
```sh
jsonconsul watch -prefix="foo" foo.json
```

This polls consul every minute for changes and outputs those values to
json. If an alternate frequency is preferred then include the
`-watch-frequency` flag.

#### Watch and output to file with timestamp
```sh
jsonconsul watch -prefix="foo" -timestamp foo.json
```

## License

The MIT License (MIT)

Copyright (c) 2015 Visual Supply, Co.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
