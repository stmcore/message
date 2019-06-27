package message

import (
	"encoding/base64"
	"encoding/json"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/kpango/glg"
	"github.com/nfnt/resize"
)

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

func init() {

	file, _ := ioutil.ReadFile("./titan.json")
	err := json.Unmarshal([]byte(file), &titanconf)

	if err != nil {
		log.Fatalln(err)
	}

	file, _ = ioutil.ReadFile("./elemental.json")
	err = json.Unmarshal([]byte(file), &elementalconf)

	if err != nil {
		log.Fatalln(err)
	}

	file, _ = ioutil.ReadFile("./mediaexcel.json")
	err = json.Unmarshal([]byte(file), &mediaexcelconf)

	if err != nil {
		log.Fatalln(err)
	}
}

func (self *Message) SetNameDotJPG(oldname string, transcoder string) {

	self.FileName = oldname

	if transcoder == "media_excel" {
		// newname := strings.Split(oldname, "_")
		// self.FileName = "mediaexcel_" + newname[len(newname)-1]
		// self.ChName = strings.Split(newname[len(newname)-1], ".")[0]
		// self.Transcoder = "media_excel"
		newname := strings.Split(oldname, "_")
		chname := strings.Split(newname[len(newname)-1], ".")[0]

		for _, v := range mediaexcelconf {
			if v.UID == chname {
				self.FileName = "mediaexcel_" + strings.ReplaceAll(v.Name, " ", "_") + ".jpg"
				self.ChName = v.Name
				self.Transcoder = v.Service
			}
		}
	} else if transcoder == "titan16" {
		for _, v := range titanconf {
			if v.Service == transcoder && strings.HasPrefix(oldname, v.UID) {
				self.FileName = transcoder + "_" + strings.ReplaceAll(v.Name, " ", "_") + ".jpg"
				self.ChName = v.Name
				self.Transcoder = "titan"
			}
		}
	} else if transcoder == "titan17" {
		for _, v := range titanconf {
			if v.Service == "titan17" && strings.HasPrefix(oldname, v.UID) {
				self.FileName = transcoder + "_" + strings.ReplaceAll(v.Name, " ", "_") + ".jpg"
				self.ChName = v.Name
				self.Transcoder = "titan"
			}
		}
	} else if transcoder == "titan25" {
		for _, v := range titanconf {
			if v.Service == "titan25" && strings.HasPrefix(oldname, v.UID) {
				self.FileName = transcoder + "_" + strings.ReplaceAll(v.Name, " ", "_") + ".jpg"
				self.ChName = v.Name
				self.Transcoder = "titan"
			}
		}
	} else if transcoder == "elemental12" {
		tmp := strings.Split(oldname, ".")
		tmp = strings.Split(tmp[0], "_")
		id := tmp[len(tmp)-1]
		for _, v := range elementalconf {
			if v.UID == id && v.Service == "elemental12" {
				self.FileName = transcoder + "_" + strings.ReplaceAll(v.Name, " ", "_") + ".jpg"
				self.ChName = v.Name
				self.Transcoder = "elemental"
			}
		}
	} else if transcoder == "elemental13" {
		tmp := strings.Split(oldname, ".")
		tmp = strings.Split(tmp[0], "_")
		id := tmp[len(tmp)-1]
		for _, v := range elementalconf {
			if v.UID == id && v.Service == "elemental13" {
				self.FileName = transcoder + "_" + strings.ReplaceAll(v.Name, " ", "_") + ".jpg"
				self.ChName = v.Name
				self.Transcoder = "elemental"
			}
		}
	} else if transcoder == "elemental14" {
		tmp := strings.Split(oldname, ".")
		tmp = strings.Split(tmp[0], "_")
		id := tmp[len(tmp)-1]
		for _, v := range elementalconf {
			if v.UID == id && v.Service == "elemental14" {
				self.FileName = transcoder + "_" + strings.ReplaceAll(v.Name, " ", "_") + ".jpg"
				self.ChName = v.Name
				self.Transcoder = "elemental"
			}
		}
	} else if transcoder == "elemental15" {
		tmp := strings.Split(oldname, ".")
		tmp = strings.Split(tmp[0], "_")
		id := tmp[len(tmp)-1]
		for _, v := range elementalconf {
			if v.UID == id && v.Service == "elemental15" {
				self.FileName = transcoder + "_" + strings.ReplaceAll(v.Name, " ", "_") + ".jpg"
				self.ChName = v.Name
				self.Transcoder = "elemental"
			}
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

	out, err := os.Create(self.Path + self.FileName)

	if err != nil {
		return err
	}

	defer out.Close()

	m := resize.Resize(120, 91, thumbnail, resize.Lanczos3)

	self.Colors = getDominantColor(m)

	self.DominantColorRef = cmp.Equal(self.Colors, colorRef1) || cmp.Equal(self.Colors, colorRef2) || cmp.Equal(self.Colors, colorRef3) || cmp.Equal(self.Colors, colorRef4)

	err = jpeg.Encode(out, m, nil)

	if err != nil {
		return err
	}

	return nil
}
