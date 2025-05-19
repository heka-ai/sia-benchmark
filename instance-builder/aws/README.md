## Ami builder

This script builds AMIs for AWS.

### Usage
Export the following environment variables:
```bash
export AWS_ACCESS_KEY_ID="your_access_key_id"
export AWS_SECRET_ACCESS_KEY="your_secret_access_key"
export AWS_REGION="your_region"
export AWS_BASE_AMI_ID="your_base_ami_id"
export AWS_SECURITY_GROUP_ID="your_security_group_id"
export AWS_SUBNET_ID="your_subnet_id"
export AWS_USER="ubuntu" # ec2-user
export BOOSTRAP_GPU_INSTANCE_TYPE="g5.xlarge"
export BOOSTRAP_CPU_INSTANCE_TYPE="t2.micro"
```

Run the script with the desired install type:
```bash
./build_aws_ami.sh --cpu
./build_aws_ami.sh --gpu
```

The script will build the AMIs and print the AMI IDs to the console.

### AMI Versions
The AMI versions are stored in the `ami_versions_cpu.txt` and `ami_versions_gpu.txt` files.

### Entrypoint API (State: In Progress)

- [ ] [CRITICAL] Secure the entrypoint API with a secret key
- [ ] [CRITICAL] Restrict the start endpoint to a specific command to start benchmark (cpu) or deploy model (gpu)
- [X] Add a script to deploy the entrypoint API to the AMI

The entrypoint API is a FastAPI application that is started by systemd at instance startup. It is installed in the `/opt/entrypoint_api.py` file and the systemd service file is installed in the `/etc/systemd/system/entrypoint_api.service` file.

To test the entrypoint API, you can use the following command:
```bash
curl -X POST http://<instance-ip or DNS>:8001/start -H "Content-Type: application/json" -d '{"command": "echo 'Hello, World!'"}'
curl -X POST http://<instance-ip or DNS>:8001/get-file -H "Content-Type: application/json" -d '{"filename": "test.txt"}'
```
