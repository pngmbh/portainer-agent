package tunnel

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/portainer/agent"
)

type pollStatusResponse struct {
	Status          string           `json:"status"`
	Port            int              `json:"port"`
	Schedules       []agent.Schedule `json:"schedules"`
	CheckinInterval float64          `json:"checkin"`
}

func (operator *Operator) poll() error {
	pollURL := fmt.Sprintf("%s/api/endpoints/%s/status", operator.key.PortainerInstanceURL, operator.key.EndpointID)
	req, err := http.NewRequest("GET", pollURL, nil)
	if err != nil {
		return err
	}

	// TODO: @@DOCUMENTATION: document the extra security layer added by the Edge ID
	req.Header.Set(agent.HTTPEdgeIdentifierHeaderName, operator.edgeID)

	resp, err := operator.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[DEBUG] [http,edge,poll] [response_code: %d] [message: Poll request failure]", resp.StatusCode)
		return errors.New("short poll request failed")
	}

	var responseData pollStatusResponse
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] [http,edge,poll] [status: %s] [port: %d] [schedule_count: %d] [checkin_interval_seconds: %f]", responseData.Status, responseData.Port, len(responseData.Schedules), responseData.CheckinInterval)

	if responseData.Status == "IDLE" && operator.tunnelClient.IsTunnelOpen() {
		log.Printf("[DEBUG] [http,edge,poll] [status: %s] [message: Idle status detected, shutting down tunnel]", responseData.Status)

		err := operator.tunnelClient.CloseTunnel()
		if err != nil {
			log.Printf("[ERROR] [http,edge,poll] [message: Unable to shutdown tunnel] [error: %s]", err)
		}
	}

	if responseData.Status == "REQUIRED" && !operator.tunnelClient.IsTunnelOpen() {
		tunnelConfig := agent.TunnelConfig{
			ServerAddr:       operator.key.TunnelServerAddr,
			ServerFingerpint: operator.key.TunnelServerFingerprint,
			Credentials:      operator.key.Credentials,
			RemotePort:       strconv.Itoa(responseData.Port),
			LocalAddr:        operator.apiServerAddr,
		}

		log.Printf("[DEBUG] [http,edge,poll] [status: %s] [port: %d] [message: Required status detected, creating reverse tunnel]", responseData.Status, responseData.Port)

		err = operator.tunnelClient.CreateTunnel(tunnelConfig)
		if err != nil {
			return err
		}

		operator.ResetActivityTimer()
	}

	err = operator.scheduleManager.Schedule(responseData.Schedules)
	if err != nil {
		log.Printf("[ERROR] [http,edge,cron] [message: an error occured during schedule management] [err: %s]", err)
	}

	if responseData.CheckinInterval != operator.pollIntervalInSeconds {
		log.Printf("[DEBUG] [http,edge,poll] [old_interval: %f] [new_interval: %f] [message: updating checkin interval]", operator.pollIntervalInSeconds, responseData.CheckinInterval)
		operator.pollIntervalInSeconds = responseData.CheckinInterval
		go operator.restartStatusPollLoop()
	}

	return nil
}