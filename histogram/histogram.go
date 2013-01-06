package histogram

import (
	// "fmt"
	"image"
	"log"
	"net/http"
	"html/template"	
	_ "image/gif"
	_ "image/png"
	_ "image/jpeg"
)

var histogramTemplate, templateErr = template.ParseGlob("templates/**")

func init() {
	if templateErr != nil {
		log.Fatal(templateErr)
		return
	}
	http.HandleFunc("/histogram", makeHistogram)
}

func makeHistogram(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == "POST" {
		processImage(w, r)
	} else {
		http.Error(w, "Only POST requests are allowed for /histogram", http.StatusInternalServerError)
	}
}

func processImage(w http.ResponseWriter, r *http.Request) {
	// Open the file.
	file, _, err := r.FormFile("image_file")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Decode the image.
	m, _, err := image.Decode(file)
	if err != nil {
		log.Fatal("Error decoding image: ", err)
	}

	bounds := m.Bounds()

	widthPixels := int(bounds.Max.X - bounds.Min.X)
	heightPixels := int(bounds.Max.Y - bounds.Min.Y)
	// init histogram
	histogram := make([][]float64, 3)
	for i := range histogram {
		histogram[i] = make([]float64, widthPixels)
	}
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
			// values should be between 0 and 1
			relR := float64(r)/65535.0
			relG := float64(g)/65535.0
			relB := float64(b)/65535.0
			histogram[0][x] += relR
			histogram[1][x] += relG
			histogram[2][x] += relB
		}
	}
	
	// create median for each x value
	for i, values := range histogram {
		for j := range values {
			histogram[i][j] = roundedPercentage(histogram[i][j], float64(heightPixels))
		}
	}
	
	if err := histogramTemplate.Execute(w, histogram); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func roundedPercentage(value, count float64) float64 {
	return float64(int(((value / count) * 100) + 0.5))
}

