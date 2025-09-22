package generator

import "text/template"

var (
	NameTemplate = template.Must(template.New("iso").Parse(
		"ubuntu-server-{{.codename}}-autoinstall"))

	XorrisoCmdUbuntu2204Template = template.Must(template.New("xorriso").Parse(
		Xorriso + " -as mkisofs -r " +
			"-V {{.Label}} " +
			"-o {{.Output}} " +
			"--grub2-mbr ../BOOT/1-Boot-NoEmul.img " +
			"-partition_offset 16 " +
			"--mbr-force-bootable " +
			"-append_partition 2 28732ac11ff8d211ba4b00a0c93ec93b ../BOOT/2-Boot-NoEmul.img " +
			"-appended_part_as_gpt " +
			"-iso_mbr_part_type a2a0d0ebe5b9334487c068b6b72699c7 " +
			"-c /boot.catalog " +
			"-b /boot/grub/i386-pc/eltorito.img " +
			"-no-emul-boot -boot-load-size 4 -boot-info-table --grub2-boot-info " +
			"-eltorito-alt-boot " +
			"-e --interval:appended_partition_2::: " +
			"-no-emul-boot ."))

	XorrisoCmdUbuntu2004Template = template.Must(template.New("xorriso").Parse(
		Xorriso + " -as mkisofs -r " +
			"-V {{.Label}} " +
			"-o {{.Output}} " +
			"-J " +
			"-b isolinux/isolinux.bin " +
			"-c isolinux/boot.cat " +
			"-no-emul-boot " +
			"-boot-load-size 4 " +
			"-isohybrid-mbr /usr/lib/ISOLINUX/isohdpfx.bin " +
			"-boot-info-table " +
			"-input-charset utf-8 " +
			"-eltorito-alt-boot " +
			"-e boot/grub/efi.img " +
			"-no-emul-boot " +
			"-isohybrid-gpt-basdat " +
			"."))

	ShellTemplate = `#!/bin/bash
# The default installation package will be downloaded to /cdrom/mnt/packages/ directory
cp /etc/apt/sources.list /etc/apt/sources.list.bak
echo 'deb [trusted=yes] file:///mnt/packages/ ./' > /etc/apt/sources.list
apt-get update
{{range .}}
apt-get install -y {{.}}
{{end}}
`
)
