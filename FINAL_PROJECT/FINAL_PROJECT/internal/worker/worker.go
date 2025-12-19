package worker

import (
	"clinic-cli/internal/model"
	"context"
	"log"
	"time"
)

func StartEmailWorker(ctx context.Context, appChan <-chan model.Appointment) {
	log.Println("Background Email Worker Started...")
	for {
		select {
		case <-ctx.Done():
			log.Println("Email Worker shutting down...")
			return
		case app := <-appChan:
			time.Sleep(500 * time.Millisecond)
			log.Printf("[WORKER] Sending confirmation email for Appointment ID %d (Patient: %d, Doctor: %d)\n",
				app.ID, app.PatientID, app.DoctorID)
		}
	}
}
