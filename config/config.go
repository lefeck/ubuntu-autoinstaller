package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Autoinstall Autoinstall `yaml:"autoinstall" json:"autoinstall"`
}

type Autoinstall struct {
	Apt           AptConfig              `yaml:"apt" json:"apt"`
	Drivers       DriversConfig          `yaml:"drivers" json:"drivers"`
	Identity      Identity               `yaml:"identity" json:"identity"`
	Kernel        KernelConfig           `yaml:"kernel" json:"kernel"`
	Keyboard      KeyboardConfig         `yaml:"keyboard" json:"keyboard"`
	Locale        string                 `yaml:"locale" json:"locale"`
	Network       NetworkConfig          `yaml:"network" json:"network"`
	SSH           SSHConfig              `yaml:"ssh" json:"ssh"`
	Storage       Storage                `yaml:"storage" json:"storage"`
	Updates       string                 `yaml:"updates" json:"updates"`
	Shutdown      string                 `yaml:"shutdown" json:"shutdown"`
	Version       int                    `yaml:"version" json:"version"`
	Packages      []string               `yaml:"packages" json:"packages"`
	EarlyCommands []string               `yaml:"early-commands" json:"early-commands"`
	LateCommands  []string               `yaml:"late-commands" json:"late-commands"`
	UserData      map[string]interface{} `yaml:"user-data" json:"user-data" `
	TimeZone      string                 `yaml:"timezone" json:"timezone"`
}

type AptConfig struct {
	DisableComponents   []string       `yaml:"disable_components" json:"disable_components"`
	DisableSuites       []string       `yaml:"disable_suites" json:"disable_suites"`
	GeoIP               bool           `yaml:"geoip" json:"geoip"`
	PreserveSourcesList bool           `yaml:"preserve_sources_list" json:"preserve_sources_list"`
	Primary             []PrimaryEntry `yaml:"primary" json:"primary"`
}

type PrimaryEntry struct {
	Arches []string `yaml:"arches" json:"arches"`
	URI    string   `yaml:"uri" json:"uri"`
}

type DriversConfig struct {
	Install bool `yaml:"install" json:"install"`
}

type KernelConfig struct {
	Package string `yaml:"package" json:"package"`
}

type KeyboardConfig struct {
	Layout  string `yaml:"layout" json:"layout"`
	Toggle  string `yaml:"toggle" json:"toggle"`
	Variant string `yaml:"variant" json:"variant"`
}

type Identity struct {
	Hostname string `yaml:"hostname" json:"hostname"`
	Password string `yaml:"password" json:"password"`
	Realname string `yaml:"realname" json:"realname"`
	Username string `yaml:"username" json:"username"`
}

type NetworkConfig struct {
	Renderer     string                 `yaml:"renderer" json:"renderer"` // network-renderer (e.g., "networkd", "NetworkManager")
	Version      int                    `yaml:"version" json:"version"`
	Ethernets    map[string]Ethernet    `yaml:"ethernets" json:"ethernets"`
	Wifis        map[string]Wifi        `yaml:"wifis,omitempty" json:"wifis,omitempty"`
	Bridges      map[string]Bridge      `yaml:"bridges,omitempty" json:"bridges,omitempty"`
	Bonds        map[string]Bond        `yaml:"bonds,omitempty" json:"bonds,omitempty"`
	Vlans        map[string]VLAN        `yaml:"vlans,omitempty" json:"vlans,omitempty"`
	Tunnels      map[string]Tunnel      `yaml:"tunnels,omitempty" json:"tunnels,omitempty"`
	Vrfs         map[string]VRF         `yaml:"vrfs,omitempty" json:"vrfs,omitempty"`
	DummyDevices map[string]DummyDevice `yaml:"dummy-devices,omitempty" json:"dummy-devices,omitempty"`
	Modems       map[string]Modem       `yaml:"modems,omitempty" json:"modems,omitempty"`
	NMDevices    map[string]NMDevice    `yaml:"nm-devices,omitempty" json:"nm-devices,omitempty"`
}

