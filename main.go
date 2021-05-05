// SPDX-FileCopyrightText: Copyright 2021 The go-tmux Authors
// SPDX-License-Identifier: BSD-3-Clause

// Command kube-tmux prints Kubernetes context and namespace to tmux status line.
package main

import (
	"fmt"
	"os"
	"text/template"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	defaultSepalater       = "/"
	defaultContextFormat   = "{{.Context}}"
	defaultNamespaceFormat = "{{.Namespace}}"
	defaultformat          = defaultContextFormat + defaultSepalater + defaultNamespaceFormat
)

type kubeContext struct {
	Context   string
	Namespace string
}

func main() {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.RawConfig()
	if err != nil {
		fmt.Fprintln(os.Stdout, "[ERROR] could not get kubeconfig")
		return
	}
	if len(config.Contexts) == 0 {
		fmt.Fprintln(os.Stdout, "[ERROR] kubeconfig is empty")
		return
	}

	curCtx := config.CurrentContext
	if curCtx == "" {
		kctx := kubeContext{
			Context: "empty",
		}
		if err := printContext(kctx, defaultContextFormat); err != nil {
			fmt.Fprintln(os.Stdout, "[ERROR] could not print kube context")
		}
		return
	}

	kctx := kubeContext{
		Context: curCtx,
	}

	curNs := config.Contexts[curCtx].Namespace
	if curNs == "" {
		curNs = corev1.NamespaceDefault
	}
	kctx.Namespace = curNs

	format := defaultformat
	if len(os.Args) > 1 {
		format = os.Args[1]
	}

	if err := printContext(kctx, format); err != nil {
		fmt.Fprintf(os.Stdout, "[ERROR] could not print kube context: %v\n", err)
	}
}

func printContext(kctx kubeContext, format string) error {
	return template.Must(template.New("kube-tmux").Parse(format)).Execute(os.Stdout, kctx)
}
