package main

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/inspursoft/wand/src/daemonworker/utils"
)

const (
	kvmToolsPath    = "/root/kvm"
	kvmRegistryPath = "/root/kvmregistry"
)

func prepareKVMHost(nodeIP, nodeSSHPort, username, password, kvmToolkitsPath, kvmRegistrySize, kvmRegistryPort string) error {
	sshPort, _ := strconv.Atoi(nodeSSHPort)
	sshHandler, err := utils.NewSecureShell(nodeIP, sshPort, username, password)
	if err != nil {
		return err
	}
	kvmToolsNodePath := filepath.Join(kvmToolkitsPath, "kvm")
	kvmRegistryNodePath := filepath.Join(kvmToolkitsPath, "kvmregistry")
	err = sshHandler.ExecuteCommand(fmt.Sprintf("mkdir -p %s %s", kvmToolsNodePath, kvmRegistryNodePath))
	if err != nil {
		return err
	}
	err = sshHandler.SecureCopy(kvmToolsPath, kvmToolsNodePath)
	if err != nil {
		return err
	}
	err = sshHandler.SecureCopy(kvmRegistryPath, kvmRegistryNodePath)
	if err != nil {
		return err
	}
	return sshHandler.ExecuteCommand(fmt.Sprintf(`
		cd %s && chmod +x kvmregistry &&nohup ./kvmregistry -size %s -port %s > kvmregistry.out 2>&1 &`,
		kvmRegistryNodePath, kvmRegistrySize, kvmRegistryPort))
}
