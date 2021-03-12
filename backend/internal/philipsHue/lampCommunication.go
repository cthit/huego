package philipsHue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/viddem/huego/internal/utilities"
	"io/ioutil"
	"net/http"
	"strconv"
)

func SetLampCall(values *utilities.LampData, config *utilities.HueConfig, lampNumber uint16) error {
	url := fmt.Sprintf("%s/lights/%d/state", config.BaseUrl, lampNumber)
	return updateLights(config, values, url)
}

func SetAllLampsCall(values *utilities.LampData, config *utilities.HueConfig) error {
	url := fmt.Sprintf("%s/groups/1/action", config.BaseUrl)
	return updateLights(config, values, url)
}

func updateLights(config *utilities.HueConfig, values *utilities.LampData, url string) error {
	client := &http.Client{}
	jsonData, err := json.Marshal(values)

	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

type LampsDataJson struct {
	Lights map[string]LampDataJson `json:"lights"`
	Groups map[string]interface{} `json:"groups"`
	Config map[string]interface{} `json:"config"`
	Schedules map[string]interface{} `json:"schedules"`
	Scenes map[string]interface{} `json:"scenes"`
	Rules map[string]interface{} `json:"rules"`
	Sensors map[string]interface{} `json:"sensors"`
	ResourceLinks map[string]interface{} `json:"resourceLinks"`
}

type LampDataJson struct {
	State utilities.LampData `json:"state"`
}

type LampWithCoordinates struct {
	Light utilities.Light `json:"light"`
	State utilities.SimpleLampData `json:"state"`
}

func GetLightsInfo(config *utilities.HueConfig) ([]LampWithCoordinates, error) {
	res, err := http.Get(config.BaseUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var lampsData LampsDataJson
	err = json.Unmarshal(body, &lampsData)
	if err != nil {
		return nil, err
	}

	var lampDatas []LampWithCoordinates
	for key, val := range lampsData.Lights {
		id, err := strconv.Atoi(key)
		if err != nil {
			return nil, err
		}

		light, err := config.GetLightFromMap(uint16(id))
		if err == nil {
			lampDatas = append(lampDatas, LampWithCoordinates{
				Light: light,
				State: val.State.Simplify(),
			})
		}
	}

	return lampDatas, nil
}
