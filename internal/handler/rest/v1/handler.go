package v1

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/the-witcher-knight/image-minimize-go/internal/service/imaging"
)

const (
	maxMemorySize = 5 * 1024 // 5MB
)

type Handler struct {
	imagingService imaging.Service
}

func New(imagingService imaging.Service) Handler {
	return Handler{
		imagingService: imagingService,
	}
}

func (hdl Handler) ReduceImageSize() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(maxMemorySize); err != nil {
			fmt.Fprintf(w, "parse multipart form error %v", err.Error())
			return
		}

		imgData, fileHeader, err := r.FormFile("file")
		if err != nil {
			fmt.Fprintf(w, "get file error %v", err.Error())
			return
		}

		resized, err := hdl.imagingService.Resize(r.Context(), imgData)
		if err != nil {
			fmt.Fprintf(w, "resize image error %v", err.Error())
			return
		}

		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(resized)))
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%v-resize.jpg", fileHeader.Filename))
		if _, err := w.Write(resized); err != nil {
			log.Println("unable to write image.")
		}
	}
}
