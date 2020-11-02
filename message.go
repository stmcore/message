package message

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/kpango/glg"
	"github.com/nfnt/resize"
)

type Elemental struct {
	Id        string `json:"Id"`
	Name      string `json:"Name"`
	Elemental string `json:"Elemental"`
}

type Message struct {
	ChName           string
	FileName         string
	Path             string
	TimeStamp        time.Time
	Colors           color.RGBA
	Transcoder       string
	DominantColorRef bool
}

type ConfigCh struct {
	UID     string
	Name    string
	Service string
}

type Titan struct {
	Name string `json:"Name"`
	UID  string `json:"UID"`
}

func logInfo(log string) {
	glg.Info(log)
}

var titanconf []ConfigCh
var elementalconf []ConfigCh
var mediaexcelconf []ConfigCh

var colorRef1 = color.RGBA{
	R: uint8(214),
	G: uint8(213),
	B: uint8(213),
}
var colorRef2 = color.RGBA{
	R: uint8(230),
	G: uint8(229),
	B: uint8(229),
}
var colorRef3 = color.RGBA{
	R: uint8(160),
	G: uint8(160),
	B: uint8(160),
}
var colorRef4 = color.RGBA{
	R: uint8(213),
	G: uint8(213),
	B: uint8(213),
}
var colorRef5 = color.RGBA{
	R: uint8(213),
	G: uint8(214),
	B: uint8(213),
}

func init() {

	file, _ := ioutil.ReadFile("./titan.json")
	err := json.Unmarshal([]byte(file), &titanconf)

	if err != nil {
		log.Fatalln(err)
	}

	// file, _ = ioutil.ReadFile("./elemental.json")
	// err = json.Unmarshal([]byte(file), &elementalconf)

	// if err != nil {
	// 	log.Fatalln(err)
	// }

	file, _ = ioutil.ReadFile("./mediaexcel.json")
	err = json.Unmarshal([]byte(file), &mediaexcelconf)

	if err != nil {
		log.Fatalln(err)
	}
}

func getNameChTitan(pid, titanNode string) Titan {
	var data Titan
	url := "http://10.18.40.73:8892/titan/" + titanNode + "/" + pid
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		glg.Info(err)
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			glg.Info(err)
		}
		json.Unmarshal(body, &data)
	}

	return data
}

func getNameChElemental(pid, elementalNode string) Elemental {
	var data Elemental
	url := "http://10.18.40.73:8891/" + "elemental/" + elementalNode + "/" + pid
	//resp, err := http.Get(url)

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)

	if err != nil {
		glg.Info(err)
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			glg.Info(err)
		}
		json.Unmarshal(body, &data)
	}

	return data

}

func (self *Message) SetNameDotJPG(oldname string, transcoder string) {

	self.FileName = oldname

	if transcoder == "media_excel" {
		// newname := strings.Split(oldname, "_")
		// self.FileName = "mediaexcel_" + newname[len(newname)-1]
		// self.ChName = strings.Split(newname[len(newname)-1], ".")[0]
		// self.Transcoder = "media_excel"
		newname := strings.Split(oldname, "_")
		var code string
		if len(newname) == 4 {
			code = newname[len(newname)-2]
		}

		if len(newname) == 2 {
			code = newname[0]
		}

		chname := strings.Split(newname[len(newname)-1], ".")[0]

		nameByFilename := code + "_" + chname

		if code != "" && code != "VEBOXRMHROHWVMANFUGTBE" {
			self.FileName = "mediaexcel_" + nameByFilename + ".jpg"
			self.ChName = code + " " + chname
		} else {
			self.FileName = "mediaexcel_" + chname + ".jpg"
			self.ChName = chname
		}

		self.Transcoder = transcoder

		// for _, v := range mediaexcelconf {
		// 	if v.UID == chname {
		// 		self.FileName = "mediaexcel_" + strings.ReplaceAll(v.Name, " ", "_") + ".jpg"
		// 		self.ChName = v.Name
		// 		self.Transcoder = v.Service
		// 	}
		// }
	} else if strings.HasPrefix(transcoder, "titan") {
		s := strings.Split(oldname, "-")
		pid := strings.Join(s[:len(s)-1], "-")

		titan := getNameChTitan(pid, transcoder)

		self.FileName = transcoder + "_" + strings.ReplaceAll(titan.Name, " ", "_") + ".jpg"
		self.ChName = titan.Name
		self.Transcoder = "titan"

	} else if strings.HasPrefix(transcoder, "elemental") {
		tmp := strings.Split(oldname, ".")
		tmp = strings.Split(tmp[0], "_")
		id := tmp[len(tmp)-1]

		elemental := getNameChElemental(id, transcoder)

		if elemental.Name == "" {
			self.FileName = oldname
			self.ChName = strings.Split(oldname, ".")[0]
			self.Transcoder = "other"
		} else {
			self.FileName = transcoder + "_" + elemental.Name + ".jpg"
			self.ChName = elemental.Name
			self.Transcoder = "elemental"
		}
	} else {
		self.FileName = oldname
		self.ChName = strings.Split(oldname, ".")[0]
		self.Transcoder = "other"
	}
}

