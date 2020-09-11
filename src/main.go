package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Alarm is the desired format
// that all alerts should be converted to
type Alarm struct {
	FaultName   string `json:"FaultName"`
	FaultSource string `json:"FaultSource"`
	Severity    string `json:"Severity,omitempty"`
	Description string `json:"Description,omitempty"`
	Reason      string `json:"Reason,omitempty"`
}

type RequestParams struct {
	ServiceType string            `json:"ServiceType"`
	Alert       map[string]string `json:"Alert"`
}

type ResponseData struct {
	Status  string `json:"Status"`
	Message string `json:"Message,omitempty"`
	Alarm   Alarm  `json:"Alarm,omitempty"`
}

// To cache the mappers we got before
// to avoid reading the same file again
var mapperCache = map[string]map[string]string{}

func convertAlert(reqBody []byte) ResponseData {

	var responseData ResponseData

	fmt.Println("Getting Request Params")
	var reqParams RequestParams

	errReq := json.Unmarshal(reqBody, &reqParams)
	if errReq != nil {
		fmt.Println("Error in Getting Request Params: ", errReq)
		responseData.Status = "Fail"
		responseData.Message = "Error in Getting Request Params"
		return responseData
	}

	// The actual content of the alert
	alertData := reqParams.Alert

	fmt.Println("Getting Alert Mapper")
	// The map that maps Alarm fields to Alert fields
	var alertMapper map[string]string

	serviceType := reqParams.ServiceType
	mapper, mapperExist := mapperCache[serviceType]

	if mapperExist {
		alertMapper = mapper
		fmt.Println("The mapper was found in cache")
	} else {
		fmt.Println("Read the mapper from the file")
		mapperFileName := "../mapping/" + strings.ToLower(serviceType) + ".json"
		mapperFileContent, _ := ioutil.ReadFile(mapperFileName)

		errMap := json.Unmarshal(mapperFileContent, &alertMapper)
		if errMap != nil {
			fmt.Println("Error in Getting Alert Mapper: ", errMap)
			responseData.Status = "Fail"
			responseData.Message = "Error in Getting Alert Mapper"
			return responseData
		}
		mapperCache[serviceType] = alertMapper
	}

	fmt.Println("################### The Mapper ################")
	fmt.Println(alertMapper)

	fmt.Println("Creating The Alarm")
	var alarm Alarm

	fn, hasFN := alertMapper["FaultName"]
	if hasFN {
		fnVal := alertData[fn]
		alarm.FaultName = fnVal
	}

	fs, hasFS := alertMapper["FaultSource"]
	if hasFS {
		fsVal := alertData[fs]
		alarm.FaultSource = fsVal
	}

	sev, hasSev := alertMapper["Severity"]
	if hasSev {
		sevVal := alertData[sev]
		alarm.Severity = sevVal
	}

	desc, hasDesc := alertMapper["Description"]
	if hasDesc {
		descVal := alertData[desc]
		alarm.Description = descVal
	}

	res, hasRes := alertMapper["Reason"]
	if hasRes {
		resVal := alertData[res]
		alarm.Reason = resVal
	}

	fmt.Println("################### The Alarm ################")
	fmt.Printf("%+v\n", alarm)

	responseData.Status = "Success"
	responseData.Alarm = alarm

	return responseData
}

func conversionAPI(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	result := convertAlert(reqBody)
	json.NewEncoder(w).Encode(result)
}

func handleRequests() {
	http.HandleFunc("/convert", conversionAPI)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	handleRequests()
}
