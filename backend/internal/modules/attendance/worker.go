package attendance

import (
	"hris-backend/pkg/logger"
	"time"

	"gorm.io/gorm"
)

type GeocodeJob struct {
	AttendanceID uint
	Latitude     float64
	Longitude    float64
	IsCheckout   bool
}

func StartGeocodeWorker(db *gorm.DB, fetcher LocationFetcher, queue <-chan GeocodeJob) {
	go func() {
		logger.Info("Geocode Worker Started....")

		rateLimiter := time.NewTicker(1500 * time.Millisecond)
		defer rateLimiter.Stop()

		for job := range queue {
			<-rateLimiter.C

			processJob(db, fetcher, job)
		}
	}()
}

func processJob(db *gorm.DB, fetcher LocationFetcher, job GeocodeJob) {
	address := fetcher.GetAddressFromCoords(job.Latitude, job.Longitude)

	logger.Infof("Address found for ID %d: %s", job.AttendanceID, address)

	columnName := "check_in_address"
	if job.IsCheckout {
		columnName = "check_out_address"
	}

	result := db.Model(&Attendance{}).Where("id = ?", job.AttendanceID).Update(columnName, address)

	if result.Error != nil {
		logger.Errorw("Failed to update address", result.Error)
	}

	logger.Infof("Success update address from job ID %d", job.AttendanceID)
}