func (self *Message) SetPath(path string) {
	self.Path = path
}

func getDominantColor(img image.Image) color.RGBA {
	var r, g, b, count float64

	rect := img.Bounds()
	for i := 0; i < rect.Max.Y; i++ {
		for j := 0; j < rect.Max.X; j++ {
			c := color.RGBAModel.Convert(img.At(j, i))
			r += float64(c.(color.RGBA).R)
			g += float64(c.(color.RGBA).G)
			b += float64(c.(color.RGBA).B)
			count++
		}
	}

	return color.RGBA{
		R: uint8(r / count),
		G: uint8(g / count),
		B: uint8(b / count),
	}
}

func (self *Message) ConvertToImage(encoded string) error {
	var err error
	var thumbnail image.Image

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(encoded))

	thumbnail, _, err = image.Decode(reader)

	if err != nil {
		return err
	}

	if _, err := os.Stat(self.Path); os.IsNotExist(err) {
		err := os.Mkdir(self.Path, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	out, err := os.Create(self.Path + self.FileName)

	if err != nil {
		return err
	}

	defer out.Close()

	m := resize.Resize(120, 91, thumbnail, resize.Lanczos3)

	self.Colors = getDominantColor(m)

	self.DominantColorRef = cmp.Equal(self.Colors, colorRef1) || cmp.Equal(self.Colors, colorRef2) || cmp.Equal(self.Colors, colorRef3) || cmp.Equal(self.Colors, colorRef4) || cmp.Equal(self.Colors, colorRef5)

	err = jpeg.Encode(out, m, nil)

	if err != nil {
		return err
	}

	return nil
}

func (self *Message) ConvertToImageOriginalSize(encoded string) error {
	var err error
	var thumbnail image.Image

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(encoded))

	thumbnail, _, err = image.Decode(reader)

	if err != nil {
		return err
	}

	if _, err := os.Stat(self.Path); os.IsNotExist(err) {
		err := os.Mkdir(self.Path, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	out, err := os.Create(self.Path + self.FileName)

	if err != nil {
		return err
	}

	defer out.Close()

	//m := resize.Resize(120, 91, thumbnail, resize.Lanczos3)

	//self.Colors = getDominantColor(m)

	//self.DominantColorRef = cmp.Equal(self.Colors, colorRef1) || cmp.Equal(self.Colors, colorRef2) || cmp.Equal(self.Colors, colorRef3) || cmp.Equal(self.Colors, colorRef4)

	err = jpeg.Encode(out, thumbnail, nil)

	if err != nil {
		return err
	}

	return nil
}

func (self *Message) ConvertToImageResize(encoded string, height, width uint) error {
	var err error
	var thumbnail image.Image

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(encoded))

	thumbnail, _, err = image.Decode(reader)

	if err != nil {
		return err
	}

	if _, err := os.Stat(self.Path); os.IsNotExist(err) {
		err := os.Mkdir(self.Path, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	out, err := os.Create(self.Path + self.FileName)

	if err != nil {
		return err
	}

	defer out.Close()

	m := resize.Resize(height, width, thumbnail, resize.Lanczos3)

	//self.Colors = getDominantColor(m)

	//self.DominantColorRef = cmp.Equal(self.Colors, colorRef1) || cmp.Equal(self.Colors, colorRef2) || cmp.Equal(self.Colors, colorRef3) || cmp.Equal(self.Colors, colorRef4)

	err = jpeg.Encode(out, m, nil)

	if err != nil {
		return err
	}

	return nil
}
