# win-dex-agent

InsightFinder Windows Agent for Digital Employee Experience (DEX) monitoring. This agent collects comprehensive system metrics from Windows endpoints and sends them to the InsightFinder platform for analysis and anomaly detection.

## Overview

The win-dex-agent is a lightweight monitoring agent designed to collect real-time performance metrics from Windows systems. It leverages both native Go libraries and Windows Performance Data Helper (PDH) APIs to gather detailed system telemetry, enabling IT teams to monitor and improve digital employee experience.

## Features

- **Multi-source Metric Collection**
  - CPU utilization and performance metrics
  - Memory usage and availability
  - Disk I/O and storage metrics
  - Network interface statistics
  - Thermal/temperature monitoring
  - Process-level metrics


- **InsightFinder Integration**
  - Direct metric streaming to InsightFinder platform
  - Automatic data formatting and submission
  - Built-in retry and error handling

## Prerequisites

- **Operating System**: Windows 10/11 or Windows Server 2016+
- **Go Version**: 1.23 or higher
- **Permissions**: Administrator privileges (required for PDH access and low-level metrics)

## Installation


### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/your-org/win-dex-agent.git
   cd win-dex-agent
   ```

2. Build the agent:
   ```bash
   go build -o win-dex-agent.exe
   ```

3. Run as administrator:
   ```bash
   # Right-click and select "Run as administrator" or use:
   runas /user:Administrator win-dex-agent.exe
   ```

## Configuration

Configure the agent using in the `main.go` file to build the executable:

```go
IFClient := insightfinder.CreateInsightFinderClient("https://app.insightfinder.com", "insightfinder_username", "insightfinder_licensekey", "insightfinder-project")
```

## Architecture

The agent consists of several key components:

- **`main.go`**: Entry point and orchestration
- **`collector/`**: Metric collection modules
    - `generalCollector.go`: Native Go-based system metrics
    - `pdhCollectorService.go`: Windows PDH counter collection
    - `generalCollectorModel.go` & `pdhDataModel.go`: Data models
    - `utils.go`: Shared utilities
- **`insightfinder/`**: InsightFinder API integration
    - `insightfinder.go`: Main API client
    - `model.go` & `projectDataModel.go`: API data structures
    - `utility.go`: Helper functions
- **`internal/`**: Internal libraries
    - `pdh/`: Windows PDH API bindings
    - `headers/`: Windows API headers
- **`cache/`**: Local data caching
- **`tool/`**: Utility tools

## Collected Metrics

### System Metrics
- CPU usage (per core and aggregate)
- Memory utilization and availability
- Disk read/write operations
- Disk space usage
- Network interface throughput
- System uptime

### Performance Counters (via PDH)
- Processor queue length
- Context switches
- Thread count
- Handle count
- Page faults
- And many more Windows-specific counters

## Running as a Service

To run the agent as a Windows service, you can use tools like NSSM (Non-Sucking Service Manager):

1. Download NSSM from https://nssm.cc/
2. Install the service:
   ```cmd
   nssm install WinDexAgent "C:\path\to\win-dex-agent.exe"
   nssm start WinDexAgent
   ```

## Troubleshooting

### Common Issues

**Permission Denied Errors**
- Ensure you're running the agent with Administrator privileges
- Some PDH counters require elevated permissions

**Connection Failures**
- Verify your InsightFinder credentials and server URL
- Check network connectivity and firewall settings
- Ensure outbound HTTPS traffic is allowed

**Missing Metrics**
- Some performance counters may not be available on all Windows versions
- Check Windows Event Viewer for PDH-related errors

### Logging

The agent logs to standard output. Redirect to a file for persistent logging:
```cmd
win-dex-agent.exe > agent.log 2>&1
```