package main

import (
	"context"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// MetricsCollector holds all our metrics instruments
type MetricsCollector struct {
	// Memory Metrics
	heapAlloc       metric.Int64ObservableGauge
	heapIdle        metric.Int64ObservableGauge
	heapInUse       metric.Int64ObservableGauge
	heapObjects     metric.Int64ObservableGauge
	heapReleased    metric.Int64ObservableGauge
	heapSys         metric.Int64ObservableGauge
	gcPauseDuration metric.Float64Histogram
	gcCount         metric.Int64Counter
	gcForced        metric.Int64Counter

	// Goroutine Metrics
	activeGoroutines metric.Int64ObservableGauge
	totalGoroutines  metric.Int64Counter

	// GC Metrics
	gcCPUFraction metric.Float64ObservableGauge
	gcSys         metric.Int64ObservableGauge

	// CPU Metrics
	goroutineExecTime metric.Int64Counter
	gcTime            metric.Int64Counter
	systemTime        metric.Int64Counter

	// Thread Metrics
	osThreads metric.Int64ObservableGauge

	// Stack Metrics
	stackInUse metric.Int64ObservableGauge
	stackSys   metric.Int64ObservableGauge

	// Mutex Metrics
	mutexWaitTime  metric.Int64Counter
	mutexLockCount metric.Int64Counter

	// Runtime stats
	lastPause uint64
	lastGC    uint32
}

func newMetricsCollector(meter metric.Meter) (*MetricsCollector, error) {
	mc := &MetricsCollector{}
	var err error

	// Memory Metrics
	mc.heapAlloc, err = meter.Int64ObservableGauge("go_memstats_heap_alloc_bytes",
		metric.WithDescription("Current bytes allocated to heap objects"))
	if err != nil {
		return nil, err
	}

	mc.heapIdle, err = meter.Int64ObservableGauge("go_memstats_heap_idle_bytes",
		metric.WithDescription("Bytes in idle (unused) spans"))
	if err != nil {
		return nil, err
	}

	mc.heapInUse, err = meter.Int64ObservableGauge("go_memstats_heap_inuse_bytes",
		metric.WithDescription("Bytes in in-use spans"))
	if err != nil {
		return nil, err
	}

	mc.heapObjects, err = meter.Int64ObservableGauge("go_memstats_heap_objects",
		metric.WithDescription("Number of allocated heap objects"))
	if err != nil {
		return nil, err
	}

	mc.heapReleased, err = meter.Int64ObservableGauge("go_memstats_heap_released_bytes",
		metric.WithDescription("Bytes released to the OS"))
	if err != nil {
		return nil, err
	}

	mc.heapSys, err = meter.Int64ObservableGauge("go_memstats_heap_sys_bytes",
		metric.WithDescription("Bytes obtained from system for heap"))
	if err != nil {
		return nil, err
	}

	mc.gcPauseDuration, err = meter.Float64Histogram("go_gc_duration_seconds",
		metric.WithDescription("GC pause duration"))
	if err != nil {
		return nil, err
	}

	mc.gcCount, err = meter.Int64Counter("go_gc_cycles_total",
		metric.WithDescription("Number of completed GC cycles"))
	if err != nil {
		return nil, err
	}

	mc.gcForced, err = meter.Int64Counter("go_gc_forced_total",
		metric.WithDescription("Number of forced GC cycles"))
	if err != nil {
		return nil, err
	}

	// Goroutine Metrics
	mc.activeGoroutines, err = meter.Int64ObservableGauge("go_goroutines",
		metric.WithDescription("Number of goroutines that currently exist"))
	if err != nil {
		return nil, err
	}

	mc.totalGoroutines, err = meter.Int64Counter("go_goroutines_created_total",
		metric.WithDescription("Total number of goroutines created"))
	if err != nil {
		return nil, err
	}

	// GC Metrics
	mc.gcCPUFraction, err = meter.Float64ObservableGauge("go_gc_cpu_fraction",
		metric.WithDescription("Fraction of CPU time used by GC"))
	if err != nil {
		return nil, err
	}

	mc.gcSys, err = meter.Int64ObservableGauge("go_memstats_gc_sys_bytes",
		metric.WithDescription("Bytes used for garbage collection system metadata"))
	if err != nil {
		return nil, err
	}

	// CPU Metrics
	mc.goroutineExecTime, err = meter.Int64Counter("go_runtime_cpu_goroutine_seconds_total",
		metric.WithDescription("Total user and system CPU time spent in goroutines"))
	if err != nil {
		return nil, err
	}

	mc.gcTime, err = meter.Int64Counter("go_gc_cpu_seconds_total",
		metric.WithDescription("Total CPU time spent in GC"))
	if err != nil {
		return nil, err
	}

	mc.systemTime, err = meter.Int64Counter("go_runtime_sys_seconds_total",
		metric.WithDescription("Total system CPU time"))
	if err != nil {
		return nil, err
	}

	// Thread Metrics
	mc.osThreads, err = meter.Int64ObservableGauge("go_threads",
		metric.WithDescription("Number of OS threads created"))
	if err != nil {
		return nil, err
	}

	// Stack Metrics
	mc.stackInUse, err = meter.Int64ObservableGauge("go_memstats_stack_inuse_bytes",
		metric.WithDescription("Bytes used by stack allocator"))
	if err != nil {
		return nil, err
	}

	mc.stackSys, err = meter.Int64ObservableGauge("go_memstats_stack_sys_bytes",
		metric.WithDescription("Bytes obtained from system for stack allocator"))
	if err != nil {
		return nil, err
	}

	// Mutex Metrics
	mc.mutexWaitTime, err = meter.Int64Counter("go_mutex_wait_seconds_total",
		metric.WithDescription("Total time spent waiting for mutex locks"))
	if err != nil {
		return nil, err
	}

	mc.mutexLockCount, err = meter.Int64Counter("go_mutex_lock_total",
		metric.WithDescription("Total number of mutex lock operations"))
	if err != nil {
		return nil, err
	}

	// Register callbacks for observable metrics
	_, err = meter.RegisterCallback(
		func(_ context.Context, o metric.Observer) error {
			var stats runtime.MemStats
			runtime.ReadMemStats(&stats)

			o.ObserveInt64(mc.heapAlloc, int64(stats.HeapAlloc))
			o.ObserveInt64(mc.heapIdle, int64(stats.HeapIdle))
			o.ObserveInt64(mc.heapInUse, int64(stats.HeapInuse))
			o.ObserveInt64(mc.heapObjects, int64(stats.HeapObjects))
			o.ObserveInt64(mc.heapReleased, int64(stats.HeapReleased))
			o.ObserveInt64(mc.heapSys, int64(stats.HeapSys))
			o.ObserveInt64(mc.activeGoroutines, int64(runtime.NumGoroutine()))
			o.ObserveFloat64(mc.gcCPUFraction, stats.GCCPUFraction)
			o.ObserveInt64(mc.gcSys, int64(stats.GCSys))
			o.ObserveInt64(mc.osThreads, int64(runtime.NumCPU()))
			o.ObserveInt64(mc.stackInUse, int64(stats.StackInuse))
			o.ObserveInt64(mc.stackSys, int64(stats.StackSys))

			// Update GC metrics
			if stats.NumGC > mc.lastGC {
				delta := stats.NumGC - mc.lastGC
				mc.gcCount.Add(context.Background(), int64(delta))
				if stats.LastGC > mc.lastPause {
					pauseNs := stats.LastGC - mc.lastPause
					mc.gcPauseDuration.Record(context.Background(), float64(pauseNs)/1e9)
				}
				mc.lastGC = stats.NumGC
				mc.lastPause = stats.LastGC
			}

			return nil
		},
		mc.heapAlloc,
		mc.heapIdle,
		mc.heapInUse,
		mc.heapObjects,
		mc.heapReleased,
		mc.heapSys,
		mc.activeGoroutines,
		mc.gcCPUFraction,
		mc.gcSys,
		mc.osThreads,
		mc.stackInUse,
		mc.stackSys,
	)

	return mc, err
}

func main() {
	// Create a Prometheus exporter
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}

	// Create a meter provider with the Prometheus exporter
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	defer provider.Shutdown(context.Background())

	// Set the global meter provider
	otel.SetMeterProvider(provider)

	// Get a meter
	meter := provider.Meter("go-runtime-metrics")

	// Create metrics collector
	collector, err := newMetricsCollector(meter)
	if err != nil {
		log.Fatal(err)
	}

	// Set up Gin server
	r := gin.Default()

	// Add test routes that affect metrics
	r.GET("/allocate-memory", func(c *gin.Context) {
		// Allocate some memory to see the effect
		_ = make([]byte, 1024*1024) // 1MB
		c.JSON(http.StatusOK, gin.H{"message": "Allocated memory"})
	})

	r.GET("/spawn-goroutines", func(c *gin.Context) {
		// Spawn 10 goroutines that sleep briefly
		for i := 0; i < 10; i++ {
			go func() {
				collector.totalGoroutines.Add(c.Request.Context(), 1)
				time.Sleep(time.Second)
			}()
		}
		c.JSON(http.StatusOK, gin.H{"message": "Spawned goroutines"})
	})

	r.GET("/simulate-mutex", func(c *gin.Context) {
		start := time.Now()
		time.Sleep(time.Millisecond * 100) // Simulate mutex wait
		collector.mutexWaitTime.Add(c.Request.Context(), time.Since(start).Nanoseconds())
		collector.mutexLockCount.Add(c.Request.Context(), 1)
		c.JSON(http.StatusOK, gin.H{"message": "Simulated mutex operation"})
	})

	// Expose Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Run the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
