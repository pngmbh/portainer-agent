package os

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/portainer/agent"
)

const (
	EnvKeyAgentHost         = "AGENT_HOST"
	EnvKeyAgentPort         = "AGENT_PORT"
	EnvKeyClusterAddr       = "AGENT_CLUSTER_ADDR"
	EnvKeyAgentSecret       = "AGENT_SECRET"
	EnvKeyCapHostManagement = "CAP_HOST_MANAGEMENT"
	EnvKeyEdge              = "EDGE"
	EnvKeyEdgeKey           = "EDGE_KEY"
	EnvKeyEdgeID            = "EDGE_ID"
	EnvKeyEdgeServerHost    = "EDGE_SERVER_HOST"
	EnvKeyEdgeServerPort    = "EDGE_SERVER_PORT"
	EnvKeyEdgePollInterval  = "EDGE_POLL_INTERVAL"
	EnvKeyEdgeSleepInterval = "EDGE_SLEEP_INTERVAL"
	EnvKeyLogLevel          = "LOG_LEVEL"
)

type EnvOptionParser struct{}

func NewEnvOptionParser() *EnvOptionParser {
	return &EnvOptionParser{}
}

func (parser *EnvOptionParser) Options() (*agent.Options, error) {
	options := &agent.Options{
		AgentServerAddr:       agent.DefaultAgentAddr,
		AgentServerPort:       agent.DefaultAgentPort,
		ClusterAddress:        os.Getenv(EnvKeyClusterAddr),
		HostManagementEnabled: false,
		SharedSecret:          os.Getenv(EnvKeyAgentSecret),
		EdgeID:                os.Getenv(EnvKeyEdgeID),
		EdgeServerAddr:        agent.DefaultEdgeServerAddr,
		EdgeServerPort:        agent.DefaultEdgeServerPort,
		EdgePollInterval:      agent.DefaultEdgePollInterval,
		EdgeSleepInterval:     agent.DefaultEdgeSleepInterval,
		LogLevel:              agent.DefaultLogLevel,
	}

	if os.Getenv(EnvKeyCapHostManagement) == "1" {
		options.HostManagementEnabled = true
	}

	if os.Getenv(EnvKeyEdge) == "1" {
		options.EdgeMode = true
	}

	if options.EdgeMode && options.EdgeID == "" {
		return nil, errors.New("missing mandatory " + EnvKeyEdgeID + " environment variable")
	}

	agentAddrEnv := os.Getenv(EnvKeyAgentHost)
	if agentAddrEnv != "" {
		options.AgentServerAddr = agentAddrEnv
	}

	agentPortEnv := os.Getenv(EnvKeyAgentPort)
	if agentPortEnv != "" {
		_, err := strconv.Atoi(agentPortEnv)
		if err != nil {
			return nil, errors.New("invalid port format in " + EnvKeyAgentPort + " environment variable")
		}
		options.AgentServerPort = agentPortEnv
	}

	edgeAddrEnv := os.Getenv(EnvKeyEdgeServerHost)
	if edgeAddrEnv != "" {
		options.EdgeServerAddr = edgeAddrEnv
	}

	edgePortEnv := os.Getenv(EnvKeyEdgeServerPort)
	if edgePortEnv != "" {
		_, err := strconv.Atoi(edgePortEnv)
		if err != nil {
			return nil, errors.New("invalid port format in " + EnvKeyEdgeServerPort + " environment variable")
		}
		options.EdgeServerPort = edgePortEnv
	}

	edgeKeyEnv := os.Getenv(EnvKeyEdgeKey)
	if edgeKeyEnv != "" {
		options.EdgeKey = edgeKeyEnv
	}

	edgePollIntervalEnv := os.Getenv(EnvKeyEdgePollInterval)
	if edgePollIntervalEnv != "" {
		_, err := time.ParseDuration(edgePollIntervalEnv)
		if err != nil {
			return nil, errors.New("invalid time duration format in " + EnvKeyEdgePollInterval + " environment variable")
		}
		options.EdgePollInterval = edgePollIntervalEnv
	}

	edgeSleepIntervalEnv := os.Getenv(EnvKeyEdgeSleepInterval)
	if edgeSleepIntervalEnv != "" {
		_, err := time.ParseDuration(edgeSleepIntervalEnv)
		if err != nil {
			return nil, errors.New("invalid time duration format in " + EnvKeyEdgeSleepInterval + " environment variable")
		}
		options.EdgeSleepInterval = edgeSleepIntervalEnv
	}

	logLevelEnv := os.Getenv(EnvKeyLogLevel)
	if logLevelEnv != "" {
		options.LogLevel = logLevelEnv
	}

	return options, nil
}