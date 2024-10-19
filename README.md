# GHBypass

Very basic HTTP bypass proxy for localhost projects

### Server launch
You must have server with wildcard DNS A entry (for example - *example.org*)

```
make build-server
make run-server --host=<example.org>
```
You can go to any subdomain on your server (for example *test.example.org*), and get client binaries from your server

### Client launch
For example:  
- I have public server on IP: *200.1.1.1* (with wildcard dns **.example.org*)
- I run backend locally with access by *local.dev:8052*
- I want get access to my backend on *test.example.org:8080*
```
.\client-windows.exe --subdomain=test --proxy=200.1.1.1:8080 --expose=local.dev:8052
```
