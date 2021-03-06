package vsphere

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// defaultAPITimeout is a default timeout value that is passed to functions
// requiring contexts, and other various waiters.
var defaultAPITimeout = time.Minute * 5

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_USER", nil),
				Description: "The user name for vSphere API operations.",
			},

			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_PASSWORD", nil),
				Description: "The user password for vSphere API operations.",
			},

			"vsphere_server": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_SERVER", nil),
				Description: "The vSphere Server name for vSphere API operations.",
			},
			"allow_unverified_ssl": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_ALLOW_UNVERIFIED_SSL", false),
				Description: "If set, VMware vSphere client will permit unverifiable SSL certificates.",
			},
			"vcenter_server": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_VCENTER", nil),
				Deprecated:  "This field has been renamed to vsphere_server.",
			},
			"client_debug": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_CLIENT_DEBUG", false),
				Description: "govomomi debug",
			},
			"client_debug_path_run": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_CLIENT_DEBUG_PATH_RUN", ""),
				Description: "govomomi debug path for a single run",
			},
			"client_debug_path": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_CLIENT_DEBUG_PATH", ""),
				Description: "govomomi debug path for debug",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"vsphere_datacenter":               resourceVSphereDatacenter(),
			"vsphere_file":                     resourceVSphereFile(),
			"vsphere_folder":                   resourceVSphereFolder(),
			"vsphere_host_port_group":          resourceVSphereHostPortGroup(),
			"vsphere_host_virtual_switch":      resourceVSphereHostVirtualSwitch(),
			"vsphere_license":                  resourceVSphereLicense(),
			"vsphere_virtual_disk":             resourceVSphereVirtualDisk(),
			"vsphere_virtual_machine":          resourceVSphereVirtualMachine(),
			"vsphere_nas_datastore":            resourceVSphereNasDatastore(),
			"vsphere_vmfs_datastore":           resourceVSphereVmfsDatastore(),
			"vsphere_virtual_machine_snapshot": resourceVSphereVirtualMachineSnapshot(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"vsphere_datacenter": dataSourceVSphereDatacenter(),
			"vsphere_host":       dataSourceVSphereHost(),
			"vsphere_vmfs_disks": dataSourceVSphereVmfsDisks(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	// Handle backcompat support for vcenter_server; once that is removed,
	// vsphere_server can just become a Required field that is referenced inline
	// in Config below.
	server := d.Get("vsphere_server").(string)

	if server == "" {
		server = d.Get("vcenter_server").(string)
	}

	if server == "" {
		return nil, fmt.Errorf(
			"One of vsphere_server or [deprecated] vcenter_server must be provided.")
	}

	config := Config{
		User:          d.Get("user").(string),
		Password:      d.Get("password").(string),
		InsecureFlag:  d.Get("allow_unverified_ssl").(bool),
		VSphereServer: server,
		Debug:         d.Get("client_debug").(bool),
		DebugPathRun:  d.Get("client_debug_path_run").(string),
		DebugPath:     d.Get("client_debug_path").(string),
	}

	return config.Client()
}
