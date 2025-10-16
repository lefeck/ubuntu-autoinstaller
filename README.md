<h1 align="center">UbuntuAutoInstaller</h1>

**UbuntuAutoInstaller** is a **web-based Ubuntu Autoinstall ISO generator** built with **Go + Gin**.  
It supports **cloud-init** automation, real-time validation, and visual configuration of system parameters.

## 📺 Video Tutorial

🎥 **[Watch the tutorial](https://youtu.be/Z3Pqv76VJcE)** - Learn how to use UbuntuAutoInstaller step by step

## Features

- [x] Automatically encrypt plaintext passwords
- [x] Import ISO files (upload custom ISO or download directly from the internet)
- [x] Generate cloud-init configurations
- [x] Validate user-data format before ISO build
- [x] Advanced networking and storage layouts
- [x] User and SSH key management
- [x] HWE kernel support and ISO integrity verification
- [x] Add support for downloading and embedding installation packages into the the ISO
- [ ] Add support for building local applications directly into the ISO


## Install

### Requirements

- **Go 1.24.5+**
- Linux(Ubuntu Linux recommended)
- Tools: `xorriso`, `p7zip/7zip`
- Network access to fetch Ubuntu ISOs

### Run Locally

Start by cloning the repository:
```bash
git clone https://github.com/lefeck/ubuntu-autoinstaller.git
cd ubuntu-autoinstaller
make run
```
You can also build using make build, which will compile in the web assets so that AutouISO can be run from anywhere:

```bash
make build
./ubuntu-autoinstaller
```
The Makefile provides several targets:
* build: Builds the project and places the binary in the current directory.
* run: Run the project.
* test: Run tests.
* deps: Install project dependencies.
* vet: Run static analysis.
* fmt: Formatting
* package: Archive artifacts (zip for Linux, tar.gz for others)
* clean: Clean up the project.
* help: Display help.
* docker-build: Build the project and place the binary in the current directory.
* docker-run: Run the project.
* docker-push: Push the project to the registry.
* compose-up: Starts the project using docker-compose.
* compose-down: Stops the project using docker-compose.
* compose-logs: View the logs of the project using docker-compose.


### Docker images(Recommended)

**Suggestion:**

* The version of your Docker image must match the version of the image you want to build; otherwise, issues such as installation failures may occur.

Docker images are available on [ACR](https://cr.console.aliyun.com) or [Docker Hub](https://hub.docker.com/).

#### Quick Start with Docker

**Step 1: Pull the image**

```bash
# Pull from Docker Hub (Ubuntu 22.04)
docker pull jetfuls/ubuntu-autoinstaller:1.0-ubuntu22.04

# Or pull from Aliyun ACR
docker pull crpi-g7nxbvns4i9rnvaf.cn-hangzhou.personal.cr.aliyuncs.com/jetfuls/ubuntu-autoinstaller:1.0-ubuntu22.04
```

**Step 2: Run the container**

You can launch a Autouiso container for trying it out with

```
# ubuntu-autoinstaller ubuntu 20.04
docker run -itd -p 8080:8080 --name ubuntu-autoinstaller-20.04  jetfuls/ubuntu-autoinstaller:1.0-ubuntu20.04
# or
docker run -itd -p 8080:8080 --name ubuntu-autoinstaller-20.04  crpi-g7nxbvns4i9rnvaf.cn-hangzhou.personal.cr.aliyuncs.com/jetfuls/ubuntu-autoinstaller:1.0-ubuntu20.04

# ubuntu-autoinstaller ubuntu 22.04
docker run -itd -p 8080:8080 --name ubuntu-autoinstaller-22.04  jetfuls/ubuntu-autoinstaller:1.0-ubuntu22.04
# or
docker run -itd -p 8080:8080 --name ubuntu-autoinstaller-22.04  crpi-g7nxbvns4i9rnvaf.cn-hangzhou.personal.cr.aliyuncs.com/jetfuls/ubuntu-autoinstaller:1.0-ubuntu22.04


# ubuntu-autoinstaller ubuntu 24.04
docker run -itd -p 8080:8080 --name ubuntu-autoinstaller-24.04  jetfuls/ubuntu-autoinstaller:1.0-ubuntu24.04
# or
docker run -itd -p 8080:8080 --name ubuntu-autoinstaller-24.04  crpi-g7nxbvns4i9rnvaf.cn-hangzhou.personal.cr.aliyuncs.com/jetfuls/ubuntu-autoinstaller:1.0-ubuntu24.04
```

**Step 3: Access the service**

Once the container is running, you can access:

- Web UI → [http://localhost:8080](http://localhost:8080)
- API Docs → [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

**Step 4: Manage the container**

```bash
# View container logs
docker logs ubuntu-autoinstaller-22.04

# Follow logs in real-time
docker logs -f ubuntu-autoinstaller-22.04

# Stop the container
docker stop ubuntu-autoinstaller-22.04

# Start the container again
docker start ubuntu-autoinstaller-22.04

# Remove the container
docker rm -f ubuntu-autoinstaller-22.04
```


### Building the Docker image 

You can build a docker image locally with the following commands:

```bash
# Build
make docker-build REGISTRY_USER=jetfuls APP_VERSION=1.0 UBUNTU_VERSION=22.04

# Run
make docker-run DOCKER_IMAGE=jetfuls/ubuntu-autoinstaller:1.0-ubuntu22.04
```
Notes:
- `REGISTRY_USER`: your Docker registry username
- `APP_VERSION`: your app version
- `UBUNTU_VERSION`: your Ubuntu version
- `DOCKER_IMAGE`: your Docker image name


### Docker Compose

If you are using `docker-compose.yml`, you can start the service stack with:

```
# Start (auto-builds the image if not present)
make compose-up # 22.04 
# or
make compose-up COMPOSE_FILE=docker-compose-24.04.yml

# View logs
make compose-logs

# Stop
make compose-down
```

### Health Check

To check the health of your Autouiso instance, use the following command:

```bash
curl -sf http://localhost:8080/health && echo OK || echo FAIL
```


### Supported Platforms

Considering that this project is designed for Ubuntu:

- In a non-Docker environment, it can only run on an Ubuntu host.

- In a Docker environment, it can run on any operating system.

### FAQ

- **Blank UI / 404?** → Check `STATIC_DIR` path
- **Health check fails?** → Inspect logs: `docker logs autouiso`
- **ISO build fails?** → Fix config validation errors, then check network


### Contribution

Contributions via Issues & PRs are welcome


### License

Licensed under the **MIT License** – see [LICENSE](LICENSE).

