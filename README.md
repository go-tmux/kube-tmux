# kube-tmux

Command kube-tmux prints Kubernetes context and namespace to tmux status line.

## Usage

```console
$ kube-tmux -h
Usage of kube-tmux:
  -ctxBg string
        Context background colour
  -ctxFg string
        Context foreground colour
  -nsBg string
        Nasespace background colour
  -nsFg string
        Nasespace foreground colour
  -sepBg string
        Separator background colour
  -sepFg string
        Separator foreground colour
  -separator string
        Separator of Context and Nasespace (default "/")
```

```console
set -g status-right "#(kube-tmux -separator=':' -ctxFg='colour011#,bright' -ctxBg 'default' -sepFg 'colour015#,bright' -nsFg 'colour075#,bright')"
```

## Benchmark

```console
$ hyperfine -S /usr/local/bin/bash --max-runs 10 --warmup 3 'kube.tmux'
Benchmark #1: kube.tmux
  Time (mean ± σ):     109.0 ms ±   2.7 ms    [User: 102.7 ms, System: 37.0 ms]
  Range (min … max):   106.8 ms … 116.3 ms    10 runs

$ hyperfine -S /usr/local/bin/bash --max-runs 10 --warmup 3 "kube-tmux -separator=':' -ctxFg='colour011#,bright' -ctxBg 'default' -sepFg 'colour015#,bright' -nsFg 'colour075#,bright'"
Benchmark #1: kube-tmux -separator=':' -ctxFg='colour011#,bright' -ctxBg 'default' -sepFg 'colour015#,bright' -nsFg 'colour075#,bright'
  Time (mean ± σ):      10.8 ms ±   0.6 ms    [User: 6.7 ms, System: 2.7 ms]
  Range (min … max):    10.3 ms …  12.1 ms    10 runs
```
