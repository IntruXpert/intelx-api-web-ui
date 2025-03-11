# 🔍 IntelX API Web UI

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8.svg)](https://golang.org/doc/go1.24)

A modern, user-friendly web interface for the IntelX API built with Go, designed to simplify intelligence gathering and data exploration.

> This project is inspired by the official [Intelligence X SDK](https://github.com/IntelligenceX/SDK).

## 🌟 Features

- 🖥️ Clean, intuitive web interface for IntelX API interactions
- 🔐 Secure API key management
- 🔍 Advanced search capabilities with filters
- 📊 Result visualization and data export options
- ⚡ Fast performance with Go's concurrency features
- 🔄 Real-time result updates
- 📱 Responsive design for desktop and mobile

## 🚀 Installation

### Prerequisites

- Go 1.24 or newer
- IntelX API credentials (obtain at https://intelx.io/account?tab=developer)

### Quick Start

```bash
# Clone the repository
git clone https://github.com/IntruXpert/intelx-api-web-ui.git
cd intelx-api-web-ui

# Build the application
go build

# Run the application
./intelx-api-web-ui
```

## ⚙️ Configuration

There are two ways to configure your IntelX API credentials:

### 1. Using .env file

Create a `.env` file in the root directory with your IntelX API credentials:

```
INTELX_API_KEY=your_api_key_here
PORT=8080
```

### 2. Edit main.go directly

Alternatively, you can directly edit the `DefaultAPIKey` variable in the `main.go` file:

```go
const DefaultAPIKey = "your_api_key_here" // Replace with your API key
```

This method doesn't require a .env file.

## 🔧 Usage

1. Start the application
2. Navigate to `http://localhost:8080` in your browser
3. Enter your search parameters
4. Explore and analyze the results

## 🧩 API Endpoints

The web UI provides access to the following IntelX API endpoints:

- `/search` - Perform searches across various data sources
- `/results` - Retrieve and filter search results
- `/file` - Download specific files from search results
- `/stats` - View usage statistics and limits

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the GNU General Public License v3.0 - see below for details:

```
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
```

For the full license text, please see [https://www.gnu.org/licenses/gpl-3.0.txt](https://www.gnu.org/licenses/gpl-3.0.txt)

## 📞 Contact

- GitHub: [@IntruXpert](https://github.com/IntruXpert)
- X.com: https://x.com/intruxpert

## 🙏 Acknowledgments

- [Intelligence X](https://github.com/IntelligenceX/SDK) for providing the official SDK and API
- All contributors who have helped make this project better 