// Common groups fields frequently shared across device types. It is embedded
// (inline) into specific device structs so its keys appear at the same level.
type Common struct {
	Match         *Match          `yaml:"match,omitempty" json:"match,omitempty"`
	Dhcp4         bool            `yaml:"dhcp4,omitempty" json:"dhcp4,omitempty"`
	Dhcp6         bool            `yaml:"dhcp6,omitempty" json:"dhcp6,omitempty"`
	Addresses     []string        `yaml:"addresses,omitempty" json:"addresses,omitempty"` // e.g., 192.168.1.10/24
	Nameservers   *Nameservers    `yaml:"nameservers,omitempty" json:"nameservers,omitempty"`
	Routes        []Route         `yaml:"routes,omitempty" json:"routes,omitempty"`                 // policy routes
	RoutingPolicy []RoutingPolicy `yaml:"routing-policy,omitempty" json:"routing-policy,omitempty"` // policy routing rules
	MTU           int             `yaml:"mtu,omitempty" json:"mtu,omitempty"`
	Optional      bool            `yaml:"optional,omitempty" json:"optional,omitempty"`
	LinkLocal     []string        `yaml:"link-local,omitempty" json:"link-local,omitempty"` // e.g., ["ipv4", "ipv6", "off"]
	WakeOnLAN     bool            `yaml:"wakeonlan,omitempty" json:"wakeonlan,omitempty"`
	MACAddress    string          `yaml:"macaddress,omitempty" json:"macaddress,omitempty"`
	Critical      bool            `yaml:"critical,omitempty" json:"critical,omitempty"`
}

// Match selects devices to apply a logical definition.
type Match struct {
	Name       string `yaml:"name,omitempty" json:"name,omitempty"` // glob supported, e.g., "en*"
	MACAddress string `yaml:"macaddress,omitempty" json:"macaddress,omitempty"`
	Driver     string `yaml:"driver,omitempty" json:"driver,omitempty"`
}

// Nameservers config.
type Nameservers struct {
	Addresses []string          `yaml:"addresses,omitempty" json:"addresses,omitempty"`
	Search    []string          `yaml:"search,omitempty" json:"search,omitempty"`
	Options   map[string]string `yaml:"options,omitempty" json:"options,omitempty"`
}

// Route represents a static route.
type Route struct {
	To     string `yaml:"to" json:"to"`                         // destination CIDR, or "default"
	Via    string `yaml:"via,omitempty" json:"via,omitempty"`   // next-hop
	From   string `yaml:"from,omitempty" json:"from,omitempty"` // source prefix
	Metric int    `yaml:"metric,omitempty" json:"metric,omitempty"`
	Table  int    `yaml:"table,omitempty" json:"table,omitempty"`
	OnLink bool   `yaml:"on-link,omitempty" json:"on-link,omitempty"`
	MTU    int    `yaml:"mtu,omitempty" json:"mtu,omitempty"`
	Type   string `yaml:"type,omitempty" json:"type,omitempty"` // e.g., "unicast" (rarely needed)
}

// RoutingPolicy for advanced policy routing (RFC 791 / 4191 style rules).
type RoutingPolicy struct {
	From     string `yaml:"from,omitempty" json:"from,omitempty"`
	To       string `yaml:"to,omitempty" json:"to,omitempty"`
	Table    int    `yaml:"table,omitempty" json:"table,omitempty"`
	Priority int    `yaml:"priority,omitempty" json:"priority,omitempty"`
	FwMark   int    `yaml:"fwmark,omitempty" json:"fwmark,omitempty"`
}

// Ethernet device.
type Ethernet struct {
	Common `yaml:",inline" json:",inline"`
}

// WifiAP describes a single SSID configuration.
type WifiAP struct {
	Password   string `yaml:"password,omitempty" json:"password,omitempty"`
	Mode       string `yaml:"mode,omitempty" json:"mode,omitempty"` // e.g., "infrastructure"
	Hidden     bool   `yaml:"hidden,omitempty" json:"hidden,omitempty"`
	BSSID      string `yaml:"bssid,omitempty" json:"bssid,omitempty"` // AP MAC, optional
	Auth       string `yaml:"auth,omitempty" json:"auth,omitempty"`   // backend-specific, rarely used
	Band       string `yaml:"band,omitempty" json:"band,omitempty"`   // e.g., "2.4GHz", "5GHz"
	Channel    int    `yaml:"channel,omitempty" json:"channel,omitempty"`
	ClientAuth string `yaml:"client-auth,omitempty" json:"client-auth,omitempty"` // tls/wpa-enterprise etc. (simplified)
}

