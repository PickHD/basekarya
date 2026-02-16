package attendance

import (
	"hris-backend/pkg/logger"
	"sync"
	"time"

	"gorm.io/gorm"
)

type GeocodeWorker interface {
	Start(workerCount int)
	Enqueue(job GeocodeJob)
	Stop()
}

type GeocodeJob struct {
	AttendanceID uint
	Latitude     float64
	Longitude    float64
	IsCheckout   bool
}

type geocodeWorker struct {
	db      *gorm.DB
	fetcher LocationFetcher
	queue   chan GeocodeJob
	wg      *sync.WaitGroup
	quit    chan bool
}

func NewGeocodeWorker(db *gorm.DB, fetcher LocationFetcher, bufferSize int) GeocodeWorker {
	return &geocodeWorker{
		db:      db,
		fetcher: fetcher,
		queue:   make(chan GeocodeJob, bufferSize),
		wg:      &sync.WaitGroup{},
		quit:    make(chan bool),
	}
}

func (w *geocodeWorker) Start(workerCount int) {
	for i := 0; i < workerCount; i++ {
		w.wg.Add(1)
		go w.runWorker(i)
	}
}

func (w *geocodeWorker) Enqueue(job GeocodeJob) {
	select {
	case w.queue <- job:
	default:
		logger.Warn("Geocode queue is full, skipping job for ID:", job.AttendanceID)
	}
}

func (w *geocodeWorker) Stop() {
	logger.Info("Stopping Geocode Workers...")
	close(w.queue)
	w.wg.Wait()
	logger.Info("All Geocode Workers stopped.")
}

func (w *geocodeWorker) runWorker(id int) {
	defer w.wg.Done()
	logger.Infof("Geocode Worker #%d Started", id)

	limiter := time.NewTicker(1500 * time.Millisecond)
	defer limiter.Stop()

	for job := range w.queue {
		<-limiter.C

		w.processJob(job)
	}
}

func (w *geocodeWorker) processJob(job GeocodeJob) {
	address := w.fetcher.GetAddressFromCoords(job.Latitude, job.Longitude)

	if address == "" {
		logger.Warnf("Empty address for ID %d, skipping update", job.AttendanceID)
		return
	}

	columnName := "check_in_address"
	if job.IsCheckout {
		columnName = "check_out_address"
	}

	result := w.db.Table("attendances").Where("id = ?", job.AttendanceID).Update(columnName, address)
	if result.Error != nil {
		logger.Errorw("Failed to update address", "error", result.Error, "id", job.AttendanceID)
	} else {
		logger.Infof("Success update address for ID %d", job.AttendanceID)
	}

}
