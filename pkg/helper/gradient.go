package helper

import (
	"fmt"
	"math/rand"
	"time"
)

// GradientPalette contains 12 curated gradient pairs from reference design
var GradientPalette = []struct {
	Start string
	End   string
}{
	// Row 1
	{"#FFA500", "#FF4500"}, // Orange to Red
	{"#87CEEB", "#1E90FF"}, // Light Blue to Blue
	{"#FFD700", "#32CD32"}, // Yellow to Green
	{"#FFDAB9", "#FF6347"}, // Peach to Coral

	// Row 2
	{"#FFB6C1", "#8A2BE2"}, // Pink to Purple
	{"#FFD700", "#FF4500"}, // Yellow to Orange-Red
	{"#B0C4DE", "#708090"}, // Light Steel to Slate Gray
	{"#00CED1", "#20B2AA"}, // Cyan to Teal

	// Row 3
	{"#DDA0DD", "#4169E1"}, // Plum to Royal Blue
	{"#00FFFF", "#3CB371"}, // Cyan to Medium Sea Green
	{"#FFFACD", "#9370DB"}, // Lemon Chiffon to Medium Purple
	{"#4169E1", "#FF8C00"}, // Royal Blue to Dark Orange
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateRandomGradient returns a random gradient color pair from the palette
func GenerateRandomGradient() (startColor, endColor string) {
	idx := rand.Intn(len(GradientPalette))
	return GradientPalette[idx].Start, GradientPalette[idx].End
}

// GenerateGradientFromSeed generates a consistent gradient based on a seed (e.g., user UUID)
// This ensures the same user always gets the same gradient
func GenerateGradientFromSeed(seed string) (startColor, endColor string) {
	// Create a deterministic hash from the seed
	hash := 0
	for _, c := range seed {
		hash = (hash*31 + int(c)) % len(GradientPalette)
	}
	if hash < 0 {
		hash = -hash
	}
	idx := hash % len(GradientPalette)
	return GradientPalette[idx].Start, GradientPalette[idx].End
}

// GenerateRandomHexColor generates a random hex color
func GenerateRandomHexColor() string {
	return fmt.Sprintf("#%02X%02X%02X", rand.Intn(256), rand.Intn(256), rand.Intn(256))
}
