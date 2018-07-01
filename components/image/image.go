package image

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"

	"github.com/Casxt/TimeLine/session"

	"github.com/Casxt/TimeLine/page"
)

//Route decide page
func Route(res http.ResponseWriter, req *http.Request) {
	var result []byte
	var status int
	//subPath := req.URL.Path[len("/image"):]
	switch {
	case req.Method == "POST":
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

	Session, _ := session.Auto(res, req)
	UserID, ok := Session.Get("UserID")

	PostReader, err := req.MultipartReader()
	if err != nil {
		//log.Println("UploadImage-MultipartReader:", err.Error())
		byteRes, _ = json.Marshal(ImgUploadRes{
			State: "Failde",
			Msg:   err.Error(),
			Hashs: make([]string, 0),
		})
		return 400, byteRes
	}

	MaxSize := int64(20 * 1024 * 1024)
	rawBuff := bytes.NewBuffer(make([]byte, MaxSize))
	Hashlist := make([]string, 9)
	Imglist := make([][]byte, 9)
	imgNum := 0
	for {
		rawBuff.Reset()
		part, err := PostReader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			byteRes, _ = json.Marshal(ImgUploadRes{
				State: "Failde",
				Msg:   err.Error(),
			})
			return 400, byteRes
		}
		if imgNum >= 9 {
			byteRes, _ = json.Marshal(ImgUploadRes{
				State: "Failde",
				Msg:   "Too many img",
			})
			return 400, byteRes
		}

		if part.FormName() != "images" {
			continue
		}

		rawSize, err := io.CopyN(rawBuff, part, MaxSize+1)
		//DetectContentType need first 512 byte,
		//in order to be safe, keep file large more than 512 bytes
		if rawSize < 512 {
			byteRes, _ = json.Marshal(ImgUploadRes{
				State: "Failde",
				Msg:   "File too small",
			})
			return 400, byteRes
		}
		if MaxSize -= rawSize; MaxSize < 0 {
			byteRes, _ = json.Marshal(ImgUploadRes{
				State: "Failde",
				Msg:   "File too big",
			})
			return 400, byteRes
		}

		if err != nil && err != io.EOF {
			byteRes, _ = json.Marshal(ImgUploadRes{
				State: "Failde",
				Msg:   err.Error(),
			})
			return 400, byteRes
		}
		//ReEncode img
		//Code before already make sure that file large than 512 bytes
		var img image.Image
		rawBytes := rawBuff.Bytes()
		switch http.DetectContentType(rawBytes) {
		case "image/jpeg":
			img, err = jpeg.Decode(bytes.NewReader(rawBytes))
			if err != nil {
				log.Println(err.Error())
				byteRes, _ = json.Marshal(ImgUploadRes{
					State: "Failde",
					Msg:   "invalid jpeg file",
				})
				return 400, byteRes
			}
		case "image/png":
			img, err = png.Decode(bytes.NewReader(rawBytes))
			if err != nil {
				log.Println(err.Error())
				byteRes, _ = json.Marshal(ImgUploadRes{
					State: "Failde",
					Msg:   "invalid jpeg file",
				})
				return 400, byteRes
			}
		default:
			byteRes, _ = json.Marshal(ImgUploadRes{
				State: "Failde",
				Msg:   "Unsupprot Format",
			})
			return 400, byteRes
		}

		rawBuff.Reset()
		jpeg.Encode(rawBuff, img, &jpeg.Options{Quality: 80})
		//Check img dumplicate
		JpgBytes := make([]byte, rawBuff.Len())
		rawBuff.Read(JpgBytes)
		Hash256 := sha256.New()
		Hash256.Write(JpgBytes)
		ImgHash := hex.EncodeToString(Hash256.Sum(nil))

		Imglist[imgNum] = JpgBytes
		Hashlist[imgNum] = ImgHash
		imgNum++
	}

	//after all img pass check, can they be storge
	for i := 0; i < imgNum; i++ {
		//storge img
		if err = page.SaveFile(Imglist[i], "static", "image", Hashlist[i]+".jpg"); err != nil {
			byteRes, _ = json.Marshal(ImgUploadRes{
				State: "Failde",
				Msg:   err.Error(),
			})
			return 400, byteRes
		}
	}

	byteRes, _ = json.Marshal(ImgUploadRes{
		State: "Success",
		Msg:   "upload success",
		Hashs: Hashlist[0:imgNum],
	})
	return 200, byteRes
}