// Wifi device. Uses `access-points` for SSIDs.
type Wifi struct {
	Common       `yaml:",inline" json:",inline"`
	AccessPoints map[string]WifiAP `yaml:"access-points,omitempty" json:"access-points,omitempty"`
}

// Bridge device.
type Bridge struct {
	Common     `yaml:",inline" json:",inline"`
	Interfaces []string      `yaml:"interfaces" json:"interfaces"`
	Parameters *BridgeParams `yaml:"parameters,omitempty" json:"parameters,omitempty"`
}

// Bridge device.
// BridgeParams contains STP and other timing options.
type BridgeParams struct {
	STP          bool `yaml:"stp,omitempty" json:"stp,omitempty"`
	ForwardDelay int  `yaml:"forward-delay,omitempty" json:"forward-delay,omitempty"`
	AgeingTime   int  `yaml:"ageing-time,omitempty" json:"ageing-time,omitempty"`
	Priority     int  `yaml:"priority,omitempty" json:"priority,omitempty"`
}

// Bond (link aggregation) device.
type Bond struct {
	Common     `yaml:",inline" json:",inline"`
	Interfaces []string    `yaml:"interfaces" json:"interfaces"`
	Parameters *BondParams `yaml:"parameters,omitempty" json:"parameters,omitempty"`
}

// BondParams covers common bonding options (active-backup, 802.3ad, etc.).
type BondParams struct {
	Mode               string `yaml:"mode,omitempty" json:"mode,omitempty"` // active-backup, balance-rr, 802.3ad...
	Primary            string `yaml:"primary,omitempty" json:"primary,omitempty"`
	MIIMonitorInterval *int   `yaml:"mii-monitor-interval,omitempty" json:"mii-monitor-interval,omitempty"`
	LACPRate           string `yaml:"lacp-rate,omitempty" json:"lacp-rate,omitempty"` // fast/slow
	TransmitHashPolicy string `yaml:"transmit-hash-policy,omitempty" json:"transmit-hash-policy,omitempty"`
	MinLinks           *int   `yaml:"min-links,omitempty" json:"min-links,omitempty"`
	GratuitousARP      *int   `yaml:"gratuitous-arp,omitempty" json:"gratuitous-arp,omitempty"`
}

// VLAN sub-interface.
type VLAN struct {
	Common `yaml:",inline" json:",inline"`
	ID     int    `yaml:"id" json:"id"`
	Link   string `yaml:"link" json:"link"` // parent interface
}

// Tunnel interface (e.g., GRE, IPIP, SIT, WireGuard* passthrough via nm-devices).
type Tunnel struct {
	Common `yaml:",inline" json:",inline"`
	Mode   string `yaml:"mode" json:"mode"` // gre, ipip, sit, gretap, etc.
	Local  string `yaml:"local,omitempty" json:"local,omitempty"`
	Remote string `yaml:"remote,omitempty" json:"remote,omitempty"`
	Key    *int   `yaml:"key,omitempty" json:"key,omitempty"` // GRE key, optional
	TTL    *int   `yaml:"ttl,omitempty" json:"ttl,omitempty"`
	Port   *int   `yaml:"port,omitempty" json:"port,omitempty"` // if applicable
}

// VRF (Virtual Routing and Forwarding) domain.
type VRF struct {
	Table      int      `yaml:"table" json:"table"`
	Interfaces []string `yaml:"interfaces,omitempty" json:"interfaces,omitempty"`
	Routes     []Route  `yaml:"routes,omitempty" json:"routes,omitempty"`
	// Policy routing specific to this VRF
	RoutingPolicy []RoutingPolicy `yaml:"routing-policy,omitempty" json:"routing-policy,omitempty"`
}

