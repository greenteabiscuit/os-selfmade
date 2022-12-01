export SHELL := /bin/bash -o pipefail

OS = $(shell uname -s)
BUILD_DIR := build
BUILD_ABS_DIR := $(CURDIR)/$(BUILD_DIR)

VBOX_VM_NAME := gopher-os
QEMU ?= qemu-system-x86_64

DIR := day4

img:
	vagrant ssh -c 'cd /home/vagrant/workspace/$(DIR); make img'

run: img
	cd $(DIR); qemu-system-i386 -drive file=haribote.img,format=raw,if=floppy -boot a

run-qemu: GC_FLAGS += -B
run-qemu: iso
	$(QEMU) -cdrom $(iso_target) -vga std -d int,cpu_reset -no-reboot

run-vbox: iso
	VBoxManage createvm --name $(VBOX_VM_NAME) --ostype "Linux_64" --register || true
	VBoxManage storagectl $(VBOX_VM_NAME) --name "IDE Controller" --add ide || true
	VBoxManage storageattach $(VBOX_VM_NAME) --storagectl "IDE Controller" --port 0 --device 0 --type dvddrive \
		--medium $(iso_target) || true
	VBoxManage startvm $(VBOX_VM_NAME)

# When building gdb target disable optimizations (-N) and inlining (l) of Go code
gdb: GC_FLAGS += -N -l
gdb: iso
	$(QEMU) -M accel=tcg -vga std -s -S -cdrom $(iso_target) &
	sleep 1
	gdb \
	    -ex 'add-auto-load-safe-path $(pwd)' \
	    -ex 'set disassembly-flavor intel' \
	    -ex 'layout split' \
	    -ex 'set arch i386:intel' \
	    -ex 'file $(kernel_target)' \
	    -ex 'target remote localhost:1234' \
	    -ex 'set arch i386:x86-64:intel' \
	    -ex 'source $(GOROOT)/src/runtime/runtime-gdb.py' \
	    -ex 'set substitute-path $(VAGRANT_SRC_FOLDER) $(shell pwd)'
	@killall $(QEMU) || true
