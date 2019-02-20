## Available Database

Currently supported database:

- MongoDB (replica set)
- Redis (master and slave)
- MySQL (master and slave)

## Docker Image

Each image of data sources:

- Inherit from `ubuntu:16.04`, install database package from PPA
- Store data into a volume (like `/var/lib/redis`)
- Expose database on convention port + 1 (like `6380`)
- Expose an Nginx proxy to master node on convention port (like `6379`)

## In Container

Use supervisor as a root process, control children:

- `nginx` proxy connections to the master node
- `control-agent startup` to query role from DeployBeta, then `exec` to database main process
- `control-agent command-loop` to poll commands from DeployBeta, change database configuration
- And other processes if necessary

## Control Agent

`/usr/bin/control-agent` written in Bash.

## Replica

A data source may have two or more nodes, there are one master node and the rest are salves, they replica from the master.

Data source node required to ask it's role from DeployBeta, and report it's role changes back to DeployBeta.

If a data source node is a salve, it must proxy all connection to the master node.

## State Machine of DeployBeta

New node:

- If there is no master on the data source, the new node will be master
- If there is one master on the data source, the new node will be a salve of the master

Poll command:

- If the node's role is unknown, ask the node to report status
- If the node's reported role is master, but the data source's master is another node, ask node change to slave
- If the node's reported role is slave, but the data source's master is this node, ask node change to master

Cronjob:

- If node doesn't run on swarm, remove this node's metadata from etcd
