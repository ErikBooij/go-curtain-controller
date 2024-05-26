package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"curtain-controller/app/config"
	"curtain-controller/app/drivers"
)

func RunAPI(apiConfig config.APIConfig, devices drivers.DeviceList) {
	mux := http.NewServeMux()

	mux.HandleFunc("/shutter/aqara", func(writer http.ResponseWriter, request *http.Request) {
		deviceId, position, err := getParameters(request)

		if err != nil {
			respondWithError(writer, http.StatusBadRequest, err)
			return
		}

		device, ok := devices.AqaraShutters[deviceId]

		if !ok {
			respondWithError(writer, http.StatusNotFound, "Device '%s' is unknown", deviceId)
			return
		}

		if !device.SetPosition(position) {
			respondWithError(writer, http.StatusInternalServerError, "Unable to set the position of the shutter at this time")
			return
		}

		writer.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/curtain/slide", func(writer http.ResponseWriter, request *http.Request) {
		deviceId, position, err := getParameters(request)

		if err != nil {
			respondWithError(writer, http.StatusBadRequest, err)
			return
		}

		device, ok := devices.SlideCurtains[deviceId]

		if !ok {
			respondWithError(writer, http.StatusNotFound, "Device '%s' is unknown", deviceId)
			return
		}

		if !device.SetPosition(position) {
			respondWithError(writer, http.StatusInternalServerError, "Unable to set the position of the curtain at this time")
			return
		}

		writer.WriteHeader(http.StatusOK)
	})

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", apiConfig.Bind, apiConfig.Port), mux); err != nil {
		panic(err)
	}
}

func getParameters(request *http.Request) (string, float64, error) {
	deviceId := request.URL.Query().Get("deviceId")
	position := request.URL.Query().Get("position")

	if deviceId == "" {
		return "", 0, errors.New("No deviceId parameter present in request")
	}

	if position == "" {
		return "", 0, errors.New("No position parameter present in request")
	}

	pos, err := strconv.ParseFloat(position, 64)

	if err != nil {
		return "", 0, fmt.Errorf("Unparseable position value (got '%s')", position)
	}

	return deviceId, pos, nil
}

func respondWithError(writer http.ResponseWriter, statusCode int, message interface{}, args ...interface{}) {
	if m, ok := message.(error); ok {
		message = m.Error()
	}

	writer.WriteHeader(statusCode)
	writer.Write([]byte(fmt.Sprintf(message.(string), args...)))
}
