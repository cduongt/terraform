resource "occi_virtual_machine" "vm" {
	image_template = "http://occi.carach5.ics.muni.cz/occi/infrastructure/os_tpl#uuid_egi_centos_7_fedcloud_warg_149"
	resource_template = "http://fedcloud.egi.eu/occi/compute/flavour/1.0#small"
	endpoint = "https://carach5.ics.muni.cz:11443"
	name = "occi.core.title=test_vm_small"
	x509 = "/tmp/x509up_u1000"
	public_key = "/home/cduongt/fedcloud.pub"
	count = 3
}

output "virtual_machine_id" {
	value = "${join(",",occi_virtual_machine.vm.*.vm_id)}"
}

output "virtual_machine_ip" {
	value = "${join(",",occi_virtual_machine.vm.*.ip_address)}"
}