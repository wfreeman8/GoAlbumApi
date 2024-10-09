package models

import (
  "fmt"
  "errors"
  "os"
  "io"
  "reflect"
	"mime/multipart"
  "encoding/json"
  "strings"
  "path"
  "regexp"
)

type ImageResize struct {
  Name string `json:"name"`
  Width int `json:"width"`
  Height int `json:"height"`
  Crop bool `json:"crop"`
}

type Album struct {
  albumFolderPath string
  AlbumId string `json:"albumid"`
  Title string `json:"title"`
  Author string `json:"author"`
  Description string `json:"description"`
  DateTaken string `json:"created_datetime"`
  Pagename string `json:"pagename"`
  FeaturedImage string `json:"featured_image"`
  ThumbnailSize string `json:"thumbnail_size"`
  LargeSize string `json:"large_size"`
  Resizes []ImageResize `json:"resizes"`
  Images []Image `json:"images"`
}

type AlbumFormatted struct {
  AlbumId string `json:"albumid"`
  Title string `json:"title"`
  Author string `json:"author"`
  Description string `json:"description"`
  DateTaken string `json:"created_datetime"`
  Pagename string `json:"pagename"`
  FeaturedImage string `json:"featured_image"`
  ThumbnailSize string `json:"thumbnail_size"`
  LargeSize string `json:"large_size"`
  Resizes []ImageResize `json:"resizes"`
  Images []ImageFormatted `json:"images"`
}


func FindAlbum(albumBasePath string, albumPagename string) (Album, error) {
  albumFolderPath := albumBasePath + "/" + albumPagename
  albumDirMeta, err := os.Stat(albumFolderPath)
  if err != nil {
    return  Album{}, err
  }

  if albumDirMeta.IsDir() {
    albumJsonFile, err := os.Open(albumFolderPath + "/album.json")
    defer albumJsonFile.Close()
    if err != nil {
      return  Album{}, err
    }

    albumJsonContent, _ := io.ReadAll(albumJsonFile)
    album := Album{
      albumFolderPath: albumFolderPath,
    }
    err = json.Unmarshal(albumJsonContent, &album)
    if err == nil {
      return album, nil
    }
  }

  return Album{}, errors.New("Album not Found")
}

func (album *Album) SetAlbumPath(albumPath string) {
  if album.albumFolderPath == "" {
    album.albumFolderPath = albumPath
  }
}


// Delete resized images, optionally delete only images of a particular image
func (album *Album) DeleteResizeImages(imageBasename string) {
  albumDirMeta, err := os.Stat(album.albumFolderPath)
  resizeImageRegex := regexp.MustCompile(`-\d+x\d+x\d.jpg$`)
  if err != nil {
    return
  }

  if albumDirMeta.IsDir() {
    albumDirEntries, err := os.ReadDir(album.albumFolderPath)
    if err != nil {
      return
    }

    for _, albumDirEntry := range albumDirEntries {
      entryExt := path.Ext(albumDirEntry.Name())
      if strings.ToLower(entryExt) == ".jpg" {
        if imageBasename != "" {
          if strings.HasPrefix(albumDirEntry.Name(), imageBasename + "-") == false {
            continue
          }
        }

        imagePath := album.albumFolderPath + "/" + albumDirEntry.Name()
        if resizeImageRegex.MatchString(albumDirEntry.Name()) {
          os.Remove(imagePath)
        }
      }
    }
  }
}

func (album *Album) CheckAndResetImagesIndex() (bool, error) {
  albumDirMeta, err := os.Stat(album.albumFolderPath)

  if err != nil {
    return false, err
  }

  if albumDirMeta.IsDir() {
    albumDirEntries, err := os.ReadDir(album.albumFolderPath)
    if err != nil {
      return false, errors.New("Invalid directory for album")
    }

    dirImageList := []string{}
    for _, albumDirEntry := range albumDirEntries {
      dirEntryExt := strings.ToLower(path.Ext(albumDirEntry.Name()))
      if !strings.Contains(albumDirEntry.Name(), "-") && dirEntryExt == ".jpg" {
        dirImageList = append(dirImageList, albumDirEntry.Name())
      }
    }
    jsonImageList := []string{}
    for _, imageMeta := range album.Images {
      jsonImageList = append(jsonImageList, path.Base(imageMeta.FilePath))
    }

    if reflect.DeepEqual(dirImageList, jsonImageList) == false {
      album.DeleteResizeImages("")
      album.IndexImagesWithDirectory()
      return true, nil
    }
  }

  return false, nil
}

