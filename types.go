package osc

import (
	"encoding/json"
	"io"
)

// Version represents an API level supported by the camera
type Version int

// Version constants
const (
	V1 Version = 1 << iota
	V2
)

// Info type holds the information about the given endpoint
type Info struct {
	Manufacturer string    //The camera manufacturer.
	Model        string    //The camera model.
	Serial       string    //Serial number.
	Firmware     string    //Current firmware version.
	Support      string    //URL for the camera’s support webpage.
	GPS          bool      //True if the the camera has GPS.
	Gyro         bool      //True if the camera has Gyroscope.
	Uptime       int       //Number of seconds since the camera boot.
	API          []string  //List of supported APIs.
	Endpoints    Endpoints //Object	A JSON object containing information about the camera’s endpoints. See the next table.
	Versions     []Version
	Vendor       map[string]interface{}
}

// Endpoints holds the available ports for server and updates communications
type Endpoints struct {
	HTTP  Ports
	HTTPS Ports
}

// Ports holds the server and updates port information
type Ports struct {
	Server  int
	Updates int
}

// convert a JSON decoded map from the response to an info type. We need to handle `_vendorSpecific` fields dynamically, which Go isn't so good
// at with JSON, so we use the stream directly to parse it ourselves
func parseInfo(r io.Reader) (*Info, error) {
	body := map[string]interface{}{}
	err := json.NewDecoder(r).Decode(&body)

	if err != nil {
		return nil, err
	}

	info := Info{
		Endpoints: Endpoints{
			HTTP: Ports{}, HTTPS: Ports{},
		},
		Vendor: map[string]interface{}{},
	}

	for k, v := range body {
		switch k {
		case "manufacturer":
			info.Manufacturer = v.(string)
		case "model":
			info.Model = v.(string)
		case "serialNumber":
			info.Serial = v.(string)
		case "firmwareVersion":
			info.Firmware = v.(string)
		case "supportUrl":
			info.Support = v.(string)
		case "gps":
			info.GPS = v.(bool)
		case "gyro":
			info.Gyro = v.(bool)
		case "uptime":
			info.Uptime = int(v.(float64))
		case "api":
			values := v.([]interface{})
			info.API = make([]string, len(values))
			for idx, v := range values {
				info.API[idx] = v.(string)
			}
		case "apiLevel":
			info.Versions = v.([]Version)
		case "endpoints":
			eps := v.(map[string]interface{})

			if value, ok := eps["httpPort"]; ok {
				info.Endpoints.HTTP.Server = int(value.(float64))
			}
			if value, ok := eps["httpUpdatesPort"]; ok {
				info.Endpoints.HTTP.Updates = int(value.(float64))
			}
			if value, ok := eps["httpsPort"]; ok {
				info.Endpoints.HTTPS.Server = int(value.(float64))
			}
			if value, ok := eps["httpsUpdatesPort"]; ok {
				info.Endpoints.HTTPS.Updates = int(value.(float64))
			}
		}

		if k[0] != '_' {
			continue
		}

		info.Vendor[k[1:]] = v

	}

	if info.Versions == nil {
		info.Versions = []Version{V1}
	}

	return &info, nil

}
