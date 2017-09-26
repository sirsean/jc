# jc (JSON Client)

A simple and fast command-line client for JSON APIs.

It saves your requests (and responses) to the filesystem so you can run
them again in the future.

It can handle GET and POST requests, with a JSON body.

## Motivation

When I develop APIs, I typically play around with them with an HTTP client
to make sure they're working the way I want. My choice for a long time has
been a Chrome extension called Postman, which is feature-rich and pretty
nice to use. I like the way it remembers what requests I've run in the
past and lets me run them again. But the more things it remembers, the
slower it gets. I've had enough waiting for its UI to be responsive.

## Installation

```
go get github.com/sirsean/jc
cd jc
go install
```

It will store its data in a directory called `.jc` in whatever directory
you run it in. This way, you have different sets of requests on a per-project
basis.

## Usage

The following commands are supported:

- `ls`: List the available requests
- `new <id>`: Create a new (empty) request
- `del <id>`: Delete an existing request
- `run <id>`: Execute an existing request
- `resp <id>`: View the result of a previously-exected request

The `id` is an arbitrary string (no spaces) that you can use to identify
your requests.

The `new` command generates a new JSON file representing your request. You
must edit this file to describe your request.

```
$ jc new get-things
.jc/get-things/request.json
```

That file contains the empty request:

```
{
    "id": "get-things",
    "url": "",
    "method": "",
    "basic_auth": {
        "username": "",
        "password": ""
    },
    "client_cert": {
        "ca_cert": "",
        "client_cert": "",
        "client_key": ""
    },
    "headers": null
}
```

Edit it to describe the request you want to make:

```
{
    "id": "get-things",
    "url": "http://localhost:8080/things",
    "method": "GET",
    "basic_auth": {
        "username": "",
        "password": ""
    },
    "client_cert": {
        "ca_cert": "",
        "client_cert": "",
        "client_key": ""
    },
    "headers": null
}
```

(Note that if you don't enter in your HTTP Basic Auth or client cert or headers,
they will be ignored.)

Now you can see your request in the list:

```
$ jc ls
get-things: GET http://localhost:8080/things
```

Now that it's filled out, you can execute your request:

```
$ jc run get-things
120ms
.jc/get-things/response.json
```

That file contains the response body of the request. Also displayed is the
duration of the request, which can be useful information.

You can view the response:

```
$ jc resp get-things | less
```

You can also do a POST request with a JSON body.

```
$ jc new add-thing
.jc/add-thing/request.json
```

Edit that file to describe a POST request:

```
{
    "id": "add-thing",
    "url": "http://localhost:8080/things",
    "method": "POST",
    "basic_auth": {
        "username": "",
        "password": ""
    },
    "headers": {
        "Content-Type": "application/json"
    }
}
```

But then also add a new file, called `.jc/add-thing/body.json`:

```
{
    "name": "Thing One",
    "description": "This is the first thing."
}
```

Now when you run your new request, it will pick up that `body.json` and send
it to the given endpoint. It will populate the `.jc/add-thing/response.json`
file with the result.
