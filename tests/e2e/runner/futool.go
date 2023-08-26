package runner

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type FutoolRunnerConfig struct {
	FuryConfigTemplate string

	ImageTag   string
	IncludeIBC bool

	EnableAutomatedUpgrade  bool
	FuryUpgradeName         string
	FuryUpgradeHeight       int64
	FuryUpgradeBaseImageTag string

	SkipShutdown bool
}

// FutoolRunner implements a NodeRunner that spins up local chains with futool.
// It has support for the following:
// - running a Fury node
// - optionally, running an IBC node with a channel opened to the Fury node
// - optionally, start the Fury node on one version and upgrade to another
type FutoolRunner struct {
	config FutoolRunnerConfig
}

var _ NodeRunner = &FutoolRunner{}

// NewFutoolRunner creates a new FutoolRunner.
func NewFutoolRunner(config FutoolRunnerConfig) *FutoolRunner {
	return &FutoolRunner{
		config: config,
	}
}

// StartChains implements NodeRunner.
// For FutoolRunner, it sets up, runs, and connects to a local chain via futool.
func (k *FutoolRunner) StartChains() Chains {
	// install futool if not already installed
	installFutoolCmd := exec.Command("./scripts/install-futool.sh")
	installFutoolCmd.Stdout = os.Stdout
	installFutoolCmd.Stderr = os.Stderr
	if err := installFutoolCmd.Run(); err != nil {
		panic(fmt.Sprintf("failed to install futool: %s", err.Error()))
	}

	// start local test network with futool
	log.Println("starting fury node")
	futoolArgs := []string{"testnet", "bootstrap", "--fury.configTemplate", k.config.FuryConfigTemplate}
	// include an ibc chain if desired
	if k.config.IncludeIBC {
		futoolArgs = append(futoolArgs, "--ibc")
	}
	// handle automated upgrade functionality, if defined
	if k.config.EnableAutomatedUpgrade {
		futoolArgs = append(futoolArgs,
			"--upgrade-name", k.config.FuryUpgradeName,
			"--upgrade-height", fmt.Sprint(k.config.FuryUpgradeHeight),
			"--upgrade-base-image-tag", k.config.FuryUpgradeBaseImageTag,
		)
	}
	// start the chain
	startFuryCmd := exec.Command("futool", futoolArgs...)
	startFuryCmd.Env = os.Environ()
	startFuryCmd.Env = append(startFuryCmd.Env, fmt.Sprintf("FURY_TAG=%s", k.config.ImageTag))
	startFuryCmd.Stdout = os.Stdout
	startFuryCmd.Stderr = os.Stderr
	log.Println(startFuryCmd.String())
	if err := startFuryCmd.Run(); err != nil {
		panic(fmt.Sprintf("failed to start fury: %s", err.Error()))
	}

	// wait for chain to be live.
	// if an upgrade is defined, this waits for the upgrade to be completed.
	if err := waitForChainStart(futoolFuryChain); err != nil {
		k.Shutdown()
		panic(err)
	}
	log.Println("fury is started!")

	chains := NewChains()
	chains.Register("fury", &futoolFuryChain)
	if k.config.IncludeIBC {
		chains.Register("ibc", &futoolIbcChain)
	}
	return chains
}

// Shutdown implements NodeRunner.
// For FutoolRunner, it shuts down the local futool network.
// To prevent shutting down the chain (eg. to preserve logs or examine post-test state)
// use the `SkipShutdown` option on the config.
func (k *FutoolRunner) Shutdown() {
	if k.config.SkipShutdown {
		log.Printf("would shut down but SkipShutdown is true")
		return
	}
	log.Println("shutting down fury node")
	shutdownFuryCmd := exec.Command("futool", "testnet", "down")
	shutdownFuryCmd.Stdout = os.Stdout
	shutdownFuryCmd.Stderr = os.Stderr
	if err := shutdownFuryCmd.Run(); err != nil {
		panic(fmt.Sprintf("failed to shutdown futool: %s", err.Error()))
	}
}
