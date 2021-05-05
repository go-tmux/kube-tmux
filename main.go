// SPDX-FileCopyrightText: Copyright 2021 The go-tmux Authors
// SPDX-License-Identifier: BSD-3-Clause

// Command kube-tmux prints Kubernetes context and namespace to tmux status line.
package main

import (
	"flag"
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

var (
	ctxFg     string
	ctxBg     string
	sepFg     string
	sepBg     string
	nsFg      string
	nsBg      string
	separator string
)

func init() {
	flag.StringVar(&ctxFg, "ctxFg", "", "Context foreground colour")
	flag.StringVar(&ctxBg, "ctxBg", "", "Context background colour")
	flag.StringVar(&sepFg, "sepFg", "", "Separator foreground colour")
	flag.StringVar(&sepBg, "sepBg", "", "Separator background colour")
	flag.StringVar(&nsFg, "nsFg", "", "Nasespace foreground colour")
	flag.StringVar(&nsBg, "nsBg", "", "Nasespace background colour")
	flag.StringVar(&separator, "separator", defaultSepalater, "Separator of Context and Nasespace")
}

func main() {
	flag.Parse()

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

	var format string
	switch {
	case flag.CommandLine.NArg() > 1:
		format = flag.CommandLine.Arg(0)
	default:
		// TODO(zchee): refactoring
		if ctxFg != "" || ctxBg != "" {
			format = "#["
		}
		if ctxFg != "" {
			format += "fg=" + ctxFg
		}
		if ctxBg != "" {
			format += ",bg=" + ctxBg
		}
		if ctxFg != "" || ctxBg != "" {
			format += "]"
		}
		format += defaultContextFormat

		if sepFg != "" || sepBg != "" {
			format += "#["
		}
		if sepFg != "" {
			format += "fg=" + sepFg
		}
		if sepBg != "" {
			format += ",bg=" + sepBg
		}
		if sepFg != "" || sepBg != "" {
			format += "]"
		}
		format += separator

		if nsFg != "" || nsBg != "" {
			format += "#["
		}
		if nsFg != "" {
			format += "fg=" + nsFg
		}
		if nsBg != "" {
			format += ",bg=" + nsBg
		}
		if nsFg != "" || nsBg != "" {
			format += "]"
		}
		format += defaultNamespaceFormat
		format += "#[fg=default#,bg=default]"
	}

	if err := printContext(kctx, format); err != nil {
		fmt.Fprintf(os.Stdout, "[ERROR] could not print kube context: %v\n", err)
	}
}

func printContext(kctx kubeContext, format string) error {
	return template.Must(template.New("kube-tmux").Parse(format)).Execute(os.Stdout, kctx)
}
