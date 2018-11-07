package util

import (
	"errors"
	"github.com/gocraft/web"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"io"
	"os"
)

func SaveImage(r *web.Request, field string, maxWidth uint, uuid string) error {
	file, handler, _ := r.FormFile(field)
	if file == nil {
		return errors.New("No file provided")
	}
	defer file.Close()
	var image image.Image
	mimeType := handler.Header["Content-Type"][0]
	switch mimeType {
	case "image/jpeg", "image/jpg":
		image, _ = jpeg.Decode(file)
	default:
		return errors.New("Wrong image extension")
	}
	out, _ := os.Create("./data/images/" + uuid + ".jpeg")
	defer out.Close()
	resized := resize.Resize(maxWidth, 0, image, resize.Bilinear)
	jpeg.Encode(out, resized, nil)
	return nil
}

func ResizeImage(infilename, outfilename string, size string) error {

	file, err := os.Open(infilename)
	defer file.Close()
	if err != nil {
		return err
	}

	image, err := jpeg.Decode(file)
	if err != nil {
		return err
	}

	out, err := os.Create(outfilename)
	defer out.Close()
	if err != nil {
		return err
	}

	switch size {
	case "small":
		resized := resize.Resize(32, 32, image, resize.Bilinear)
		jpeg.Encode(out, resized, nil)
	case "728x90":
		resized := resize.Resize(728, 0, image, resize.Bilinear)
		jpeg.Encode(out, resized, nil)
	case "200x200":
		resized := resize.Resize(200, 0, image, resize.Bilinear)
		jpeg.Encode(out, resized, nil)
	case "230x230":
		resized := resize.Resize(230, 230, image, resize.Bilinear)
		jpeg.Encode(out, resized, nil)
	case "300x300":
		resized := resize.Resize(300, 0, image, resize.Bilinear)
		jpeg.Encode(out, resized, nil)
	}

	return nil
}

func ServeImage(filename, size string, w web.ResponseWriter, r *web.Request) error {
	w.Header().Set("Content-type", "image/jpeg")
	w.Header().Set("Cache-control", "public, max-age=259200")

	originalFilename := "./data/images/" + filename + ".jpeg"
	resizedFilename := "./data/images/" + filename + "_" + size + ".jpeg"

	switch size {

	case "normal":
		file, err := os.Open(originalFilename)
		defer file.Close()
		if err != nil {
			return err
		}
		io.Copy(w, file)

	case "200x200", "728x90", "small", "300x300", "230x230":
		resizedFile, err := os.Open(resizedFilename)
		defer resizedFile.Close()

		if err != nil { // image does not exist. resize
			ResizeImage(originalFilename, resizedFilename, size)

			newlyResizedFile, err := os.Open(resizedFilename)
			defer newlyResizedFile.Close()
			if err != nil {
				return err
			}

			io.Copy(w, newlyResizedFile)
		} else {
			io.Copy(w, resizedFile)
		}

	}

	return nil

}
