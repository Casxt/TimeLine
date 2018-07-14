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

	"github.com/Casxt/TimeLine/database"
	"github.com/Casxt/TimeLine/static"
	"github.com/Casxt/TimeLine/tools"
)

//Route decide static
func Route(res http.ResponseWriter, req *http.Request) {
	var result []byte
	var status int
	subPath := req.URL.Path[len("/image"):]
	switch {
	case req.Method == "POST":
		status, result = UploadImage(res, req)
		//match /image/8c2c78eb41c26bc571a004895427300c187c8f2c2f3c0600a73773b685a8ee0c
	case tools.CheckImgHash(subPath):
		status, result = GetImage(res, req)
	default:
		status, result, _ = static.GetPage("components", "image", "line.html")
	}
	res.WriteHeader(status)
	res.Write(result)
}

//GetImage retuen img if user have permission
func GetImage(res http.ResponseWriter, req *http.Request) (status int, byteRes []byte) {

	UserID, _ := tools.GetLoginStateOfCookie(req)
	if UserID == "" {
		resByte, _ := json.Marshal(map[string]string{
			"State": "Failde",
			"Msg":   "User Not Sign In",
		})
		return 400, resByte
	}
	//imgName = 8c2c78eb41c26bc571a004895427300c187c8f2c2f3c0600a73773b685a8ee0c
	imgName := req.URL.Path[len("/image/") : len("/image/")+64]
	_, err := database.CheckImgInfo(imgName, UserID)
	if err != nil {
		resByte, _ := json.Marshal(map[string]string{
			"State": "Failde",
			"Msg":   "User Not Have This Img",
		})
		return 400, resByte
	}
	status, byteRes, _ = static.GetFile("static", "image", imgName+".jpg")
	return status, byteRes
}

//UploadImage will procss img uploaded, resize and caculate hash and storge
// max size is 20MB
// TODO: Limit img num of user
func UploadImage(res http.ResponseWriter, req *http.Request) (status int, byteRes []byte) {
	type ImgUploadRes struct {
		State string
		Msg   string
		Hashs []string
	}

	//CheckUser
	UserID, _ := tools.GetLoginStateOfCookie(req)
	//Check User Login At first
	//For Some Reason, client Cannot receive this msg
	//I don't know why
	//TODO: fix this bug
	if UserID == "" {
		byteRes, _ = json.Marshal(ImgUploadRes{
			State: "Failde",
			Msg:   "User Not SignIn",
		})
		return 400, byteRes
	}

	PostReader, err := req.MultipartReader()
	if err != nil {
		//log.Println("UploadImage-MultipartReader:", err.Error())
		byteRes, _ = json.Marshal(ImgUploadRes{
			State: "Failde",
			Msg:   err.Error(),
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
			log.Println(err.Error())
			return 400, tools.JSONMarshal(ImgUploadRes{
				State: "Failde",
				Msg:   "File too small",
			})
		}

		if part.FormName() != "images" {
			continue
		}

		if imgNum >= 9 {
			return 400, tools.JSONMarshal(ImgUploadRes{
				State: "Failde",
				Msg:   "Too many img",
			})
		}

		rawSize, err := io.CopyN(rawBuff, part, MaxSize+1)
		//DetectContentType need first 512 byte,
		//in order to be safe, keep file large more than 512 bytes
		if rawSize < 512 {
			return 400, tools.JSONMarshal(ImgUploadRes{
				State: "Failde",
				Msg:   "File too small",
			})
		}
		if MaxSize -= rawSize; MaxSize < 0 {
			return 400, tools.JSONMarshal(ImgUploadRes{
				State: "Failde",
				Msg:   "File too big",
			})
		}

		if err != nil && err != io.EOF {
			log.Println(err.Error())
			return 400, tools.JSONMarshal(ImgUploadRes{
				State: "Failde",
				Msg:   "Unknow Error",
			})
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
				return 400, tools.JSONMarshal(ImgUploadRes{
					State: "Failde",
					Msg:   "invalid jpeg file",
				})
			}
		case "image/png":
			img, err = png.Decode(bytes.NewReader(rawBytes))
			if err != nil {
				log.Println(err.Error())
				return 400, tools.JSONMarshal(ImgUploadRes{
					State: "Failde",
					Msg:   "invalid jpeg file",
				})
			}
		default:
			return 400, tools.JSONMarshal(ImgUploadRes{
				State: "Failde",
				Msg:   "Unsupprot Format",
			})
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

	database.CreateImage(UserID, Hashlist[0:imgNum])

	//after all img pass check, can they be storge
	for i := 0; i < imgNum; i++ {
		//storge img
		if err = static.SaveFile(Imglist[i], "static", "image", Hashlist[i]+".jpg"); err != nil {
			log.Println(err.Error())
			return 400, tools.JSONMarshal(ImgUploadRes{
				State: "Failde",
				Msg:   "File Storge Failed",
			})
		}
	}

	byteRes, _ = json.Marshal(ImgUploadRes{
		State: "Success",
		Msg:   "upload success",
		Hashs: Hashlist[0:imgNum],
	})
	return 200, byteRes
}
