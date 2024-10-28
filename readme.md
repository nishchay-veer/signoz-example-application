# Go Runtime Metrics with OpenTelemetry and SigNoz

This project demonstrates how to collect and visualize Go runtime metrics using OpenTelemetry, pushing the data to SigNoz for monitoring and analysis. The application collects various runtime metrics including memory usage, goroutine stats, garbage collection metrics, and more.

## 📊 Features

- Comprehensive Go runtime metrics collection
- OpenTelemetry integration
- SigNoz visualization
- Docker and Docker Compose setup
- Automatic metric collection and export
- Pre-configured dashboard templates

## 🔧 Prerequisites

Before you begin, ensure you have the following installed:
- Go 1.21 or later
- Docker and Docker Compose
- A SigNoz Cloud account
- Git

## 🏗️ Metrics Collected

The application collects the following metrics categories:

### Memory Metrics
- Heap Allocation
- Heap Idle
- Heap In Use
- Heap Objects
- Heap Released
- Heap System

### Goroutines Metrics
- Active Goroutines
- Total Goroutines Created

### Garbage Collection Metrics
- GC Pause Duration
- GC Count
- GC Forced
- CPU Fraction
- GC System

### CPU Metrics
- Goroutine Execution Time
- GC Time
- System Time

### Additional Metrics
- OS Threads
- Stack Metrics
- Mutex and Semaphores

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/go-runtime-metrics.git
cd go-runtime-metrics
```

### 2. Configure Environment Variables

Create a `.env` file in the project root:

```env
SIGNOZ_TOKEN=your_signoz_token_here
```

### 3. Update SigNoz Configuration

In `otel-config.yaml`, update the following fields:
- Replace `{region}` with your SigNoz cloud region
- Verify the endpoint URL matches your SigNoz instance

### 4. Build and Run

```bash
# Download Go dependencies
go mod download

# Start the application using Docker Compose
docker-compose up --build
```

The application will be available at `http://localhost:8080`

## 📁 Project Structure

```
go-runtime-metrics/
├── main.go              # Main application file
├── Dockerfile           # Docker configuration
├── docker-compose.yml   # Docker Compose configuration
├── otel-config.yaml     # OpenTelemetry Collector configuration
├── .env                 # Environment variables (create this)
└── README.md           # This file
```

## ⚙️ Configuration

### Application Configuration
- The application runs on port 8080
- Metrics are collected every 10 seconds
- OpenTelemetry collector runs on ports 4317 (gRPC) and 4318 (HTTP)

### SigNoz Dashboard Setup

1. Log in to your SigNoz Cloud account
2. Create a new dashboard
3. Import the provided dashboard template
4. Verify metrics are being received

## 📊 Available Dashboards

The repository includes dashboard templates for:
- Memory Usage Overview
- Goroutine Statistics
- Garbage Collection Analysis
- CPU Usage Metrics
- Thread and Stack Analysis

## 🔍 Monitoring Tips

1. Monitor heap usage patterns to optimize memory allocation
2. Watch goroutine counts to detect potential leaks
3. Track GC metrics to optimize garbage collection
4. Monitor CPU usage patterns
5. Keep an eye on mutex contention

## 🛠️ Troubleshooting

### Common Issues

1. **Metrics not appearing in SigNoz**
   - Verify your SigNoz token is correct
   - Check the OpenTelemetry collector logs
   - Ensure the endpoint URL is correct

2. **Application won't start**
   - Verify all environment variables are set
   - Check Docker logs for errors
   - Ensure required ports are available

3. **High resource usage**
   - Adjust metric collection interval
   - Review batch processing settings
   - Check for goroutine leaks

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- OpenTelemetry Community
- SigNoz Team
- Go Team for runtime metrics support

## 📞 Support

For support:
1. Open an issue in the repository
2. Check the SigNoz documentation
3. Review OpenTelemetry documentation

## 🔄 Updates

Check the releases page for updates and changelog.