# Late

Late is a lightweight utility that facilitates remote command execution on target systems.

It sets up an HTTP server on port 8000, serving a WebSocket connection at the /ws route.

Clients can send the following
JSON data to the WebSocket to execute a command-line ping:

```JSON
{
  "type": "shell",
  "key": "execute",
  "value": "ping google.com"
}
```

Additionally, you can terminate a command session at any time using the following route:

```http request
DELETE http://localhost:8000/kill/{session id}
```

Feel free to customize and enhance this README as needed for your project!