// SPDX-FileCopyrightText: Copyright 2021 The go-tmux Authors
// SPDX-License-Identifier: BSD-3-Clause

// Command kube-tmux prints Kubernetes context and namespace to tmux status line.
package main

import (
	"log"
	"os"
	"path/filepath"
	"text/template"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	defaultSepalater = string(filepath.Separator)
	defaultformat    = "{{.Context}}" + defaultSepalater + "{{.Namespace}}"
)

type currentContext struct {
	Context   string
	Namespace string
}

func main() {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.RawConfig()
	if err != nil {
		os.Exit(1)
	}
	if len(config.Contexts) == 0 {
		log.Fatal("empty context in kubeconfig")
	}

	curCtx := config.CurrentContext
	curNs := config.Contexts[curCtx].Namespace
	if curNs == "" {
		curNs = corev1.NamespaceDefault
	}

	format := defaultformat
	if len(os.Args) > 1 {
		format = os.Args[1]
	}

	cctx := currentContext{
		Context:   curCtx,
		Namespace: curNs,
	}
	template.Must(template.New("tmux").Parse(format)).Execute(os.Stdout, cctx)
}