// DummyDevice (dummy interface).
type DummyDevice struct {
	Common `yaml:",inline" json:",inline"`
}

// Modem (cellular) — typically requires NetworkManager renderer.
type Modem struct {
	Common `yaml:",inline" json:",inline"`
	APN    string `yaml:"apn,omitempty" json:"apn,omitempty"`
	PIN    string `yaml:"pin,omitempty" json:"pin,omitempty"`
}

// NMDevice allows passthrough of raw NetworkManager connection keys that netplan
// doesn't model directly.
type NMDevice struct {
	Name        string                 `yaml:"name,omitempty" json:"name,omitempty"`
	UUID        string                 `yaml:"uuid,omitempty" json:"uuid,omitempty"`
	Device      string                 `yaml:"device,omitempty" json:"device,omitempty"`
	Passthrough map[string]interface{} `yaml:"passthrough,omitempty" json:"passthrough,omitempty"`
}

type SSHConfig struct {
	AllowPW        bool     `yaml:"allow-pw" json:"allow-pw"`
	AuthorizedKeys []string `yaml:"authorized-keys" json:"authorized-keys"`
	InstallServer  bool     `yaml:"install-server" json:"install-server"`
}

type Storage struct {
	Config []StorageConfig `yaml:"config" json:"config"`
	Swap   SwapConfig      `yaml:"swap" json:"swap"`
	Grub   GrubConfig      `yaml:"grub" json:"grub"`
}

type StorageConfig struct {
	Type       string     `yaml:"type" json:"type"`
	ID         string     `yaml:"id" json:"id"`
	Device     string     `yaml:"device,omitempty" json:"device,omitempty"`
	Number     int        `yaml:"number,omitempty" json:"number,omitempty"`
	Size       int64      `yaml:"size,omitempty" json:"size,omitempty"`
	Flag       string     `yaml:"flag,omitempty" json:"flag,omitempty"`
	GrubDevice bool       `yaml:"grub_device,omitempty" json:"grub_device,omitempty"`
	Ptable     string     `yaml:"ptable,omitempty" json:"ptable,omitempty"`
	Match      *DiskMatch `yaml:"match,omitempty" json:"match,omitempty"`
	Name       string     `yaml:"name,omitempty" json:"name,omitempty"`
	Preserve   bool       `yaml:"preserve" json:"preserve"`
	Volume     string     `yaml:"volume,omitempty" json:"volume,omitempty"`
	Fstype     string     `yaml:"fstype,omitempty" json:"fstype,omitempty"` // e.g. "ext4", "btrfs", "xfs"
	Volgroup   string     `yaml:"volgroup,omitempty" json:"volgroup,omitempty"`
	Path       string     `yaml:"path,omitempty" json:"path,omitempty"`
	Devices    []string   `yaml:"devices,omitempty" json:"devices,omitempty"`
	Key        string     `yaml:"key,omitempty" json:"key,omitempty"`         // For dm_crypt encryption
	Dm_name    string     `yaml:"dm_name,omitempty" json:"dm_name,omitempty"` // For dm_crypt encryption
	KeyFile    string     `yaml:"keyfile,omitempty" json:"keyfile,omitempty"` // For dm_crypt encryption
	Wipe       string     `yaml:"wipe,omitempty" json:"wipe,omitempty"`       // Supported: superblock, superblock-recursive, pvremove, zero, random
}

type DiskMatch struct {
	Size            string `yaml:"size,omitempty" json:"size,omitempty"`     // e.g. "largest", "smallest", "100G"
	Name            string `yaml:"name,omitempty" json:"name,omitempty"`     //  e.g. "sda"
	Model           string `yaml:"model,omitempty" json:"model,omitempty"`   // Disk model string (e.g., Samsung SSD 860)
	Serial          string `yaml:"serial,omitempty" json:"serial,omitempty"` // serial number
	Path            string `yaml:"path,omitempty" json:"path,omitempty"`     // /dev/disk/by-path/...
	WWN             string `yaml:"wwn,omitempty" json:"wwn,omitempty"`       // World Wide Name
	FirmwareVersion string `yaml:"firmware_version,omitempty" json:"firmware_version,omitempty"`
}

