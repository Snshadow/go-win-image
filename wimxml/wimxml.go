// Package wimxml contains xml structure used in .wim file.
package wimxml

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"strconv"
	"unicode/utf16"
)

type Wim struct {
	XMLName    xml.Name `xml:"WIM"`
	TotalBytes uint64   `xml:"TOTALBYTES"`
	Images     []Image  `xml:"IMAGE"`
}

type Image struct {
	Index                uint32   `xml:"INDEX,attr"`
	DirCount             uint32   `xml:"DIRCOUNT"`
	FileCount            uint32   `xml:"FILECOUNT"`
	TotalBytes           uint64   `xml:"TOTALBYTES"`
	HardlinkBytes        uint64   `xml:"HARDLINKBYTES"`
	CreationTime         Filetime `xml:"CREATIONTIME"`
	LastModificationTime Filetime `xml:"LASTMODIFICATIONTIME"`
	WimBoot              uint32   `xml:"WIMBOOT"`
	Windows              Windows  `xml:"WINDOWS"`
	Name                 string   `xml:"NAME"`
	Description          string   `xml:"DESCRIPTION"`
	Flags                string   `xml:"FLAGS"`
	DisplayName          string   `xml:"DISPLAYNAME"`
	DisplayDescription   string   `xml:"DISPLAYDESCRIPTION"`
}

type Filetime struct {
	HighPart uint32 `xml:"HIGHPART"`
	LowPart  uint32 `xml:"LOWPART"`
}

type _filetime struct {
	HighPart string `xml:"HIGHPART"`
	LowPart  string `xml:"LOWPART"`
}

func (f Filetime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	ff := _filetime{
		HighPart: "0x" + strconv.FormatUint(uint64(f.HighPart), 16),
		LowPart:  "0x" + strconv.FormatUint(uint64(f.LowPart), 16),
	}

	return e.EncodeElement(ff, start)
}

func (f *Filetime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var ff _filetime
	var err error
	if err = d.DecodeElement(&ff, &start); err != nil {
		return err
	}

	tempHigh, err := strconv.ParseUint(ff.HighPart, 0, 32)
	if err != nil {
		return err
	}
	f.HighPart = uint32(tempHigh)
	tempLow, err := strconv.ParseUint(ff.LowPart, 0, 32)
	if err != nil {
		return err
	}
	f.LowPart = uint32(tempLow)

	return nil
}

type Windows struct {
	Arch             uint32        `xml:"ARCH"`
	ProductName      string        `xml:"PRODUCTNAME"`
	EditionID        string        `xml:"EDITIONID"`
	InstallationType string        `xml:"INSTALLATIONTYPE"`
	ServicingData    ServicingData `xml:"SERVICINGDATA"`
	ProductType      string        `xml:"PRODUCTTYPE"`
	ProductSuite     string        `xml:"PRODUCTSUITE"`
	Languages        Languages     `xml:"LANGUAGES"`
	Version          Version       `xml:"VERSION"`
	SystemRoot       string        `xml:"SYSTEMROOT"`
}

type ServicingData struct {
	GdrDurRevision    string `xml:"GDRDUREVISION"`
	PKeyConfigVersion string `xml:"PKEYCONFIGVERSION"`
	ImageState        string `xml:"IMAGESTATE"`
}

type Languages struct {
	Language []string `xml:"LANGUAGE"`
	Fallback Fallback `xml:"FALLBACK"`
	Default  string   `xml:"DEFAULT"`
}

type Fallback struct {
	Language string `xml:"LANGUAGE,attr"`
	Value    string `xml:",chardata"`
}

type Version struct {
	Major   uint32 `xml:"MAJOR"`
	Minor   uint32 `xml:"MINOR"`
	Build   uint32 `xml:"BUILD"`
	SpBuild uint32 `xml:"SPBUILD"`
	SpLevel uint32 `xml:"SPLEVEL"`
	Branch  string `xml:"BRANCH"`
}

func BytesToWimXml(buf []byte) (*Wim, error) {
	if len(buf) < 2 {
		return nil, errors.New("buffer size is too small")
	}

	// skip BOM(0xff, 0xfe)
	if buf[0] == 0xff && buf[1] == 0xfe {
		buf = buf[2:]
	}

	xmlReader := bytes.NewReader(buf)
	xmlData := make([]uint16, len(buf)/2)
	binary.Read(xmlReader, binary.LittleEndian, xmlData)
	decoded := utf16.Decode(xmlData)

	wimInfo := &Wim{}

	err := xml.Unmarshal([]byte(string(decoded)), wimInfo)
	if err != nil {
		return nil, err
	}

	return wimInfo, nil
}

func WimXmlToBytes(info *Wim) ([]byte, error) {
	infoXml, err := xml.Marshal(info)
	if err != nil {
		return nil, err
	}

	encoded := utf16.Encode([]rune(string(infoXml)))

	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, encoded)

	// add BOM(UTF-16 LE) at the front
	infoXml = append([]byte{0xff, 0xfe}, buf.Bytes()...)

	return infoXml, nil
}
