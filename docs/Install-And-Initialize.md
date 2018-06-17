## Setup host server

Currently supported OS:

- Ubuntu 16.04 x64

Deploybeta required at least 2 host server (`deploy1` and `deploy2` below), make sure you already configured public key authentication on remote (`ssh-copy-id`) and configured they in `~/.ssh/config`:

```
ControlPath ~/.ssh/controlmasters/%r@%h:%p
ControlMaster auto
ControlPersist 60m

Host deploy1
  HostName 173.255.223.150
  User root

Host deploy2
  HostName 23.239.2.45
  User root
```

## Configure DNS

```
# used for provide Deploybeta Web UI
dashboard.deploybeta.io    A    173.255.223.150
dashboard.deploybeta.io    A    23.239.2.45

# public and private address
deploy1.deploybeta.io            A    173.255.223.150
deploy2.deploybeta.io            A    23.239.2.45
deploy1-internal.deploybeta.io   A    192.168.204.133
deploy2-internal.deploybeta.io   A    192.168.167.227

# used for internal use
es.deploybeta.io         CNAME    deploy1-internal.deploybeta.io
registry.deploybeta.io   CNAME    deploy2-internal.deploybeta.io

# used for provide users' websites
*.deploybeta.site   A    173.255.223.150
*.deploybeta.site   A    23.239.2.45

# used for etcd discovery
_etcd-server._tcp.deploybeta.io  SRV  1 10 2380 deploy1-internal.deploybeta.io
_etcd-server._tcp.deploybeta.io  SRV  2 10 2380 deploy2-internal.deploybeta.io
```

## Review ansible playbook

Write all nodes to `playbook/hosts`:

```
deploy1
deploy2
```

You may need to review or modify these files to match your requirement:

- `playbook/config.yml`
- `playbook/deploybeta.yml`
- `playbook/infrastructures.yml`

## Setup infrastructures

Login to each host and install `python-simplejson`:

```
apt-get install -y python-simplejson
```

Run `ansible-playbook infrastructures.yml` in `playbook`.

## Initialize docker swarm

- Run `docker swarm init` on one of your hosts (swarm leader), you may need to use `--advertise-addr` specify an intranet address.
- Run the command returned above on the rest hosts.
- Run `docker node promote <hostname>` for each host on swarm leader.

## Setup Deploybeta

Run `ansible-playbook deploybeta.yml` in `playbook`.

## Special tips

### Linode

Private address:

Add a private address to your hosts and enable `network helper` in advanced settings.

`error creating vxlan interface: operation not supported`:

- Run `apt-get install linux-signed-generic-lts-wily` to rollback kernel.
- Use `Grub2` in boot Settings of advanced settings.
