package agent

type (
	// AgentOptions are the options used to start an agent.
	AgentOptions struct {
		Port                  string
		ClusterAddress        string
		HostManagementEnabled bool
		SharedSecret          string
		DockerClientTimeout   int
	}

	// ClusterMember is the representation of an agent inside a cluster.
	ClusterMember struct {
		IPAddress string
		Port      string
		NodeName  string
		NodeRole  string
	}

	// AgentMetadata is the representation of the metadata object used to decorate
	// all the objects in an aggregated response.
	AgentMetadata struct {
		Agent struct {
			NodeName string `json:"NodeName"`
		} `json:"Agent"`
	}

	// PciDevice is the representation of a physical pci device on a host
	PciDevice struct {
		Vendor string
		Name   string
	}

	// PhysicalDisk is the representation of a physical disk on a host
	PhysicalDisk struct {
		Vendor string
		Size   uint64
	}

	// HostInfo is the representation of the collection of host information
	HostInfo struct {
		PCIDevices    []PciDevice
		PhysicalDisks []PhysicalDisk
	}

	// ClusterService is used to manage a cluster of agents.
	ClusterService interface {
		Create(advertiseAddr, joinAddr string, tags map[string]string) error
		Members() []ClusterMember
		Leave()
		GetMemberByRole(role string) *ClusterMember
		GetMemberByNodeName(nodeName string) *ClusterMember
	}

	// DigitalSignatureService is used to validate digital signatures.
	DigitalSignatureService interface {
		VerifySignature(signature, key string) (bool, error)
	}

	// InfoService is used to retrieve information from a Docker environment.
	InfoService interface {
		GetInformationFromDockerEngine() (map[string]string, error)
		GetContainerIpFromDockerEngine(containerName string) (string, error)
	}

	// TLSService is used to create TLS certificates to use enable HTTPS.
	TLSService interface {
		GenerateCertsForHost(host string) error
	}

	// SystemService is used to get info about the host
	SystemService interface {
		GetDiskInfo() ([]PhysicalDisk, error)
		GetPciDevices() ([]PciDevice, error)
	}
)

const (
	// AgentVersion represents the version of the agent.
	AgentVersion = "1.2.1"
	// APIVersion represents the version of the agent's API.
	APIVersion = "2"
	// DefaultListenAddr is the default address used by the web server.
	DefaultListenAddr = "0.0.0.0"
	// DefaultAgentPort is the default port exposed by the web server.
	DefaultAgentPort = "9001"
	// DefaultLogLevel is the default logging level.
	DefaultLogLevel = "INFO"
	// DefaultDockerClientTimeout is the default timeout in seconds for Docker client requests.
	DefaultDockerClientTimeout = 10
	// SupportedDockerAPIVersion is the minimum Docker API version supported by the agent.
	SupportedDockerAPIVersion = "1.24"
	// HTTPTargetHeaderName is the name of the header used to specify a target node.
	HTTPTargetHeaderName = "X-PortainerAgent-Target"
	// HTTPManagerOperationHeaderName is the name of the header used to specify that
	// a request must target a manager node.
	HTTPManagerOperationHeaderName = "X-PortainerAgent-ManagerOperation"
	// HTTPSignatureHeaderName is the name of the header containing the digital signature
	// of a Portainer instance.
	HTTPSignatureHeaderName = "X-PortainerAgent-Signature"
	// HTTPPublicKeyHeaderName is the name of the header containing the public key
	// of a Portainer instance.
	HTTPPublicKeyHeaderName = "X-PortainerAgent-PublicKey"
	// HTTPResponseAgentHeaderName is the name of the header that is automatically added
	// to each agent response.
	HTTPResponseAgentHeaderName = "Portainer-Agent"
	// HTTPResponseAgentApiVersion is the name of the header that will have the
	// Portainer Agent API Version.
	HTTPResponseAgentApiVersion = "Portainer-Agent-API-Version"
	// PortainerAgentSignatureMessage is the unhashed content that is signed by the Portainer instance.
	// It is used by the agent during the signature verification process.
	PortainerAgentSignatureMessage = "Portainer-App"
	// ResponseMetadataKey is the JSON field used to store any Portainer related information in
	// response objects.
	ResponseMetadataKey = "Portainer"
	// MemberTagKeyAgentPort is the name of the label storing information about the port exposed
	// by the agent.
	MemberTagKeyAgentPort = "AgentPort"
	// MemberTagKeyNodeName is the name of the label storing information about the name of the
	// node where the agent is running.
	MemberTagKeyNodeName = "NodeName"
	// MemberTagKeyNodeRole is the name of the label storing information about the role of the
	// node where the agent is running.
	MemberTagKeyNodeRole = "NodeRole"
	// ApplicationTagMode is the name of the label storing information about the mode of the application, either
	// "standalone" or "swarm".
	ApplicationTagMode = "Mode"
	// NodeRoleManager represents a manager node.
	NodeRoleManager = "manager"
	// NodeRoleWorker represents a worker node.
	NodeRoleWorker = "worker"
	// TLSCertPath is the default path to the TLS certificate file.
	TLSCertPath = "cert.pem"
	// TLSKeyPath is the default path to the TLS key file.
	TLSKeyPath = "key.pem"
)
