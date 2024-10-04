# GoAlbumApi API Reference

## Configuration Endpoints

GET /config returns meta information about the gallery

```
GET /config


**response**

200 OK
{
  "gallery_title": "Parades Gallery"
}

```

POST /config is only accessible for initial creation of the config.json and is disabled afterwards.

```
POST /config
Content-Type: application/json

{
  "password": "TestPassword"
  "albums_base_path": "albums"
  "gallery_title": "Parades Gallery"
}

**response**

201 Created
{
  "success": true
}

```


PUT /config enables updating the password and gallery_title. Must provide the current password in the Authorization header to succeed

```
PUT /config
Content-Type: application/json
Authorization: TestPassword

{
  "password": "TestPassword*123"
  "gallery_title": "Parades Gallery"
}

200 OK
{
  "success": true
}

```

## Albums Endpoints

GET /albums provides a summary list of all available albums.

```
GET /albums

**response**

200 OK
[
  {
    "title": "Example Album",
    "pagename": "example",
    "description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. ",
    "imagecount": 1,
    "author": "Example"
  }
]
```

POST /albums creates a new album in the gallery. Not all fields are required. pagename field is the Album's URL slug for easy access to the album.

```
POST /albums
Content-Type: application/json
Authorization: TestPassword
{
  "title": "Burlingame Pets Parade",
  "pagename": "petsparade",
  "author": "TJ",
  "created_datetime": "2024-09-28 10:00:00",

  "resizes": [
    {
      "Name": "thumbnail",
      "Width": 500,
      "Height": 500,
      "Crop": true
    },
    {
      "Name": "large",
      "Width": 1500,
      "Height": 1000,
      "Crop": false
    },
    {
      "Name": "larger",
      "Width": 2250,
      "Height": 1500,
      "Crop": false
    }
  ],
  "description": "Watch the pets from Burlingame community parade along Burlingame Ave."
}


**response**

201 Created
{
  "albumid": "175152a7-04f3-4ef0-b322-3172fca9f534",
  "title": "Pets Parade",
  "author": "TJ",
  "description": "Watch the pets from Burlingame community parade along Burlingame Ave.",
  "created_datetime": "2024-09-28 10:00:00",
  "pagename": "petsparade",
  "featured_image": "",
  "thumbnail_size": "",
  "large_size": "",
  "resizes": [
    {
      "name": "thumbnail",
      "width": 500,
      "height": 500,
      "crop": true
    },
    {
      "name": "large",
      "width": 1500,
      "height": 1000,
      "crop": false
    },
    {
      "name": "larger",
      "width": 2250,
      "height": 1500,
      "crop": false
    }
  ],
  "images": null
}
```

## Album Endpoints

GET  /album/{pagename} retrieves the full album meta data

```
GET /album/petsparade

**response**

200 OK
{
  "albumid": "175152a7-04f3-4ef0-b322-3172fca9f534",
  "title": "Pets Parade",
  "author": "TJ",
  "description": "Watch the pets from Burlingame community parade along Burlingame Ave.",
  "created_datetime": "2024-09-28 10:00:00",
  "pagename": "petsparade",
  "featured_image": "",
  "thumbnail_size": "",
  "large_size": "",
  "resizes": [
    {
      "name": "thumbnail",
      "width": 500,
      "height": 500,
      "crop": true
    },
    {
      "name": "large",
      "width": 1500,
      "height": 1000,
      "crop": false
    },
    {
      "name": "larger",
      "width": 2250,
      "height": 1500,
      "crop": false
    }
  ],
  "images": [
    {
      "filename": "8k1a8388.jpg",
      "title": "8K1A8388.JPG",
      "width": 3360,
      "height": 2240
    }
  ]
}
```

PUT  /album/{pagename} updates the album metadata. If resizes are changed all cached image resizes are deleted. The album pagename and images field can not be altered

```
PUT /album/petsparade
Content-Type: application/json
Authorization: TestPassword

{
  "title": "Pets Parade",
  "author": "Anonymous",
  "description": "Watch 100s of pets from the Burlingame community march along Burlingame Ave.",
  "created_datetime": "2024-09-28 10:00:00",
  "featured_image": "",
  "thumbnail_size": "",
  "large_size": "",
  "resizes": [
    {
      "name": "thumbnail",
      "width": 500,
      "height": 500,
      "crop": true
    },
    {
      "name": "large",
      "width": 1500,
      "height": 1000,
      "crop": false
    },
    {
      "name": "larger",
      "width": 2250,
      "height": 1500,
      "crop": false
    }
  ]
}


**response**

200 OK
{
  "albumid": "175152a7-04f3-4ef0-b322-3172fca9f534",
  "title": "Pets Parade",
  "author": "Anonymous",
  "description": "Watch 100s of pets from the Burlingame community march along Burlingame Ave.",
  "created_datetime": "2024-09-28 10:00:00",
  "pagename": "petsparade",
  "featured_image": "",
  "thumbnail_size": "",
  "large_size": "",
  "resizes": [
    {
      "name": "thumbnail",
      "width": 500,
      "height": 500,
      "crop": true
    },
    {
      "name": "large",
      "width": 1500,
      "height": 1000,
      "crop": false
    },
    {
      "name": "larger",
      "width": 2250,
      "height": 1500,
      "crop": false
    }
  ],
  "images": [
    {
      "filename": "8k1a8388.jpg",
      "title": "8K1A8388.JPG",
      "width": 3360,
      "height": 2240
    }
  ]
}
```

## Album Images

POST  /album/{album_pagename}
/images allows uploading an image to the album


```
POST  /album/petsparade/images
multipart/form-data; boundary=----WebKitFormBoundaryIBOjf86V3WxcaBql
Authorization: TestPassword

----WebKitFormBoundaryIBOjf86V3WxcaBql
{binary_image_data}


----WebKitFormBoundaryIBOjf86V3WxcaBql

**response**

200 OK
[
  {
    "filename": "8k1a8388.jpg",
    "title": "8K1A8388.JPG",
    "width": 3360,
    "height": 2240
  }
]
```

## Album Images

GET /img/{album_pagename}/{imagename}-{size_name}.jpg retries the images resized according the resized image

```
POST  /img/petsparade/8k1a8388-large.jpg

**response**

200 OK
{binary_image_data}
```


