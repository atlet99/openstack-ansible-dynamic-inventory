# OpenStack Dynamic Inventory

***Inspired by the logic and design from the project [openstack_ansible_dynamic_inventory](https://github.com/MaksimRudakov/openstack_ansible_dynamic_inventory.git) by [Maxim Rudakov](https://github.com/MaksimRudakov).***

This project is a Go-based implementation of a dynamic inventory for OpenStack environments, designed to organize and filter servers into groups based on metadata tags. It is ideal for automating infrastructure management, particularly with tools like Ansible, where dynamic inventories are beneficial.

## Table of Contents

- [OpenStack Dynamic Inventory](#openstack-dynamic-inventory)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Installation](#installation)
  - [Configuration](#configuration)
    - [.env file](#env-file)
  - [Usage](#usage)
  - [Environment Variables](#environment-variables)
  - [Project Structure](#project-structure)
  - [Localization READMEs](#localization-readmes)
  - [License](#license)

## Features

* `Dynamic inventory generation:` automatically retrieves server details from OpenStack and structures them into a dynamic inventory.
* `Metadata filtering:` filters servers based on specified metadata, allowing selection of specific environments.
* `Grouping by metadata:` creates groups based on metadata, enabling organized management of servers by roles, environments, or other attributes.
* `Easy integration:` output in JSON format compatible with Ansible and similar tools.

## Installation

1. Clone this repository:
```shell
git clone https://github.com/atlet99/openstack-ansible-dynamic-inventory.git; cd openstack-ansible-dynamic-inventory
```
2. Install dependencies:
```shell
go mod tidy
```

## Configuration

Configuration is managed via environment variables, with an optional `.env` file. These variables include both OpenStack connection credentials and filtering options for the inventory. You can find all necessary variables in the [Environment Variables](#environment-variables) section below.

### .env file

Create `.env` file in the project root to store your environment variables (also, you can easily copy `.env.template`):
```shell
OS_AUTH_URL=https://your-openstack-url:5000/v3
OS_USERNAME=your-username
OS_PASSWORD=your-password
OS_PROJECT_DOMAIN_ID="default"
OS_REGION_NAME=your-region
OS_DOMAIN_NAME="default"
ENVIRONMENT_TAG=environment
ENVIRONMENT_VALUE=test
BASE_GROUP_NAME=prod-servers
```
***Note:*** ensure `.env` is included in `.gitignore` to avoid committing sensitive information.

## Usage

1. To list the full inventory:
```shell
go run cmd/main.go --list
```
2. To request details for a specific host (returns empty JSON for compatibility):
```shell
go run cmd/main.go --host <your-hostname>
```

## Environment Variables

Each environment variable serves a specific purpose in connecting to OpenStack and managing the dynamic inventory.

```plaintext
# OpenStack authentication URL, typically provided by your OpenStack environment.
# This URL includes the identity (Keystone) API endpoint for authentication.
OS_AUTH_URL=https://your-openstack-url:5000/v3

# Username for OpenStack API access.
# Ensure this user has the necessary permissions to list and manage resources.
OS_USERNAME=your-username

# Password for the specified OpenStack user.
OS_PASSWORD=your-password

# Domain ID for the OpenStack project, often "default" unless specified otherwise by the environment.
# This domain is associated with the project name where resources are managed.
OS_PROJECT_DOMAIN_ID="default"

# Region name within OpenStack, specifying the regional data center you are connecting to.
# Useful in multi-region setups to ensure the connection is directed to the correct data center.
OS_REGION_NAME=your-region

# Domain name for the OpenStack user, typically "default" unless a specific domain is required.
OS_DOMAIN_NAME="default"

# Metadata tag key used to filter servers in the inventory.
# The script will select only the servers that have this key-value pair in their metadata.
ENVIRONMENT_TAG=environment

# Metadata tag value that must match for a server to be included in the inventory.
# Together with ENVIRONMENT_TAG, it defines the environment (e.g., "test", "prod") to which the servers belong.
ENVIRONMENT_VALUE=test

# Base group name to which all selected servers are assigned in the inventory.
# This name organizes servers into a logical group for easy management.
BASE_GROUP_NAME=prod-servers
```

## Project Structure
```
openstack-ansible-dynamic-inventory
├── .env.template
├── .gitignore
├── LICENSE
├── README.md
├── cmd
│   └── main.go
├── go.mod
├── go.sum
├── localization
│   ├── kz-KZ
│   │   └── README-kz-KZ.md
│   └── ru-RU
│       └── README-ru-RU.md
└── pkg
    ├── inventory
    │   ├── inventory.go
    │   └── openstack.go
    └── utils
        └── json.go
```

## Localization READMEs

| Language                   |
| -------------------------- |
| [Russian](localization/ru-RU/README-ru-RU.md)|
| [Kazakh](localization/kz-KZ/README-kz-KZ.md) |


## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) for details.