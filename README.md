# Satoshi Radio Pool API (for CK Pool)
The Satoshi Radio Pool API is a Go application designed to manage and provide data on your ckpool instance. It reads pool status data and user-specific logs, storing this information in a PostgreSQL database to support an API for retrieving pool and user data.

## Features

- **API Endpoints**: Access current pool and user statistics via a RESTful API.
- **Data Collection**: Reads from a specified pool status file and user-specific log files every 5 minutes and saves the information to a PostgreSQL database.
- **User and Worker Management**: Track individual user and worker statistics.

## Table of Contents

- [Getting Started](#getting-started)
- [Installation](#installation)
- [Configuration](#configuration)
- [API Documentation](#api-documentation)
- [Contributing](#contributing)

## Getting Started

### Prerequisites

Ensure you have the following installed:

- Go 1.22.2 or higher
- PostgreSQL

### Installation

### Option 1: use the relase binary

Download the latest release from the [releases page](https://github.com/satoshiradio/satoshi-radio-pool-api/releases)

### Option 2: build from source

1. Clone the repository:

```bash
git clone https://github.com/satoshiradio/satoshi-radio-pool-api.git
cd satoshi-radio-pool-api
```

2. Install dependencies:

```bash
go mod download
```

3. Build the application:

```bash
go build -o satoshi-radio-pool-api main.go
```

### Running the Application

Configuring the application is done through environment variables. Update your environment variables to connect to PostgreSQL or set up directly on the system environment if not using a `.env` file.

### Option 1: use an `.env` file

Create a `.env` file in the root directory of the project with the following variables:

```env
POSTGRES_USER=username
POSTGRES_PASSWORD=password
POSTGRES_DB=ckpool
POSTGRES_HOST=localhost
POSTGRES_PORT=5432

POOL_BASE_PATH=/path/to/ckpool
```
you can then start the server by running:

```bash
./satoshi-radio-pool-api
```


### Option 2: set environment variables directly on the system

```bash
export POSTGRES_USER=username
export POSTGRES_PASSWORD=password
export POSTGRES_DB=ckpool
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432

POOL_BASE_PATH=/path/to/ckpool
```
you can then start the server by running:

```bash
./satoshi-radio-pool-api
```

### Option 3: use systemctl

Create a service file in `/etc/systemd/system/satoshi-radio-pool-api.service`:

```ini
[Unit]
Description= Satoshi Radio Pool API
After=network.target

[Service]
# The path to the application binary
ExecStart=/opt/pool/ck-pool-api

# The working directory (where the pool and user directories are located)
WorkingDirectory=/opt/pool

# Restart the service automatically if it crashes
Restart=always
RestartSec=10

# Optional: Set the user and group under which the service runs

# Set environment variables if necessary
Environment=GO_ENV=production
Environment="POSTGRES_USER=user"
Environment="POSTGRES_PASSWORD=password"
Environment="POSTGRES_DB=ckpool_api"
Environment="POSTGRES_HOST=localhost"
Environment="POSTGRES_PORT=5432"
Environment="POOL_BASE_PATH=/path/to/ckpool"

[Install]
# Start the service when the system boots
WantedBy=multi-user.target
```
if you are not using the service, you can start the server with:

```bash
systemctl start satoshi-radio-pool-api
```

## API Documentation

The following are the main API endpoints:

1. Pool Status

GET /api/v1/pool

Fetch the latest mining pool status data.

example response:

```json
{
  "runtime": 123456789,
  "lastupdate": 123456789,
  "users": 50,
  "workers": 150,
  "hashrate1m": "100 GH/s",
  "hashrate5m": "120 GH/s",
  "hashrate15m": "115 GH/s",
  "hashrate1hr": "110 GH/s",
  "hashrate6hr": "105 GH/s",
  "hashrate1d": "102 GH/s",
  "hashrate7d": "98 GH/s",
  "diff": 1.0,
  "accepted": 2000000,
  "rejected": 1000,
  "bestshare": 900000,
  "sps1m": 0.5,
  "sps5m": 0.7,
  "sps15m": 0.6,
  "sps1h": 0.65,
}
```

2. Pool Hashrates

GET /api/v1/pool/hashrates

Retrieve the current hashrate statistics of the pool.

Example response:

```json
[
  {
    "hashrate1m": "100 GH/s",
    "hashrate5m": "120 GH/s",
    "hashrate15m": "115 GH/s",
    "hashrate1hr": "110 GH/s",
    "hashrate6hr": "105 GH/s",
    "hashrate1d": "102 GH/s",
    "hashrate7d": "98 GH/s",
    "saved_at": 2024-10-30T10:06:32.345409Z"
  },
  .....
]
````

3. User Data

GET /api/v1/users/{username}

Get detailed statistics for a specified user.

Example response:

```json
{
  "hashrate1m": "100 GH/s",
  "hashrate5m": "120 GH/s",
  "hashrate1hr": "115 GH/s",
  "hashrate1d": "110 GH/s",
  "hashrate7d": "105 GH/s",
  "lastshare": 123456789,
  "workers": 1,
  "shares": 1000000,
  "bestshare": 1000.3234234,
  "bestever": 1000,
  "authorized": 123456789,
  "worker": [
    {
      "workername": "bc1........worker1",
      "hashrate1m": "100 GH/s",
      "hashrate5m": "120 GH/s",
      "hashrate1hr": "115 GH/s",
      "hashrate1d": "110 GH/s",
      "hashrate7d": "105 GH/s",
      "lastshare": 123456789,
      "shares": 1000000,
      "bestshare": 1000.3234234,
      "bestever": 1000,
    },
    .....
  ]


}
```

4. User Hashrates

GET /api/v1/users/{username}/hashrates

Retrieve hashrate data for a specific user.

Example data:

```json
[
  {
    "hashrate1m": "100 GH/s",
    "hashrate5m": "120 GH/s",
    "hashrate15m": "115 GH/s",
    "hashrate1hr": "110 GH/s",
    "hashrate6hr": "105 GH/s",
    "hashrate1d": "102 GH/s",
    "hashrate7d": "98 GH/s",
    "saved_at": 2024-10-30T10:06:32.345409Z
  },
  .....
]
```

5. User Worker Data

GET /api/v1/users/{username}/workers/{workername}/hashrates

Retrieve detailed hashrate data for a specific worker under a user.

Example data:

  ```json
  {
    "hashrate1m": "100 GH/s",
    "hashrate5m": "120 GH/s",
    "hashrate15m": "115 GH/s",
    "hashrate1hr": "110 GH/s",
    "hashrate6hr": "105 GH/s",
    "hashrate1d": "102 GH/s",
    "hashrate7d": "98 GH/s",
    "saved_at": 2024-10-30T10:06:32.345409Z
  }
  ```



## CORS

CORS is enabled on all routes, allowing for requests from any origin.

Running the Application


## Contributing

Contributions are welcome! Fork the repo, create a new branch for your feature or bugfix, and submit a pull request.

## License

GNU Public license V3. See included LICENCE for details.
