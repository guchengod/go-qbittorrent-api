# qBittorrent API Client in Go

[![Go Version](https://img.shields.io/badge/go-1.23.2+-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](https://opensource.org/licenses/MIT)


A comprehensive Go client for interacting with the qBittorrent Web API. This library provides a simple and efficient way to manage your qBittorrent instance programmatically, supporting all major API endpoints.

## Features

- **Full API Coverage**: Supports all qBittorrent Web API endpoints, including torrent management, application preferences, RSS, and search.
- **Easy to Use**: Simple and intuitive API design for seamless integration into your Go projects.
- **Authentication**: Handles login and session management automatically.
- **Error Handling**: Robust error handling and status code checking for reliable operation.
- **Extensible**: Easily extendable to support future qBittorrent API updates.

## Installation

To install the package, use `go get`:

```bash
go get github.com/guchengod/qbittorrent-api-go
```

## Quick Start

Here's a quick example to get you started:

```go
package main

import (
	"fmt"
	"github.com/guchengod/qbittorrent-api-go/qbittorrent"
)

func main() {
	client := qbittorrent.NewDefaultClient("http://localhost:8080")
	err := client.Login("admin", "adminadmin")
	if err != nil {
		fmt.Println("Login failed:", err)
		return
	}
	defer client.Logout()

	version, err := client.GetApplicationVersion()
	if err != nil {
		fmt.Println("Failed to get application version:", err)
		return
	}
	fmt.Println("Application version:", version)
}
```

## Supported Endpoints

- **Authentication**: Login, Logout
- **Application**: Get version, Shutdown, Get/set preferences
- **Torrent Management**: Add, pause, resume, delete, recheck, reannounce torrents
- **RSS**: Add feeds, manage items, set auto-downloading rules
- **Search**: Start, stop, and manage searches
- **Sync**: Get main data, torrent peers data
- **Transfer Info**: Get global transfer info, set speed limits
- **Logs**: Get main log, peer log

## Documentation

For detailed documentation, refer to the [qBittorrent Web API Documentation](https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-(qBittorrent-4.1)).

## Contributing

Contributions are welcome! If you'd like to contribute, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Commit your changes and push to your fork.
4. Submit a pull request with a detailed description of your changes.

## Star History

<a href="https://star-history.com/#guchengod/go-qbittorrent-api&Date">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=guchengod/go-qbittorrent-api&type=Date&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=guchengod/go-qbittorrent-api&type=Date" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=guchengod/go-qbittorrent-api&type=Date" />
 </picture>
</a>

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---
