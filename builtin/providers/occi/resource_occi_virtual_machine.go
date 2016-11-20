package occi

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"os/exec"
	"strconv"
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
			"storage_size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  0,
			},
			"vm_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_id": &schema.Schema{
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

	// arguments for VM creation
	endpoint := d.Get("endpoint").(string)
	image_template := d.Get("image_template").(string)
	resource_template := d.Get("resource_template").(string)
	proxy_file := d.Get("x509").(string)
	vm_name := d.Get("name").(string)
	context_file := d.Get("context").(string)

	// create VM
	cmd_name := "occi"
	cmd_args_create := []string{"-e", endpoint, "-n", "x509", "-x", proxy_file, "-X", "-a", "create", "-r", "compute", "-M", image_template, "-M", resource_template, "-t", vm_name, "-T", strings.Join([]string{"user_data=file:///", context_file}, "")}

	if cmdOut, err = exec.Command(cmd_name, cmd_args_create...).Output(); err != nil {
		return fmt.Errorf("Error while creating virtual machine: %s", err.Error())
	}
	compute_id_address := strings.Split(string(cmdOut), "/")
	compute_id := strings.Trim(compute_id_address[len(compute_id_address)-1], "\n")
	compute := strings.Join([]string{"/compute/", compute_id}, "")
	d.Set("vm_id", compute)
	d.SetId(compute)

	// get IP address
	cmd_args_describe := []string{"-e", endpoint, "-n", "x509", "-x", proxy_file, "-X", "-a", "describe", "-r", compute}

	if cmdOut, err = exec.Command(cmd_name, cmd_args_describe...).Output(); err != nil {
		return fmt.Errorf("Error while trying to get IP address: %s", err.Error())
	}

	byte_array := bytes.Fields(cmdOut)
	for i, line := range byte_array {
		if bytes.Contains(line, []byte("occi.networkinterface.address")) {
			ip_line := string(byte_array[i+2][:])
			d.Set("ip_address", ip_line)
			break
		}
	}

	// if storage variable is set, create storage
	storage_size := d.Get("storage_size").(int)
	if storage_size > 0 {
		random, _ := rand.Int(rand.Reader, big.NewInt(32)) // random hash for name, as storage name must be unique
		hash := md5.New()
		hash.Write([]byte(random.String()))
		random_hash := hex.EncodeToString(hash.Sum(nil))
		storage_params := strings.Join([]string{"occi.storage.size=", "'num(", strconv.Itoa(storage_size), ")',occi.core.title=storage_terraform", "_", random_hash}, "")
		cmd_args_storage := []string{"-e", endpoint, "-n", "x509", "-x", proxy_file, "-X", "-a", "create", "-r", "storage", "-t", storage_params}
		if cmdOut, err = exec.Command(cmd_name, cmd_args_storage...).Output(); err != nil {
			return fmt.Errorf("Error while creating storage: %s", err.Error())
		}
		storage_id_split := strings.Split(string(cmdOut), "/")
		storage_id_trim := strings.Trim(storage_id_split[len(storage_id_split)-1], "\n")
		storage_id := strings.Join([]string{"/storage/", storage_id_trim}, "")
		d.Set("storage_id", storage_id)

		// link storage to VM
		cmd_args_storage_link := []string{"-e", endpoint, "-n", "x509", "-x", proxy_file, "-X", "-a", "link", "-r", compute, "-j", storage_id}
		if cmdOut, err = exec.Command(cmd_name, cmd_args_storage_link...).Output(); err != nil {
			return fmt.Errorf("Error while linking storage %s to VM %s: %s", compute, storage_id, err.Error())
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

	// destroy VM
	cmd_name := "occi"
	cmd_args := []string{"-e", endpoint, "-n", "x509", "-x", proxy_file, "-X", "-a", "delete", "-r", vm_id}
	if _, err = exec.Command(cmd_name, cmd_args...).Output(); err != nil {
		return fmt.Errorf("Error while destroying VM %s: %s", vm_id, err.Error())
	}

	// if storage has been provisioned, destroy it too
	storage_id := d.Get("storage_id").(string)
	if storage_id != "" {
		cmd_args_storage := []string{"-e", endpoint, "-n", "x509", "-x", proxy_file, "-X", "-a", "delete", "-r", storage_id}
		if _, err = exec.Command(cmd_name, cmd_args_storage...).Output(); err != nil {
			return fmt.Errorf("Error while destroying storage %s: %s", storage_id, err.Error())
		}
	}
	return nil
}