type SwapConfig struct {
	Swap int `yaml:"swap" json:"swap"`
}

type GrubConfig struct {
	ReorderUEFI bool `yaml:"reorder_uefi" json:"reorder_uefi"`
}

// NewDefaultConfig returns a default Config suitable as a starting template.
func NewDefaultConfig() *Config {
	return &Config{
		Autoinstall: Autoinstall{
			Version: 1,
			Apt: AptConfig{
				DisableComponents:   []string{},
				DisableSuites:       []string{"security"},
				GeoIP:               true,
				PreserveSourcesList: false,
				Primary: []PrimaryEntry{
					{
						Arches: []string{"amd64", "i386"},
						URI:    "http://archive.ubuntu.com/ubuntu",
					},
					{
						Arches: []string{"default"},
						URI:    "http://ports.ubuntu.com/ubuntu-ports",
					},
				},
			},
			Drivers: DriversConfig{
				Install: false,
			},
			Identity: Identity{
				Hostname: "ubuntu-server",
				Password: "$6$FnHPy.4giHfDcI8r$vllY94AceZFt20M6WzxXKfOLuycO7o8QD8RJJwVs/DSGcHNNF4ilbGrukXPMA3irAdiimHo0fl2zUdkbRO9vh1",
				Realname: "ubuntu",
				Username: "ubuntu",
			},
			Kernel: KernelConfig{
				Package: "linux-generic",
			},
			Keyboard: KeyboardConfig{
				Layout:  "us",
				Toggle:  "",
				Variant: "",
			},
			Locale: "en_US.UTF-8",
			Network: NetworkConfig{
				Version:  2,
				Renderer: "networkd",
				Ethernets: map[string]Ethernet{
					"ens160": Ethernet{
						Common: Common{
							Dhcp4:     false,
							Dhcp6:     false,
							Addresses: []string{"192.168.6.100/24"},
							Routes: []Route{
								{
									To:  "default",
									Via: "192.168.6.1",
								},
							},
							Nameservers: &Nameservers{
								Addresses: []string{"114.114.114.114"},
								Search:    []string{},
							},
							Match: &Match{
								Name: "ens160",
							},
						},
					},
				},
			},
			SSH: SSHConfig{
				AllowPW:        true,
				AuthorizedKeys: []string{},
				InstallServer:  true,
			},
			Storage: Storage{
				Config: []StorageConfig{
					{
						Type:       "disk",
						ID:         "disk0",
						GrubDevice: true,
						Ptable:     "gpt",
						Wipe:       "superblock-recursive",
						Match: &DiskMatch{
							Size: "largest",
						},
						Name:     "",
						Preserve: false,
					},
					{
						Type:     "partition",
						ID:       "bios-grub-part",
						Device:   "disk0",
						Number:   1,
						Size:     1048576,
						Flag:     "bios_grub",
						Preserve: false,
						Wipe:     "superblock",
					},
					{
						Type:     "partition",
						ID:       "boot-part",
						Device:   "disk0",
						Number:   2,
						Size:     2147483648,
						Preserve: false,
						Wipe:     "superblock",
					},
					{
						Type:     "format",
						ID:       "boot-fs",
						Volume:   "boot-part",
						Fstype:   "ext4",
						Preserve: false,
					},
					{
						Type:     "partition",
						ID:       "pv-part",
						Device:   "disk0",
						Number:   3,
						Size:     -1,
						Preserve: false,
						Wipe:     "superblock",
					},
					{
						Type:     "lvm_volgroup",
						ID:       "vg0",
						Name:     "ubuntu-vg",
						Devices:  []string{"pv-part"},
						Preserve: false,
					},
					{
						Type:     "lvm_partition",
						ID:       "lv-swap",
						Volgroup: "vg0",
						Name:     "swap",
						Size:     1073741824,
						Preserve: false,
						Wipe:     "superblock",
					},
					{
						Type:     "format",
						ID:       "fs-swap",
						Volume:   "lv-swap",
						Fstype:   "swap",
						Preserve: false,
					},
					{
						Type:     "lvm_partition",
						ID:       "lv-root",
						Volgroup: "vg0",
						Name:     "ubuntu-lv",
						Size:     -1,
						Preserve: false,
						Wipe:     "superblock",
					},
					{
						Type:     "format",
						ID:       "fs-root",
						Volume:   "lv-root",
						Fstype:   "ext4",
						Preserve: false,
					},
					{
						Type:     "dm_crypt",
						ID:       "dm_crypt-0",
						Volume:   "lv-root",
						Key:      "ubuntu",
						Preserve: false,
					},
					{
						Type:     "mount",
						ID:       "mount-root",
						Path:     "/",
						Device:   "fs-root",
						Preserve: false,
					},
					{
						Type:     "mount",
						ID:       "mount-boot",
						Path:     "/boot",
						Device:   "boot-fs",
						Preserve: false,
					},
					{
						Type:     "mount",
						ID:       "mount-swap",
						Path:     "",
						Device:   "fs-swap",
						Preserve: false,
					},
				},
				Swap: SwapConfig{
					Swap: 0,
				},
				Grub: GrubConfig{
					ReorderUEFI: false,
				},
			},
			Updates:  "security",
			Shutdown: "reboot",
			Packages: []string{},
			LateCommands: []string{
				"cp -rp /cdrom/mnt /target/",
				"chmod +x /target/mnt/script/install-pkgs.sh",
				"curtin in-target --target=/target -- /mnt/script/install-pkgs.sh",
				"chmod +x /target/mnt/script/config.sh",
				"curtin in-target --target=/target -- /mnt/script/config.sh",
				"cp /cdrom/rc-local.service /target/lib/systemd/system/rc-local.service",
				"curtin in-target --target=/target -- ln -s /lib/systemd/system/rc-local.service /etc/systemd/system/rc-local.service",
				"cp -p /cdrom/rc.local /target/etc/rc.local",
				"chmod +x /target/etc/rc.local",
				"systemctl daemon-reload",
				"cp -rp /cdrom/mnt/wsa  /target/opt",
				"cp /cdrom/mnt/wsa/wsa.service /target/lib/systemd/system/wsa.service",
				"curtin in-target --target=/target -- ln -sn /lib/systemd/system/wsa.service  /etc/systemd/system/multi-user.target.wants/wsa.service",
			},
		},
	}
}

