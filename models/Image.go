package models

import(
  "os"
  "errors"
  "io/ioutil"
  "log"
  "fmt"
  "math"
  "image"
  "image/jpeg"
  "image/png"
	"mime/multipart"
  "path"
  "strings"
  "golang.org/x/image/draw"

)

type Image struct {
  FilePath string `json:"filepath"`
  Title string `json:"title"`
  Width int `json:"width"`
  Height int `json:"height"`
}

type ImageFormatted struct {
  FileName string `json:"filename"`
  Title string `json:"title"`
  Width int `json:"width"`
  Height int `json:"height"`
}


func getImageDataByPath(imageFilePath string) (Image) {
  imageDataFile, _ := os.Open(imageFilePath)
  defer imageDataFile.Close()

  jpegSource, err := jpeg.Decode(imageDataFile)

  imageMetaData := Image{}
  if err != nil {
    return imageMetaData
  }
  imageMetaData.FilePath = imageDataFile.Name()
  imageMetaData.Title = path.Base(imageDataFile.Name())
  imageMetaData.Height = jpegSource.Bounds().Max.Y
  imageMetaData.Width = jpegSource.Bounds().Max.X
  return imageMetaData
}

func SaveImage(imageFile *multipart.FileHeader, filePath string) (Image, error) {
  imageType := strings.ToLower(path.Ext(imageFile.Filename))
  var err error
  if imageType == ".jpg" || imageType == ".png" {
    imageFileData, err := imageFile.Open()
    defer imageFileData.Close()

    var imageSource image.Image
    if err == nil {
      if imageType == ".jpg" {
        imageSource, err = jpeg.Decode(imageFileData)
      } else if imageType == ".png" {
       imageSource, err = png.Decode(imageFileData)
      } else {
        return Image{}, errors.New("Invalid image file type. must be png or jpeg")
      }
      
      if err == nil {
        newImageFs, err := os.Create(filePath)
        defer newImageFs.Close()
        if err != nil {
          return Image{}, err
        }
        err = jpeg.Encode(newImageFs, imageSource, &jpeg.Options{Quality: 100})
        if err == nil {
          imageMetaData := Image{}
          imageMetaData.FilePath = filePath
          imageMetaData.Title = path.Base(imageFile.Filename)
          imageMetaData.Height = imageSource.Bounds().Max.Y
          imageMetaData.Width = imageSource.Bounds().Max.X
          return imageMetaData, nil
        }
        

      }
    }

  } else {
    err = errors.New("Image is invalid format")
  }
  return Image{}, err
}

func (albumImage *Image) GetBytes() ([]byte) {
  fs, err := ioutil.ReadFile(albumImage.FilePath)
  if err != nil {
    log.Fatal(err)
  }
  return fs
}

func (albumImage *Image) GetResizeBytes(resize ImageResize) ([]byte) {

  newWidth := resize.Width
  newHeight := resize.Height
  crop := resize.Crop
  
  imageSourceFs, err := os.Open(albumImage.FilePath)
  defer imageSourceFs.Close()
  if err != nil {
    log.Fatal(err)
  }

  albumImagePath := path.Dir(albumImage.FilePath)
  imageFilename := path.Base(albumImage.FilePath)
  imageExt := path.Ext(imageFilename)
  imageFilenameBase := strings.TrimSuffix(imageFilename, imageExt)
  cropStr := "0"
  if crop {
    cropStr = "1"
  }

  resizedImageFilename := fmt.Sprintf("%s-%dx%dx%s%s", imageFilenameBase, newWidth, newHeight, cropStr, imageExt)
  resizedImageFilePath := path.Join(albumImagePath, resizedImageFilename)
  newImageBytes, err := ioutil.ReadFile(resizedImageFilePath)
  if err == nil {
    return newImageBytes
  }

  newImageFs, _ := os.Create(resizedImageFilePath)
  defer newImageFs.Close()

  imageData, _ := jpeg.Decode(imageSourceFs)
  orgWidth := float64(imageData.Bounds().Max.X)
  orgHeight := float64(imageData.Bounds().Max.Y)

  var newImage draw.Image

  orgImageRatio :=  orgWidth / orgHeight
  newImageRatio := float64(newWidth) / float64(newHeight)
  orgImageRect := imageData.Bounds()
  if crop {
    if orgImageRatio > newImageRatio {
      newImageWidth := float64(newHeight) * orgImageRatio
      orgImageWidthOffset := (newImageWidth - float64(newWidth)) / 2.0
      widthRatio := orgWidth / newImageWidth
      orgImageWidthCut := int(math.Round(orgImageWidthOffset * widthRatio))

      orgXCenterLeft := orgImageWidthCut
      orgXCenterRight := int(orgWidth) - orgImageWidthCut
      orgImageRect = image.Rect(orgXCenterLeft, 0  , orgXCenterRight, imageData.Bounds().Max.Y)

    } else {
      orgImageHeightRatio := 1.0 / orgImageRatio
      newImageHeight := float64(newWidth) * orgImageHeightRatio

      orgImageHeightOffset := (newImageHeight - float64(newHeight)) / 2.0
      heightRatio := orgHeight / newImageHeight

      orgImageHeightCut := int(math.Round(orgImageHeightOffset * heightRatio))
      orgXMiddleTop := orgImageHeightCut
      orgXMiddleBottom := int(orgHeight) - orgImageHeightCut
      orgImageRect = image.Rect(0, orgXMiddleTop, imageData.Bounds().Max.X, orgXMiddleBottom)
    }
    
  } else {
    if orgImageRatio > newImageRatio {
      newHeight = int(math.Round((orgHeight / orgWidth) * float64(newWidth)))
    } else {
      newWidth = int(math.Round(orgImageRatio * float64(newHeight)))
    }
  }
  newImage = image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
  draw.NearestNeighbor.Scale(newImage, newImage.Bounds(), imageData, orgImageRect, draw.Over, nil)

  jpeg.Encode(newImageFs, newImage, &jpeg.Options{Quality: 70})
  newImageFs.Close()

  imageBytes, err := ioutil.ReadFile(resizedImageFilePath)
  if err == nil {
    return imageBytes
  } else {
    fmt.Println("err-------")
    fmt.Println(err)
  }
  return []byte{}
}

func (albumImage *Image) GetImageFormatted() ImageFormatted {
  return ImageFormatted {
    FileName: path.Base(albumImage.FilePath),
    Title: albumImage.Title,
    Width: albumImage.Width,
    Height: albumImage.Height,
  }
}