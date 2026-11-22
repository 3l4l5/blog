---
title: 11/22勉強したもの(ansible)
date: 2026-11-22T21:58:25.8129+09:00
description: ""
categories: tech
tagas:
    - ansible
publish: false
---

## ansible

ansibleは多くのサーバーに対して同じ設定を行いたいときのためのツール。

ツールを使わずに設定を行う場合は、一つひとつのサーバーに対してsshしながらコマンドを流していかなくてはならないが、ansibleを使うことで、あらかじめ設定しておいたコマンドを、コマンド一つで流すことができる。

シンプルな構成なら、shell scriptで十分だが、diffを見れたりなどの運用面での嬉しい機能がくっついてくる。


### 基本構成

基本構成は以下

```
ansible/
├── ansible.cfg      # Ansible 自体の設定（inventory の場所など）
├── inventory/       # インフラの一覧など
│   └── hosts.yaml
└── playbooks/       # 実際に流すコマンドなど
    └── hogehoge.yaml
```

ansible.cfg

```conf
[defaults]
inventory = inventory/hosts.yaml
roles_path = roles
retry_files_enabled = False
```

hosts.yaml

```yaml
all:
  vars:
    ansible_user: xxxx
  children:
    k8s_nodes:
      children:
        control_plane:
          hosts:
            srkr001:
              ansible_host: xxx.xxx.xxx.xxx
        workers:
          hosts:
            srkr002:
              ansible_host: xxx.xxx.xxx.xxx
            srkr003:
              ansible_host: xxx.xxx.xxx.xxx
            srkr004:
              ansible_host: xxx.xxx.xxx.xxx

```

のような感じに設定

```sh
ansible-playbook playbooks/k8s-node.yaml --check --diff -K
```

これで実際にplaybookを流した場合の結果が出る。

--checkでDry Run,　--diff で変更差分が表示される

面白い。iscsiを使えるように全てのserverにinstallしようとしてたけど、これは便利だ！

今後、ロールや一部のnodeにだけ入れるなどのことが必要になった場合にどうかくかなども必要に応じてみていきたい。