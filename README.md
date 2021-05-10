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
        Namespace background colour
  -nsFg string
        Namespace foreground colour
  -sepBg string
        Separator background colour
  -sepFg string
        Separator foreground colour
  -separator string
        Separator of Context and Namespace (default "/")
```

I’m using the follows:

```console
set -g status-right "#(kube-tmux -separator=':' -ctxFg='colour011#,bright' -ctxBg 'default' -sepFg 'colour015#,bright' -nsFg 'colour075#,bright')"
```

## Benchmark

### Benchmark command:

```console
$ kube.tmux
#[fg=blue]⎈ #[fg=colour]#[fg=]gke_asia-northeast1_cluster#[fg=colour250]:#[fg=]kube-system

$ kube-tmux -separator=':' -ctxFg='blue' -sepFg 'colour250' -nsFg 'default'
#[fg=blue]gke_asia-northeast1_cluster#[fg=colour250]:#[fg=default]kube-system#[fg=default#,bg=default]
```

### Note

- `go version`
  - `go version devel go1.17-543e098320 Wed May 5 19:17:46 2021 +0000 X:staticlockranking darwin/amd64`
- tmux `default-shell` is `/bin/sh`.
- On darwin, `echo 3 > /proc/sys/vm/drop_caches` similar command is `sync && sudo purge`

### Result

```console
$ hyperfine --shell /bin/sh --prepare 'sync && sudo purge' --max-runs 10 --warmup 3 'kube.tmux'
Benchmark #1: kube.tmux
  Time (mean ± σ):     275.2 ms ±  54.0 ms    [User: 111.2 ms, System: 74.3 ms]
  Range (min … max):   224.3 ms … 368.1 ms    10 runs

$ hyperfine --shell /bin/sh --prepare 'sync && sudo purge' --max-runs 10 --warmup 3 "kube-tmux -separator=':' -ctxFg='blue' -sepFg 'colour250' -nsFg 'default'"
Benchmark #1: kube-tmux -separator=':' -ctxFg='blue' -sepFg 'colour250' -nsFg 'default'
  Time (mean ± σ):     102.2 ms ±  64.6 ms    [User: 8.0 ms, System: 35.6 ms]
  Range (min … max):    12.5 ms … 178.0 ms    10 runs
```

| Command                                                                     | Mean [ms]    | Min [ms] | Max [ms] | Relative |
|:----------------------------------------------------------------------------|-------------:|---------:|---------:|---------:|
| `kube.tmux`                                                                 | 290.2 ± 62.9 | 223.6    | 382.2    | 1.00     |
| `kube-tmux -separator=':' -ctxFg='blue' -sepFg 'colour250' -nsFg 'default'` |  84.2 ± 59.6 |  15.9    | 156.5    | 1.00     |
