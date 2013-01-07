package histogram

import (
	// "fmt"
	"html/template"
	"image"
	"os"
	"io"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
)

var templates, templateErr = template.ParseGlob("templates/**")

func init() {
	if templateErr != nil {
		log.Fatal(templateErr)
		return
	}
	setupHandlers();
}

func setupHandlers() {
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/histogram", histogram)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {	
	http.ServeFile(w, r, "public/html/index.html")
}

func histogram(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != "POST" {
		http.Error(w, "Only POST requests allowed", http.StatusInternalServerError)
		return
	}
	
	var file io.ReadCloser	
	var err error
	
	// we will either receive example_image_file or image_data
	if r.FormValue("example_image_file") != "" {
		file, err = os.Open("public/images/" + r.FormValue("example_image_file"))
	} else {
		file, _, err = r.FormFile("image_data")
	}
		
	if err != nil { 
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	defer file.Close()
	
	// Decode the image.
	m, _, err := image.Decode(file)
	if err != nil {
		errMsg := "Error decoding image: \n" + err.Error()
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	
	histogram := generateHistogramForImage(m)
	
	if err := templates.Execute(w, histogram); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func generateHistogramForImage(m image.Image) [][]float64 {

	bounds := m.Bounds()
	widthPixels := int(bounds.Max.X - bounds.Min.X)
	heightPixels := int(bounds.Max.Y - bounds.Min.Y)
	histogram := make([][]float64, 3)
	
	for i := range histogram {
		histogram[i] = make([]float64, widthPixels)
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
			// values should be between 0 and 1
			relR := float64(r) / 65535.0
			relG := float64(g) / 65535.0
			relB := float64(b) / 65535.0
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

	return histogram
}

func roundedPercentage(value, count float64) float64 {
	return float64(int(((value / count) * 100) + 0.5))
}
