package main

import (
	"bytes"
	"image"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

type dimensions struct {
	Width  uint
	Height uint
}

type ratio struct {
	Width  float32
	Height float32
}

type target struct {
	Name    string
	Display dimensions
}

type taggedImage struct {
	Name  string
	Image image.Image
}

func getRatio(pixelResolution dimensions) ratio {
	var width, height float32
	if pixelResolution.Width > pixelResolution.Height {
		width = 1
		height = float32(pixelResolution.Height) / float32(pixelResolution.Width)

	} else if pixelResolution.Height > pixelResolution.Width {
		height = 1
		width = float32(pixelResolution.Width) / float32(pixelResolution.Height)
	} else {
		width = 1
		height = 1
	}
	return ratio{
		Width:  width,
		Height: height,
	}
}

func getCropDimension(imageDimensions dimensions, cropRatio ratio) dimensions {
	return dimensions{
		Width:  uint(float32(imageDimensions.Width) * cropRatio.Width),
		Height: uint(float32(imageDimensions.Height) * cropRatio.Height),
	}
}

// pathExists returns whether the given file or directory exists
// https://stackoverflow.com/a/10510783
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func main() {

	imageNames := []string{
		"Spicy_Dark",
		"Spicy_Dark_solo",
	}

	targetDimensions := []target{
		target{
			Name: "iPhone 8, 7, 6, 6s",
			Display: dimensions{
				Width:  750,
				Height: 1334,
			},
		},

		target{
			Name: "iPhone 8 Plus, 7 Plus, 6 Plus, 6s Plus",
			Display: dimensions{
				Width:  1242,
				Height: 2208,
			},
		},

		target{
			Name: "iPhone Xr",
			Display: dimensions{
				Width:  828,
				Height: 1792,
			},
		},

		target{
			Name: "iPhone Xs",
			Display: dimensions{
				Width:  1125,
				Height: 2436,
			},
		},
		target{
			Name: "iPhone Xs",
			Display: dimensions{
				Width:  1242,
				Height: 2688,
			},
		},
		target{
			Name: "HD Landscape",
			Display: dimensions{
				Width:  1280,
				Height: 720,
			},
		},
		target{
			Name: "HD Portrait",
			Display: dimensions{
				Width:  720,
				Height: 1280,
			},
		},
		target{
			Name: "FHD Landscape",
			Display: dimensions{
				Width:  1920,
				Height: 1080,
			},
		},
		target{
			Name: "FHD Portrait",
			Display: dimensions{
				Width:  1440,
				Height: 2560,
			},
		},
		target{
			Name: "QHD-WQHD Landscape",
			Display: dimensions{
				Width:  2560,
				Height: 1440,
			},
		},
		target{
			Name: "QHD-WQHD Portrait",
			Display: dimensions{
				Width:  1080,
				Height: 1920,
			},
		},
		target{
			Name: "QHD 21:9 Ultrawide Landscape",
			Display: dimensions{
				Width:  3440,
				Height: 1440,
			},
		},
		target{
			Name: "4K Landscape",
			Display: dimensions{
				Width:  3840,
				Height: 2160,
			},
		},
		target{
			Name: "4K Portrait",
			Display: dimensions{
				Width:  2160,
				Height: 3840,
			},
		},
	}

	if exists, _ := pathExists("./resized/"); exists {
		os.RemoveAll("./resized/")
	}

	for _, folderName := range targetDimensions {
		os.MkdirAll("./resized/"+folderName.Name+"/", os.ModePerm)
	}

	images := []taggedImage{}

	for _, imageName := range imageNames {
		imageData, err := ioutil.ReadFile(imageName + ".png")
		if err != nil {
			log.Println(err)
			continue
		}
		imageDecoded, _, err := image.Decode(bytes.NewReader(imageData))
		image := taggedImage{
			Name:  imageName,
			Image: imageDecoded,
		}
		if err != nil {
			log.Println("Image " + imageName + " does not seem to be valid:")
			log.Println(err)
			continue
		}
		images = append(images, image)
	}

	for _, targetDevice := range targetDimensions {

		for _, workingImage := range images {

			imageDimensions := dimensions{
				Width:  uint(workingImage.Image.Bounds().Dx()),
				Height: uint(workingImage.Image.Bounds().Dy()),
			}

			cropDimensions := getCropDimension(imageDimensions, getRatio(targetDevice.Display))

			croppedImage, err := cutter.Crop(workingImage.Image, cutter.Config{
				Width:  int(cropDimensions.Width),
				Height: int(cropDimensions.Height),
				Mode:   cutter.Centered,
			})

			if err != nil {
				log.Println("Image could not be cropped.")
				log.Println(err)
				continue
			}

			/*
				croppedImageName := "./cropped/" + workingImage.Name + " (" + targetDevice.Name + ").png"

				outputFileCropped, err := os.Create(croppedImageName)

				if err != nil {
					log.Println("Image output could not be created.")
					log.Println(err)
					continue
				}

				png.Encode(outputFileCropped, croppedImage)

				outputFileCropped.Close()
			*/

			resizedImage := resize.Resize(
				targetDevice.Display.Width,
				targetDevice.Display.Height,
				croppedImage,
				resize.Lanczos3,
			)

			if err != nil {
				log.Println("Image could not be resized.")
				log.Println(err)
				continue
			}

			resizedImageName := "./resized/" + targetDevice.Name + "/" + workingImage.Name + ".png"

			outputFileResized, err := os.Create(resizedImageName)

			if err != nil {
				log.Println("Image output could not be created.")
				log.Println(err)
				continue
			}

			png.Encode(outputFileResized, resizedImage)

			outputFileResized.Close()

		}
	}

}
