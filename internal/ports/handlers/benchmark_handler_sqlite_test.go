package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	dbadapters "github.com/brunojet/go-infra-backend/internal/database/adapters"
	dbcontracts "github.com/brunojet/go-infra-backend/internal/database/contracts"
	"github.com/brunojet/go-infra-backend/internal/ports/repositories"
	svcimpl "github.com/brunojet/go-infra-backend/internal/ports/services"
	svc "github.com/brunojet/go-infra-backend/internal/ports/services/contracts"
	"github.com/gin-gonic/gin"
	gormlogger "gorm.io/gorm/logger"
)

type discardResponseWriter struct {
	header http.Header
}

func newDiscardResponseWriter() *discardResponseWriter {
	return &discardResponseWriter{header: make(http.Header)}
}

func (w *discardResponseWriter) Header() http.Header { return w.header }

func (w *discardResponseWriter) Write(p []byte) (int, error) { return len(p), nil }

func (w *discardResponseWriter) WriteHeader(statusCode int) {}

var benchmarkURL = &url.URL{Path: "/v1/complex"}

func newSimpleRequest(method, path string) *http.Request {
	return &http.Request{
		Method: method,
		URL:    &url.URL{Path: path},
		Header: make(http.Header, 1),
	}
}

func newJSONRequest(body []byte) *http.Request {
	req := &http.Request{
		Method:        http.MethodPost,
		URL:           benchmarkURL,
		Header:        make(http.Header, 1),
		Body:          io.NopCloser(bytes.NewReader(body)),
		ContentLength: int64(len(body)),
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func listWithSizeHandler(service svc.Service[ComplexDTO, ComplexModel], size int) gin.HandlerFunc {
	return func(c *gin.Context) {
		list, err := service.List(c.Request.Context(), size)
		if err != nil {
			SetResponseFromError(c, err)
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

func seedComplexRows(b *testing.B, service svc.Service[ComplexDTO, ComplexModel], n int) {
	ctx := context.Background()
	for i := 0; i < n; i++ {
		dto := genComplexDTO()
		dto.Name = dto.Name + "-" + strconv.Itoa(i)
		if err := service.Create(ctx, dto); err != nil {
			b.Fatalf("seed create failed: %v", err)
		}
	}
}

func seedOneComplexRow(b *testing.B, service svc.Service[ComplexDTO, ComplexModel]) string {
	dto := genComplexDTO()
	if err := service.Create(context.Background(), dto); err != nil {
		b.Fatalf("seed create failed: %v", err)
	}
	if dto.ID == "" {
		b.Fatal("seed create returned empty id")
	}
	return dto.ID
}

const latencySampleEvery = 10
const memSampleEvery = 200

type heapPeak struct {
	maxAlloc     uint64
	maxHeapInuse uint64
}

func (p *heapPeak) sample(ms runtime.MemStats) {
	if ms.Alloc > p.maxAlloc {
		p.maxAlloc = ms.Alloc
	}
	if ms.HeapInuse > p.maxHeapInuse {
		p.maxHeapInuse = ms.HeapInuse
	}
}

func atomicMaxUint64(addr *uint64, v uint64) {
	for {
		cur := atomic.LoadUint64(addr)
		if v <= cur {
			return
		}
		if atomic.CompareAndSwapUint64(addr, cur, v) {
			return
		}
	}
}

func reportHeapPeak(b *testing.B, p heapPeak) {
	b.ReportMetric(float64(p.maxAlloc)/(1024*1024), "heap_peak_alloc_MB")
	b.ReportMetric(float64(p.maxHeapInuse)/(1024*1024), "heap_peak_inuse_MB")
}

type latencySummary struct {
	samplesNs []int64
}

func (s *latencySummary) add(ns int64) {
	s.samplesNs = append(s.samplesNs, ns)
}

func percentileSorted(sorted []int64, p float64) int64 {
	if len(sorted) == 0 {
		return 0
	}
	if p <= 0 {
		return sorted[0]
	}
	if p >= 1 {
		return sorted[len(sorted)-1]
	}
	// Nearest-rank method: k = ceil(p*N)
	k := int(math.Ceil(p*float64(len(sorted)))) - 1
	if k < 0 {
		k = 0
	}
	if k >= len(sorted) {
		k = len(sorted) - 1
	}
	return sorted[k]
}

func reportPerfSummary(b *testing.B, label string, elapsed time.Duration, tx int, samples []int64) {
	if elapsed > 0 {
		b.ReportMetric(float64(tx)/elapsed.Seconds(), "tps")
	}
	if len(samples) == 0 {
		b.Logf("%s: transacoes=%d, tps=%.2f, sem amostras de latencia", label, tx, float64(tx)/elapsed.Seconds())
		return
	}
	sort.Slice(samples, func(i, j int) bool { return samples[i] < samples[j] })
	p90 := percentileSorted(samples, 0.90)
	p95 := percentileSorted(samples, 0.95)
	p99 := percentileSorted(samples, 0.99)
	b.ReportMetric(float64(p90)/1e6, "p90_ms")
	b.ReportMetric(float64(p95)/1e6, "p95_ms")
	b.ReportMetric(float64(p99)/1e6, "p99_ms")
	b.Logf("%s: transacoes=%d, elapsed=%s, tps=%.2f, amostras=%d (1/%d ops), p90=%.3fms p95=%.3fms p99=%.3fms",
		label,
		tx,
		elapsed,
		float64(tx)/elapsed.Seconds(),
		len(samples),
		latencySampleEvery,
		float64(p90)/1e6,
		float64(p95)/1e6,
		float64(p99)/1e6,
	)
}

// ComplexModel é uma entidade rica para testes de performance.
type ComplexModel struct {
	ID         int64 `gorm:"primaryKey,autoIncrement"`
	Name       sql.NullString
	Bio        sql.NullString `gorm:"type:text"`
	Age        sql.NullInt64
	Active     sql.NullBool
	Score      sql.NullFloat64
	TagsJSON   []byte `gorm:"type:blob"`
	AddrStreet sql.NullString
	AddrCity   sql.NullString
	AddrZip    sql.NullString
	CreatedAt  sql.NullTime
}

func (c ComplexModel) TableName() string { return "complex_models" }

// DTOs usados pelo handler/service
type AddressDTO struct {
	Street string `json:"street"`
	City   string `json:"city"`
	Zip    string `json:"zip"`
}

type ComplexDTO struct {
	ID        string     `json:"id,omitempty"`
	Name      string     `json:"name"`
	Bio       string     `json:"bio"`
	Age       int        `json:"age"`
	Active    bool       `json:"active"`
	Score     float64    `json:"score"`
	Tags      []string   `json:"tags"`
	Address   AddressDTO `json:"address"`
	CreatedAt time.Time  `json:"created_at"`
}

// Mapper para ComplexDTO <-> ComplexModel
type complexMapper struct{}

func (m *complexMapper) GetModelKey(id string) (map[string]any, error) {
	if id == "" {
		return map[string]any{}, nil
	}
	parsed, err := svcimpl.StringToInt64(id)
	if err != nil {
		return map[string]any{}, err
	}
	return map[string]any{"id": parsed}, nil
}

func (m *complexMapper) ToModel(dto *ComplexDTO) (ComplexModel, error) {
	if dto == nil {
		return ComplexModel{}, nil
	}
	var id int64
	if dto.ID != "" {
		parsed, err := svcimpl.StringToInt64(dto.ID)
		if err != nil {
			return ComplexModel{}, err
		}
		id = parsed
	}

	tagsJSON, err := svcimpl.ToJSONBytes(dto.Tags)
	if err != nil {
		return ComplexModel{}, err
	}

	nm := svcimpl.ToNullString(dto.Name)
	bio := svcimpl.ToNullString(dto.Bio)
	age := svcimpl.ToNullInt(dto.Age)
	active := svcimpl.ToNullBool(dto.Active)
	score := svcimpl.ToNullFloat(dto.Score)
	street := svcimpl.ToNullString(dto.Address.Street)
	city := svcimpl.ToNullString(dto.Address.City)
	zip := svcimpl.ToNullString(dto.Address.Zip)
	created := svcimpl.ToNullTime(&dto.CreatedAt)

	return ComplexModel{
		ID:         id,
		Name:       nm,
		Bio:        bio,
		Age:        age,
		Active:     active,
		Score:      score,
		TagsJSON:   tagsJSON,
		AddrStreet: street,
		AddrCity:   city,
		AddrZip:    zip,
		CreatedAt:  created,
	}, nil
}

func (m *complexMapper) ToDTO(model *ComplexModel, dto *ComplexDTO) {
	if dto == nil {
		return
	}
	if model == nil {
		*dto = ComplexDTO{}
		return
	}
	var tags []string
	if len(model.TagsJSON) > 0 {
		_ = svcimpl.FromJSONBytes(model.TagsJSON, &tags)
	}
	dto.ID = svcimpl.Int64ToString(model.ID)
	dto.Name = svcimpl.FromNullString(model.Name)
	dto.Bio = svcimpl.FromNullString(model.Bio)
	dto.Age = svcimpl.FromNullInt(model.Age)
	dto.Active = svcimpl.FromNullBool(model.Active)
	dto.Score = svcimpl.FromNullFloat(model.Score)
	dto.Tags = tags
	if model.AddrStreet.Valid || model.AddrCity.Valid || model.AddrZip.Valid {
		dto.Address = AddressDTO{Street: svcimpl.FromNullString(model.AddrStreet), City: svcimpl.FromNullString(model.AddrCity), Zip: svcimpl.FromNullString(model.AddrZip)}
	}
	if t := svcimpl.FromNullTime(model.CreatedAt); t != nil {
		dto.CreatedAt = *t
	}
}

// Constrói um DTO de teste com conteúdo relativamente pesado.
func genComplexDTO() *ComplexDTO {
	tags := make([]string, 100)
	ts := time.Now().Format("150405")
	for i := 0; i < 100; i++ {
		tags[i] = "tag-very-long-value-" + ts + "-" + strconv.Itoa(i)
	}
	var sb strings.Builder
	for i := 0; i < 200; i++ {
		sb.WriteString(" lorem ipsum dolor sit amet ")
	}
	bio := sb.String()
	return &ComplexDTO{
		Name:   "Benchmark Name",
		Bio:    bio,
		Age:    42,
		Active: true,
		Score:  12345.6789,
		Tags:   tags,
		Address: AddressDTO{
			Street: "Rua Exemplo 123",
			City:   "Cidade",
			Zip:    "00000-000",
		},
		CreatedAt: time.Now(),
	}
}

func setupBenchmarkDB(b *testing.B) (dbcontracts.Database, *repositories.GenericGormRepository[ComplexModel], svc.Service[ComplexDTO, ComplexModel]) {
	adapter, err := dbadapters.NewSQLite("memory")
	if err != nil {
		b.Fatalf("failed opening sqlite adapter: %v", err)
	}
	db := adapter.GormDB()
	// Silence standard logger output and GORM logger for cleaner benchmark measurements
	log.SetOutput(io.Discard)
	db.Logger = gormlogger.Default.LogMode(gormlogger.Silent)
	gin.SetMode(gin.ReleaseMode)
	if err := db.AutoMigrate(&ComplexModel{}); err != nil {
		_ = adapter.Close()
		b.Fatalf("auto migrate failed: %v", err)
	}
	repo := repositories.NewGormRepository[ComplexModel](db)
	mapper := &complexMapper{}
	service := svcimpl.NewServiceImpl(repo, mapper)
	return adapter, repo, service
}

func snapshotMemStats() runtime.MemStats {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return ms
}

func reportRuntimeSummary(b *testing.B, elapsed time.Duration, ops int, before runtime.MemStats, after runtime.MemStats, afterGC runtime.MemStats) {
	if ops <= 0 || elapsed <= 0 {
		return
	}
	allocDelta := float64(after.TotalAlloc - before.TotalAlloc)
	gcDelta := float64(after.NumGC - before.NumGC)
	pauseDeltaNs := float64(after.PauseTotalNs - before.PauseTotalNs)

	b.ReportMetric(allocDelta/float64(ops), "alloc_B/op_rt")
	b.ReportMetric((allocDelta/elapsed.Seconds())/(1024*1024), "alloc_MB/s")
	b.ReportMetric(gcDelta/float64(ops), "gc/op")
	b.ReportMetric((pauseDeltaNs/float64(ops))/1e6, "gc_pause_ms/op")
	b.ReportMetric(float64(afterGC.Alloc)/(1024*1024), "heap_live_MB")
}

func Benchmark_Handler_Create_SQLite(b *testing.B) {
	b.ReportAllocs()
	adapter, _, service := setupBenchmarkDB(b)
	defer adapter.Close()
	handler := NewGenericHandler[ComplexModel](service)

	dto := genComplexDTO()
	body, _ := json.Marshal(dto)

	// warmup GC and capture baseline
	runtime.GC()
	before := snapshotMemStats()

	b.ResetTimer()
	var lat latencySummary
	var peak heapPeak
	for i := 0; i < b.N; i++ {
		var start time.Time
		if i%latencySampleEvery == 0 {
			start = time.Now()
		}
		if i%memSampleEvery == 0 {
			peak.sample(snapshotMemStats())
		}
		w := newDiscardResponseWriter()
		c, _ := gin.CreateTestContext(w)
		c.Request = newJSONRequest(body)
		c.Request = c.Request.WithContext(context.Background())
		handler.Create(c)
		if c.Writer.Status() >= 400 {
			b.Fatalf("unexpected status %d", c.Writer.Status())
		}
		if !start.IsZero() {
			lat.add(time.Since(start).Nanoseconds())
		}
	}

	b.StopTimer()
	elapsed := b.Elapsed()
	reportPerfSummary(b, "sequencial", elapsed, b.N, lat.samplesNs)
	reportHeapPeak(b, peak)
	after := snapshotMemStats()
	runtime.GC()
	afterGC := snapshotMemStats()
	reportRuntimeSummary(b, elapsed, b.N, before, after, afterGC)
}

func Benchmark_Handler_Create_SQLite_Parallel(b *testing.B) {
	b.ReportAllocs()
	adapter, _, service := setupBenchmarkDB(b)
	defer adapter.Close()
	handler := NewGenericHandler[ComplexModel](service)
	dto := genComplexDTO()
	body, _ := json.Marshal(dto)

	runtime.GC()
	before := snapshotMemStats()

	b.ResetTimer()
	var hadError uint32
	var latAll latencySummary
	var latMu sync.Mutex
	var peakAlloc uint64
	var peakInuse uint64
	b.RunParallel(func(pb *testing.PB) {
		var local []int64
		var localPeak heapPeak
		op := 0
		for pb.Next() {
			op++
			var start time.Time
			if op%latencySampleEvery == 0 {
				start = time.Now()
			}
			if op%memSampleEvery == 0 {
				localPeak.sample(snapshotMemStats())
			}
			w := newDiscardResponseWriter()
			c, _ := gin.CreateTestContext(w)
			c.Request = newJSONRequest(body)
			c.Request = c.Request.WithContext(context.Background())
			handler.Create(c)
			if c.Writer.Status() >= 400 {
				atomic.StoreUint32(&hadError, 1)
			}
			if !start.IsZero() {
				local = append(local, time.Since(start).Nanoseconds())
			}
		}
		if len(local) > 0 {
			latMu.Lock()
			latAll.samplesNs = append(latAll.samplesNs, local...)
			latMu.Unlock()
		}
		if localPeak.maxAlloc > 0 {
			atomicMaxUint64(&peakAlloc, localPeak.maxAlloc)
		}
		if localPeak.maxHeapInuse > 0 {
			atomicMaxUint64(&peakInuse, localPeak.maxHeapInuse)
		}
	})
	if atomic.LoadUint32(&hadError) == 1 {
		b.Fatal("unexpected HTTP error status during parallel benchmark")
	}

	b.StopTimer()
	elapsed := b.Elapsed()
	reportPerfSummary(b, "paralelo", elapsed, b.N, latAll.samplesNs)
	reportHeapPeak(b, heapPeak{maxAlloc: atomic.LoadUint64(&peakAlloc), maxHeapInuse: atomic.LoadUint64(&peakInuse)})
	after := snapshotMemStats()
	runtime.GC()
	afterGC := snapshotMemStats()
	reportRuntimeSummary(b, elapsed, b.N, before, after, afterGC)
}

func Benchmark_Handler_GetByID_SQLite(b *testing.B) {
	b.ReportAllocs()
	adapter, _, service := setupBenchmarkDB(b)
	defer adapter.Close()
	handler := NewGenericHandler[ComplexModel](service)

	// Seed a single row to read repeatedly.
	id := seedOneComplexRow(b, service)

	runtime.GC()
	before := snapshotMemStats()

	b.ResetTimer()
	var lat latencySummary
	var peak heapPeak
	for i := 0; i < b.N; i++ {
		var start time.Time
		if i%latencySampleEvery == 0 {
			start = time.Now()
		}
		if i%memSampleEvery == 0 {
			peak.sample(snapshotMemStats())
		}
		w := newDiscardResponseWriter()
		c, _ := gin.CreateTestContext(w)
		c.Request = newSimpleRequest(http.MethodGet, "/v1/complex/"+id).WithContext(context.Background())
		c.Params = gin.Params{{Key: "id", Value: id}}
		handler.GetByID(c)
		if c.Writer.Status() >= 400 {
			b.Fatalf("unexpected status %d", c.Writer.Status())
		}
		if !start.IsZero() {
			lat.add(time.Since(start).Nanoseconds())
		}
	}

	b.StopTimer()
	elapsed := b.Elapsed()
	reportPerfSummary(b, "get_by_id", elapsed, b.N, lat.samplesNs)
	reportHeapPeak(b, peak)
	after := snapshotMemStats()
	runtime.GC()
	afterGC := snapshotMemStats()
	reportRuntimeSummary(b, elapsed, b.N, before, after, afterGC)
}

func Benchmark_Handler_List100_SQLite(b *testing.B) {
	b.ReportAllocs()
	adapter, _, service := setupBenchmarkDB(b)
	defer adapter.Close()
	listH := listWithSizeHandler(service, 100)

	// Seed a fixed dataset of 100 rows so List returns 100 items consistently.
	seedComplexRows(b, service, 100)

	runtime.GC()
	before := snapshotMemStats()

	b.ResetTimer()
	var lat latencySummary
	var peak heapPeak
	for i := 0; i < b.N; i++ {
		var start time.Time
		if i%latencySampleEvery == 0 {
			start = time.Now()
		}
		if i%memSampleEvery == 0 {
			peak.sample(snapshotMemStats())
		}
		w := newDiscardResponseWriter()
		c, _ := gin.CreateTestContext(w)
		c.Request = newSimpleRequest(http.MethodGet, "/v1/complex").WithContext(context.Background())
		listH(c)
		if c.Writer.Status() >= 400 {
			b.Fatalf("unexpected status %d", c.Writer.Status())
		}
		if !start.IsZero() {
			lat.add(time.Since(start).Nanoseconds())
		}
	}

	b.StopTimer()
	elapsed := b.Elapsed()
	reportPerfSummary(b, "list_100", elapsed, b.N, lat.samplesNs)
	reportHeapPeak(b, peak)
	after := snapshotMemStats()
	runtime.GC()
	afterGC := snapshotMemStats()
	reportRuntimeSummary(b, elapsed, b.N, before, after, afterGC)
}
