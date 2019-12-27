package enrollprovision

import "github.com/mattrax/Mattrax/pkg/xml"

type WapProvisioningDoc struct {
	XMLName        xml.Name            `xml:"wap-provisioningdoc"`
	Version        string              `xml:"version,attr"`
	Characteristic []WapCharacteristic `xml:"characteristic"`
}

type WapCharacteristic struct {
	Type            string `xml:"type,attr,omitempty"`
	Params          []WapParameter
	Characteristics []WapCharacteristic `xml:"characteristic,omitempty"`
}

type WapParameter struct {
	XMLName  xml.Name `xml:"parm"`
	Name     string   `xml:"name,attr,omitempty"`
	Value    string   `xml:"value,attr,omitempty"`
	DataType string   `xml:"datatype,attr,omitempty"`
}
