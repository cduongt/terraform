package occi

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceVirtualMachineCreate,
		Read:   resourceVirtualMachineRead,
		Delete: resourceVirtualMachineDelete,

		Schema: map[string]*schema.Schema{
			"endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"x509": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"image_template": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_template": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"context": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vm_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		cmdOut []byte
		err    error
	)

	endpoint := d.Get("endpoint").(string)
	image_template := d.Get("image_template").(string)
	resource_template := d.Get("resource_template").(string)
	proxy_file := d.Get("x509").(string)
	vm_name := d.Get("name").(string)
	context_file := d.Get("context").(string)

	cmd_name := "occi"
	cmd_args_create := []string{"-e", endpoint, "-n", "x509", "-x", proxy_file, "-X", "-a", "create", "-r", "compute", "-M", image_template, "-M", resource_template, "-t", vm_name, "-T", strings.Join([]string{"user_data=file:///", context_file}, "")}

	if cmdOut, err = exec.Command(cmd_name, cmd_args_create...).Output(); err != nil {
		return err
	}

	compute_id_address := strings.Split(string(cmdOut), "/")
	compute_id := strings.Trim(compute_id_address[len(compute_id_address)-1], "\n")
	compute := strings.Join([]string{"/compute/", compute_id}, "")
	d.Set("vm_id", compute)
	d.SetId(compute)
	cmd_args_describe := []string{"-e", endpoint, "-n", "x509", "-x", proxy_file, "-X", "-a", "describe", "-r", compute}

	if cmdOut, err = exec.Command(cmd_name, cmd_args_describe...).Output(); err != nil {
		return err
	}

	byte_array := bytes.Fields(cmdOut)
	for i, line := range byte_array {
		if bytes.Contains(line, []byte("occi.networkinterface.address")) {
			ip_line := string(byte_array[i+2][:])
			d.Set("ip_address", ip_line)
			break
		}
	}

	return nil
}

func resourceVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		_   []byte
		err error
	)
	endpoint := d.Get("endpoint").(string)
	proxy_file := d.Get("x509").(string)
	vm_id := d.Get("vm_id").(string)

	cmd_name := "occi"
	cmd_args := []string{"-e", endpoint, "-n", "x509", "-x", proxy_file, "-X", "-a", "delete", "-r", vm_id}
	if _, err = exec.Command(cmd_name, cmd_args...).Output(); err != nil {
		return err
	}

	return nil
}
