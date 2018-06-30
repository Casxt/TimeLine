package image

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"

	"github.com/Casxt/TimeLine/page"
)

//Route decide page
func Route(res http.ResponseWriter, req *http.Request) {
	var result []byte
	var status int
	//subPath := req.URL.Path[len("/index"):]
	switch {
	case req.Method == "POST":
		var jsonRes map[string]string
		status, result = UploadImage(res, req)
	default:
		result, status, _ = page.GetPage("components", "image", "line.html")
	}
	res.WriteHeader(status)
	res.Write(result)
}

//UploadImage will procss img uploaded, resize and caculate hash and storge
// max size is 20MB
func UploadImage(res http.ResponseWriter, req *http.Request) (status int, byteRes []byte) {
	type ImgUploadRes struct {
		State string
		Msg   string
		Hashs []string
	}

	MaxSize := 20 * 1024 * 1024
	PostReader, err := req.MultipartReader()
	if err != nil {
		log.Println(err.Error())
		byteRes, _ = json.Marshal(ImgUploadRes{
			State: "Failde",
			Msg:   err.Error(),
			Hashs: make([]string, 0),
		})
		return 400, byteRes
	}
	//TODO 考虑复用？
	rawBytes := make([]byte, 10*1024*1024)
	for {
		part, err := PostReader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			byteRes, _ = json.Marshal(ImgUploadRes{
				State: "Failde",
				Msg:   err.Error(),
				Hashs: make([]string, 0),
			})
			return 400, byteRes
		}

		rawSize, err := part.Read(rawBytes)
		//DetectContentType need first 512 byte,
		//in order to be safe, keep file large more than 512 bytes
		if rawSize < 512 {
			byteRes, _ = json.Marshal(ImgUploadRes{
				State: "Failde",
				Msg:   "File too small",
				Hashs: make([]string, 0),
			})
			return 400, byteRes
		}
		if MaxSize -= rawSize; MaxSize < 0 {
			byteRes, _ = json.Marshal(ImgUploadRes{
				State: "Failde",
				Msg:   "File too big",
				Hashs: make([]string, 0),
			})
			return 400, byteRes
		}
		if err != nil {
			byteRes, _ = json.Marshal(ImgUploadRes{
				State: "Failde",
				Msg:   err.Error(),
				Hashs: make([]string, 0),
			})
			return 400, byteRes
		}
		//Code before already make sure that file large than 512 bytes
		var img image.Image
		switch http.DetectContentType(rawBytes) {
		case "image/jpeg":
			img, err = jpeg.Decode(bytes.NewReader(rawBytes))
			if err != nil {
				byteRes, _ = json.Marshal(ImgUploadRes{
					State: "Failde",
					Msg:   "invalid jpeg file",
					Hashs: make([]string, 0),
				})
				return 400, byteRes
			}
		default:
			byteRes, _ = json.Marshal(ImgUploadRes{
				State: "Failde",
				Msg:   "Unsupprot Format",
				Hashs: make([]string, 0),
			})
			return 400, byteRes
		}
		buff := bytes.NewBuffer(rawBytes)
		buff.Reset()
		jpeg.Encode(buff, img, &jpeg.Options{Quality: 80})

		JpgBytes := buff.Bytes()
		Hash256 := sha256.New()
		Hash256.Write(JpgBytes)
		ImgHash := hex.EncodeToString(Hash256.Sum(nil))

		if err = page.SaveFile(JpgBytes, "components", "line", ImgHash+".jpg"); err != nil {
			byteRes, _ = json.Marshal(ImgUploadRes{
				State: "Failde",
				Msg:   err.Error(),
				Hashs: make([]string, 0),
			})
			return 400, byteRes
		}

	}

	byteRes, _ = json.Marshal(ImgUploadRes{
		State: "Success",
		Msg:   err.Error(),
		Hashs: make([]string, 0),
	})
	return 200, byteRes
}