func (album *Album) IndexImagesWithDirectory() {
  var albumImagesArray []Image
  albumDirEntries, _ := os.ReadDir(album.albumFolderPath)

  // ensures that JPG extension is lower case
  for _, albumDirEntry := range albumDirEntries {
    entryExt := path.Ext(albumDirEntry.Name())
    if entryExt == ".JPG" {
      imagePath := album.albumFolderPath + "/" + albumDirEntry.Name()

      // because golang won't rename just for case change must rename twice to workaround
      tmpImagePath := album.albumFolderPath + "/" + strings.ToLower(albumDirEntry.Name()) + ".tmp"
      err := os.Rename(imagePath, tmpImagePath)
      err = os.Rename(tmpImagePath, album.albumFolderPath + "/" + strings.ToLower(albumDirEntry.Name()))
      if err != nil {
        panic(err)
      }
    }
  }

  // reset directory index for resorting correctly after filename lowercase
  albumDirEntries, _ = os.ReadDir(album.albumFolderPath)
  for _, albumDirEntry := range albumDirEntries {
    entryExt := path.Ext(albumDirEntry.Name())
    if entryExt == ".jpg" {
      imagePath := album.albumFolderPath + "/" + albumDirEntry.Name()
      imageMetaData := getImageDataByPath(imagePath)
      if imageMetaData != (Image{}) {
        albumImagesArray = append(albumImagesArray, imageMetaData)
      }
    }
  }
  album.Images = albumImagesArray
}

func (album *Album) SaveUploadedImage(imageFile *multipart.FileHeader) (error) {
  fileExt := path.Ext(imageFile.Filename)
  newFilename := strings.ToLower(strings.TrimSuffix(imageFile.Filename, fileExt))
  cleanFilenameRe := regexp.MustCompile(`[^a-z0-9]`)
  newFilename = cleanFilenameRe.ReplaceAllString(newFilename, "")
  if len(newFilename) < 5 {
    return errors.New("Filename is too short, must be greater than 5 characters")
  }

  newImagePath := album.albumFolderPath + "/" + newFilename + ".jpg"

  // Save Image to album folder and add to album JSON
  newImage, err := SaveImage(imageFile, newImagePath)
  if err != nil {
    return err
  }

  currentImageEntry, _ := album.RetrieveImage(path.Base(newImage.FilePath))

  if currentImageEntry.Title != "" {
    currentImageEntry.Title = newImage.Title
    currentImageEntry.Width = newImage.Width
    currentImageEntry.Height = newImage.Height

    // delete old resized images
    album.DeleteResizeImages(newFilename)
  } else {
    album.Images = append(album.Images, newImage)
  }
  album.Save()
  return nil

}

func (album *Album) RetrieveImage(imageFilename string) (*Image, error) {
  if imageFilename != "" {
    for idx, imageData := range album.Images {
      if path.Base(imageData.FilePath) == imageFilename {
        return &album.Images[idx], nil
      }
    }
  }
  return &Image{}, errors.New("Image file does not exist")
}

func (album *Album) FindResizeByName(resizeName string) (ImageResize, error) {
  for _, resize := range album.Resizes {
    if resize.Name == resizeName {
      return resize, nil
    }
  }
  return ImageResize{}, errors.New("Image Size does not exist")
}

func (album *Album) FindResizeBySize(width int, height int, crop bool) (ImageResize, error) {
  for _, resize := range album.Resizes {
    if resize.Width == width && resize.Height == height && resize.Crop == crop {
      return resize, nil
    }
  }
  return ImageResize{}, errors.New("Image Size does not exist")
}

func (album *Album) Save() (error) {
  albumJsonPath := album.albumFolderPath + "/album.json"
  _, err := os.Stat(albumJsonPath)
  albumJsonFile := new(os.File)
  var openFileErr error
  if errors.Is(err, os.ErrNotExist) {
    albumJsonFile, openFileErr = os.Create(albumJsonPath)
  } else {
    albumJsonFile, openFileErr = os.OpenFile(albumJsonPath, os.O_TRUNC | os.O_WRONLY, 0644)
  }
  if openFileErr != nil {
    return openFileErr
  }

  albumJson, jsonErr := json.Marshal(album)
  if jsonErr != nil {
    return jsonErr
  }
  albumJsonFile.Write(albumJson)
  albumJsonFile.Close()
  fmt.Println("album saved")
  return nil
}

func (album *Album) GetAlbumFormatted() (AlbumFormatted) {
  var formattedImagesArr []ImageFormatted
  for _, albumImg := range album.Images {
    formattedImagesArr = append(formattedImagesArr, albumImg.GetImageFormatted())
  }

  return AlbumFormatted{
    AlbumId: album.AlbumId,
    Title: album.Title,
    Author: album.Author,
    Description: album.Description,
    DateTaken: album.DateTaken,
    Pagename: album.Pagename,
    FeaturedImage: album.FeaturedImage,
    ThumbnailSize: album.ThumbnailSize,
    LargeSize: album.LargeSize,
    Resizes: album.Resizes,
    Images: formattedImagesArr,
  }
}
