---
title: おうちkubernetes上でMySQLを起動しようとしてもうまくいかなかった
date: 2026-11-23T21:38:00.077527+09:00
description: ""
categories: tech
tagas:
    - proxmox
    - 勉強メモ
publish: true
---

Claudeくんに教えてもらったものをメモ。

## おうちkubernetes上でMySQLを起動しようとしてもうまくいかなかった

おうちkubernetes上にMySQLを立てようとしたら以下のエラーが発生

```
❯ kubectl logs mysql-0 -n mysql-test           
Fatal glibc error: CPU does not support x86-64-v2
```

???

glibcとは、Linuxで広く使われているCの標準ライブラリ。ファイルやメモリなどOSの機能を使うための関数をまとめたもので、Pythonでいうとosモジュールに近い。

x86-64-v2とは、命令セットのレベル。x86-64-v1, v2, v3などがある。同じアーキテクチャの中で、どこまでの命令を持つかを表しているもの。v1よりもv2が、v2よりもv3が多くの命令を持っている。

現在の状況は、glibcから、今のCPUがx86-64-v2をサポートしてないよ　というエラーが出た。という状況。

もう少し噛み砕くと、x86-64-v2に含まれる命令のうち、このCPUが持っていないものがある　という意味だそう。

つまり、kubernetesのclusterの上で実際にMySQLのdocker imageを実行しようとした時に、x86-64-v2の命令が理解できる前提で作られているイメージなのに、実際に動かそうとしている場所がx86-64-v2を理解できないCPUであると言ってるのである。

なぜそんなことが起きるのだろうか。

### 原因

今回の原因は、今回動かしているkubernetes clusterはProxmox上のVMで動かしているものであることと、ProxmoxがこのVMに対してCPUをkvm64という種類として見せていることが原因だった。

kvm64とは、Proxmoxが仮想的にVMに割り当てるCPUの種類で、命令セットでいうところのx86-64-v1相当の命令セットしか持っていない。

そのため、MySQLのimageがx86-64-v2を要求しているのに、それを用意できないためにエラーが発生してしまったのである。

※ ちなみに、MySQLのimageはv2を要求しているが、MySQLはv2を必須としてない。今回のエラーはMySQLのベースイメージであるOracle Linux 9がv2を要求しているために発生している。ubuntuやOracle Linux8や7はv1で動作するため、base imageをubuntuにした上でMySQLをinstallしてbuildしてあげれば動く。

例：oracle linux9のpodを立てた時→v2がないため失敗する

```
❯ kubectl run ol9-test --rm -it --restart=Never --image=oraclelinux:9-slim -- cat /etc/os-release

Fatal glibc error: CPU does not support x86-64-v2
If you don't see a command prompt, try pressing enter.
pod "ol9-test" deleted from default namespace
pod default/ol9-test terminated (Error)
```

例：oracle linux8のpodを立てた時→成功
```
kubernetes-manifests on  main via ❄️  impure (nix-shell-env) took 11s 
❯ kubectl run ol8-test --rm -it --restart=Never --image=oraclelinux:8-slim -- cat /etc/os-release

NAME="Oracle Linux Server"
VERSION="8.10"
ID="ol"
ID_LIKE="fedora"
VARIANT="Server"
VARIANT_ID="server"
VERSION_ID="8.10"
PLATFORM_ID="platform:el8"
PRETTY_NAME="Oracle Linux Server 8.10"
ANSI_COLOR="0;31"
CPE_NAME="cpe:/o:oracle:linux:8:10:server"
HOME_URL="https://linux.oracle.com/"
BUG_REPORT_URL="https://github.com/oracle/oracle-linux"

ORACLE_BUGZILLA_PRODUCT="Oracle Linux 8"
ORACLE_BUGZILLA_PRODUCT_VERSION=8.10
ORACLE_SUPPORT_PRODUCT="Oracle Linux"
ORACLE_SUPPORT_PRODUCT_VERSION=8.10
```

### 解決策

この問題は、Proxmoxで動いてるVMのCPUのモデルをkvm64からx86-64-v2-AESに変えることで解決する。

現在使ってるProxmoxのバージョンが古いため、デフォルトで作られるVMのCPUのモデルがkvm64であるが、Proxmox VE 8.0からはx86-64-v2-AESがデフォルトで作られる。ちなみに、実際にProxmoxが作動しているCPU以上のレベルの命令セットを指定することはできないので注意。そこまでエミュレートしてくれるわけじゃないんだね。

### 感想

普段だったらClaudeの言いなりになって、肉MCPとなりProxmoxやterminalをいじるところだったが、ちゃんと立ち止まってエラーをふかぼってみた。

従来だと、一つひとつの概念の前提が理解できずに挫折していたけど、自分の理解が間違っているとClaudeくんが教えてくれるから正しいところに導かれる気がする。

今一度、CPUについての理解が深まったので調べてよかった。