// prometheusTestServer 是一个模拟真实生产环境的 Prometheus 指标服务器。
// 启动后会持续向 :2112/metrics 暴露模拟数据，包含正常、偶发错误、告警三种状态。
// 同时模拟 5 个 HTTP 接口供 Prometheus 抓取 QPS / 延迟 / 错误率等指标。
package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ──────────────────────────────────────────────
// 指标定义
// ──────────────────────────────────────────────

var (
	// HTTP 请求总数（按接口、方法、状态码分组）
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "HTTP 请求总数",
		},
		[]string{"handler", "method", "status_code"},
	)

	// HTTP 请求延迟分布
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP 请求耗时分布（秒）",
			Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
		},
		[]string{"handler", "method"},
	)

	// 当前活跃连接数
	activeConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "active_connections",
		Help: "当前活跃连接数",
	})

	// 错误率（每分钟）
	errorRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "error_rate_per_minute",
			Help: "每分钟错误率（百分比）",
		},
		[]string{"handler"},
	)

	// CPU 使用率（模拟）
	cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "process_cpu_usage_percent",
		Help: "进程 CPU 使用率（模拟）",
	})

	// 内存使用量（模拟）
	memoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "process_memory_bytes",
		Help: "进程内存使用量（字节，模拟）",
	})

	// 数据库连接池使用情况
	dbPoolUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_pool_connections_used",
		Help: "数据库连接池已用连接数",
	})
	dbPoolMax = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_pool_connections_max",
		Help: "数据库连接池最大连接数",
	})

	// 队列积压深度
	queueDepth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "message_queue_depth",
			Help: "消息队列积压深度",
		},
		[]string{"queue"},
	)
)

// 模拟的 API 接口列表
var handlers = []string{
	"/api/v1/order",
	"/api/v1/user",
	"/api/v1/payment",
	"/api/v1/inventory",
	"/api/v1/notification",
}

func init() {
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		activeConnections,
		errorRate,
		cpuUsage,
		memoryUsage,
		dbPoolUsed,
		dbPoolMax,
		queueDepth,
	)
	dbPoolMax.Set(100)
}

// ──────────────────────────────────────────────
// 模拟数据生成（场景：正常 → 偶发错误 → 告警）
// ──────────────────────────────────────────────

type scenario int

const (
	scNormal   scenario = iota // 正常
	scDegraded                 // 轻微降级
	scAlert                    // 告警（高错误率 / 高延迟）
)

func currentScenario(t time.Time) scenario {
	// 强制全天候处于告警状态，方便测试
	return scAlert
}

func scenarioName(s scenario) string {
	switch s {
	case scNormal:
		return "正常"
	case scDegraded:
		return "降级"
	case scAlert:
		return "⚠️  告警"
	}
	return "未知"
}

// simulateMetrics 每秒模拟一批请求并更新所有指标
func simulateMetrics() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		sc := currentScenario(t)
		fmt.Printf("[%s] 当前场景: %s\n", t.Format("15:04:05"), scenarioName(sc))

		// 基础 QPS 参数
		baseQPS := 50.0
		baseErrorPct := 0.01  // 1%
		baseLatencyMs := 80.0 // 毫秒

		switch sc {
		case scDegraded:
			baseErrorPct = 0.08 // 8%
			baseLatencyMs = 300
		case scAlert:
			baseQPS = 120        // 流量激增
			baseErrorPct = 0.35  // 35% 错误率，触发告警
			baseLatencyMs = 1500 // 高延迟
		}

		// 更新 CPU / 内存
		cpuBase := 25.0
		switch sc {
		case scDegraded:
			cpuBase = 60
		case scAlert:
			cpuBase = 85
		}
		cpuUsage.Set(cpuBase + rand.Float64()*10 - 5)
		memUsage := 200*1024*1024 + rand.Float64()*50*1024*1024
		if sc == scAlert {
			memUsage += 300 * 1024 * 1024
		}
		memoryUsage.Set(memUsage)

		// 更新数据库连接池
		dbUsed := 20.0 + rand.Float64()*20
		if sc == scAlert {
			dbUsed = 90 + rand.Float64()*10
		}
		dbPoolUsed.Set(dbUsed)

		// 更新活跃连接
		conn := baseQPS * (0.8 + rand.Float64()*0.4)
		activeConnections.Set(conn)

		// 更新队列积压
		for _, q := range []string{"order_queue", "notify_queue", "payment_queue"} {
			depth := rand.Float64() * 10
			if sc == scAlert {
				depth = 500 + rand.Float64()*200
			} else if sc == scDegraded {
				depth = 50 + rand.Float64()*50
			}
			queueDepth.WithLabelValues(q).Set(depth)
		}

		// 为每个接口模拟请求
		for _, h := range handlers {
			qps := baseQPS * (0.5 + rand.Float64())
			reqCount := int(qps)

			errRate := baseErrorPct * (0.5 + rand.Float64())

			// 按 sin 波增加延迟波动
			wave := math.Sin(float64(t.Unix())/10) * 0.2
			latency := baseLatencyMs * (1 + wave + rand.Float64()*0.3) / 1000.0

			successCount := int(float64(reqCount) * (1 - errRate))
			errorCount := reqCount - successCount

			httpRequestsTotal.WithLabelValues(h, "GET", "200").Add(float64(successCount))
			if errorCount > 0 {
				httpRequestsTotal.WithLabelValues(h, "GET", "500").Add(float64(errorCount))
			}
			httpRequestDuration.WithLabelValues(h, "GET").Observe(latency)
			errorRate.WithLabelValues(h).Set(errRate * 100)
		}
	}
}

func main() {
	go simulateMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Prometheus Test Server 运行中")
		fmt.Fprintln(w, "访问 /metrics 查看指标")
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "场景循环（每 3 分钟）：")
		fmt.Fprintln(w, "  0~90s   : 正常")
		fmt.Fprintln(w, "  90~150s : 轻微降级（错误率 8%，延迟 300ms）")
		fmt.Fprintln(w, "  150~180s: 告警（错误率 35%，延迟 1.5s，CPU 85%，队列积压）")
	})

	addr := ":2112"
	log.Printf("Prometheus 测试服务器启动，监听 %s", addr)
	log.Printf("Prometheus scrape_configs 配置示例：")
	log.Printf("  - job_name: 'oncall-test'")
	log.Printf("    static_configs:")
	log.Printf("      - targets: ['localhost:2112']")
	log.Fatal(http.ListenAndServe(addr, nil))
}
