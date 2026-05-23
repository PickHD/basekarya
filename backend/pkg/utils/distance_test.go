package utils

import (
	"math"
	"testing"
)

func TestCalculateDistance_SamePoint(t *testing.T) {
	d := CalculateDistance(-6.200000, 106.816666, -6.200000, 106.816666)
	if d > 1 {
		t.Errorf("expected ~0 for same point, got %f", d)
	}
}

func TestCalculateDistance_JakartaToBandung(t *testing.T) {
	jakartaLat, jakartaLon := -6.2088, 106.8456
	bandungLat, bandungLon := -6.9175, 107.6191

	d := CalculateDistance(jakartaLat, jakartaLon, bandungLat, bandungLon)
	dKm := d / 1000.0

	if math.Abs(dKm-150) > 50 {
		t.Errorf("expected ~150km Jakarta-Bandung, got %.1fkm", dKm)
	}
}
