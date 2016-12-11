resource "occi_virtual_machine" "vm" {
	image_template = "http://occi.carach5.ics.muni.cz/occi/infrastructure/os_tpl#uuid_egi_ubuntu_server_14_04_lts_fedcloud_warg_131"
	resource_template = "http://fedcloud.egi.eu/occi/compute/flavour/1.0#large"
	endpoint = "https://carach5.ics.muni.cz:11443"
	name = "test_vm_small"
	x509 = "/tmp/x509up_u1000"
	init_file = "/home/cduongt/context"
	count = 3
}

output "virtual_machine_id" {
	value = "${join(",",occi_virtual_machine.vm.*.vm_id)}"
}

output "virtual_machine_ip" {
	value = "${join(",",occi_virtual_machine.vm.*.ip_address)}"
}