// LoadConfig reads a YAML config from a file path.
func LoadConfig(path string) (*Config, error) {
	// Read file content
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Unmarshal YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %v", err)
	}

	return &config, nil
}

// SaveConfig marshals and writes a config to a YAML file.
func SaveConfig(config *Config, path string) error {
	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// Write file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// Validate performs basic validation for top-level config.
func (c *Config) Validate() error {
	if c.Autoinstall.Version == 0 {
		return fmt.Errorf("version must be non-zero")
	}

	if err := c.Autoinstall.Identity.Validate(); err != nil {
		return fmt.Errorf("invalid identity config: %v", err)
	}

	if err := c.Autoinstall.Network.Validate(); err != nil {
		return fmt.Errorf("invalid network config: %v", err)
	}

	if err := c.Autoinstall.Storage.Validate(); err != nil {
		return fmt.Errorf("invalid storage config: %v", err)
	}

	return nil
}

// Validate performs basic checks for identity section.
func (i *Identity) Validate() error {
	if i.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if i.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if i.Hostname == "" {
		return fmt.Errorf("hostname cannot be empty")
	}
	return nil
}

// Validate performs basic checks for network section.
func (n *NetworkConfig) Validate() error {
	if n.Version != 2 {
		return fmt.Errorf("network config version must be 2")
	}
	if len(n.Ethernets) == 0 {
		return fmt.Errorf("at least one ethernet interface is required")
	}
	return nil
}

// Validate performs basic checks for storage section.
func (s *Storage) Validate() error {
	if len(s.Config) == 0 {
		return fmt.Errorf("at least one storage config is required")
	}
	return nil
}
