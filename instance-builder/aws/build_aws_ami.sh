#!/bin/bash
set -e

BOOTSTRAP_GPU_INSTANCE_TYPE=g5.xlarge
BOOTSTRAP_CPU_INSTANCE_TYPE=t2.micro
AWS_BASE_AMI_ID=ami-009604998d7aa26d4
AWS_REGION=us-east-1 # ami ID may change with region
AWS_USER=ubuntu #ec2-user

####################
### Section 1: Environment Setup #####
####################

# Load environment variables
if [ -f .env ]; then
    echo "Loading environment variables from .env"
    source .env
else
    echo "Warning: .env file not found, continuing with existing environment variables"
fi

echo "Starting AMI build at $(date)"

####################
### Section 2: Parse Arguments #####
####################

# Parse arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --cpu) INSTALL_TYPE="cpu"; shift ;;
        --gpu) INSTALL_TYPE="gpu"; shift ;;
        *) echo "Unknown parameter: $1"; exit 1 ;;
    esac
done

echo "Building AMI type: $INSTALL_TYPE"

####################
### Section 3: Validate Environment #####
####################

# Check for required environment variables
required_vars=("AWS_ACCESS_KEY_ID" "AWS_SECRET_ACCESS_KEY" "AWS_BASE_AMI_ID" 
               "AWS_SUBNET_ID" "AWS_SECURITY_GROUP_ID" "AWS_REGION")

if [ "$INSTALL_TYPE" = "cpu" ]; then
    required_vars+=("BOOTSTRAP_CPU_INSTANCE_TYPE")
fi
if [ "$INSTALL_TYPE" = "gpu" ]; then
    required_vars+=("BOOTSTRAP_GPU_INSTANCE_TYPE")
fi

for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        echo "Error: Required environment variable $var is not set"
        exit 1
    fi
done

####################
### Section 4: AWS CLI Setup #####
####################

# Install AWS CLI if needed
if ! command -v aws &> /dev/null; then
    echo "Installing AWS CLI..."
    curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
    apt-get install -y unzip
    unzip awscliv2.zip > /dev/null 2>&1
    ./aws/install
    rm -rf aws awscliv2.zip
fi

if ! aws sts get-caller-identity >/dev/null 2>&1; then
    echo "Error: Invalid AWS credentials"
    exit 1
fi

####################
### Section 5: Key Pair Management Functions #####
####################

generate_key_pair() {
    KEY_NAME="temp-key-$(date +%s)"
    KEY_FILE="/tmp/${KEY_NAME}.pem"
    
    # Create a new key pair
    aws ec2 create-key-pair \
        --key-name "${KEY_NAME}" \
        --query 'KeyMaterial' \
        --region $AWS_REGION \
        --output text > "${KEY_FILE}"
    
    chmod 400 "${KEY_FILE}"
    echo "${KEY_NAME}"
}

delete_key_pair() {
    local key_name=$1
    local key_file="/tmp/${key_name}.pem"
    
    aws ec2 delete-key-pair --key-name "${key_name}" --region $AWS_REGION
    rm -f "${key_file}"
}

####################
### Section 6: Version Management #####
####################

# Update version tracking
update_version_file() {
    local ami_type=$1
    local version_file="ami_versions_${ami_type}.txt"
    # Ensure the version file exists
    touch "$version_file"
    # Generate a simple AMI name based on the type and current timestamp
    ami_name="ami_benchmarks_${ami_type}_$(date '+%Y%m%d%H%M%S')"
    # Return the AMI name
    echo "$ami_name"
}


####################
### Section 7: AMI Creation #####
####################

# Launch EC2 instance
launch_instance() {
    local instance_type=$1
    local ami_type=$2
    
    echo "Launching instance for $ami_type AMI..."
    
    # Create temporary key pair
    KEY_NAME=$(generate_key_pair)

    echo "KEY_NAME value: '$KEY_NAME'"

    aws ec2 describe-key-pairs --key-names "$KEY_NAME" --region $AWS_REGION
    
    # Define cleanup function
    cleanup() {
        echo "Performing cleanup..."
        if [ -n "$INSTANCE_ID" ]; then
            echo "Terminating instance $INSTANCE_ID"
            aws ec2 terminate-instances --instance-ids "$INSTANCE_ID" --region $AWS_REGION || true
        fi
        if [ -n "$TEST_INSTANCE_ID" ]; then
            echo "Terminating test instance $TEST_INSTANCE_ID"
            aws ec2 terminate-instances --instance-ids "$TEST_INSTANCE_ID" --region $AWS_REGION || true
        fi
        if [ -n "$KEY_NAME" ]; then
            echo "Deleting key pair $KEY_NAME"
            delete_key_pair "$KEY_NAME" || true
        fi
    }

    trap cleanup EXIT

    # Launch instance with the key pair
    INSTANCE_ID=$(aws ec2 run-instances \
        --image-id "$AWS_BASE_AMI_ID" \
        --instance-type "$instance_type" \
        --subnet-id "$AWS_SUBNET_ID" \
        --security-group-ids "$AWS_SECURITY_GROUP_ID" \
        --key-name "$KEY_NAME" \
        --region "$AWS_REGION" \
        --associate-public-ip-address \
        --block-device-mappings '[{"DeviceName":"/dev/sda1","Ebs":{"VolumeSize":50}}]' \
        --query 'Instances[0].InstanceId' \
        --output text)

    echo "Waiting for instance $INSTANCE_ID to be running..."
    aws ec2 wait instance-running --instance-ids "$INSTANCE_ID" --region $AWS_REGION

    # Get instance public IP
    INSTANCE_IP=$(aws ec2 describe-instances \
        --instance-ids "$INSTANCE_ID" \
        --region $AWS_REGION \
        --query 'Reservations[0].Instances[0].PublicIpAddress' \
        --output text)

    wait_for_ssh_access
}

# Wait for SSH access to be available
wait_for_ssh_access() {
    echo "Waiting for SSH to be available..."
    for i in {1..20}; do
        if ssh -o StrictHostKeyChecking=no -o ConnectTimeout=15 \
        -i "/tmp/${KEY_NAME}.pem" "${AWS_USER}"@"${INSTANCE_IP}" exit 2>/dev/null; then
            echo "SSH connection established"
            break
        fi
        echo "Attempt $i of 20: SSH connection failed, retrying after 30 seconds..."
        if [ $i -eq 30 ]; then
            echo "Timeout waiting for SSH access"
            exit 1
        fi
        sleep 30
    done
}

# Install requirements on the instance
install_requirements() {
    local install_script=$1

    scp -o StrictHostKeyChecking=no -i "/tmp/${KEY_NAME}.pem" -r \
        "ec2/" "${AWS_USER}@${INSTANCE_IP}:/home/ubuntu/"

    echo "Installing the API"
    ssh -o StrictHostKeyChecking=no -i "/tmp/${KEY_NAME}.pem" "${AWS_USER}@${INSTANCE_IP}" \
        "cd /home/ubuntu/ec2/api && chmod +x install.sh && ./install.sh"
    
    echo "Installing the $INSTALL_TYPE requirements"
    ssh -o StrictHostKeyChecking=no -i "/tmp/${KEY_NAME}.pem" "${AWS_USER}@${INSTANCE_IP}" \
        "chmod +x /home/ubuntu/ec2/$INSTALL_TYPE/install.sh && cd /home/ubuntu/ec2/$INSTALL_TYPE && sudo ./install.sh"

    if [ $? -ne 0 ]; then
        echo "Installation failed"
        exit 1
    fi
}

# Setup API service
setup_api_service() {
    echo "Setting up entrypointAPI service..."

    # Verify service status
    sleep 10 && echo "Checking service status..."

    ssh -o StrictHostKeyChecking=no -i "/tmp/${KEY_NAME}.pem" "${AWS_USER}@${INSTANCE_IP}" \
        "sudo systemctl status api.service && \
        sudo journalctl -u api.service --no-pager -n 50"

    # Verify API service is running
    echo "Verifying API service is running..."
    ssh -o StrictHostKeyChecking=no -i "/tmp/${KEY_NAME}.pem" "${AWS_USER}@${INSTANCE_IP}" \
        "curl -X GET http://localhost:8001/health"

    echo "API service setup complete"
}

# Launch a small instance from the newly created AMI to test the API service
launch_test_instance() {
    local ami_id=$1
    
    echo "Launching test instance from AMI: $ami_id"
    TEST_INSTANCE_ID=$(aws ec2 run-instances \
        --image-id "$ami_id" \
        --instance-type "t2.nano" \
        --subnet-id "$AWS_SUBNET_ID" \
        --security-group-ids "$AWS_SECURITY_GROUP_ID" \
        --key-name "$KEY_NAME" \
        --region $AWS_REGION \
        --associate-public-ip-address \
        --block-device-mappings '[{"DeviceName":"/dev/sda1","Ebs":{"VolumeSize":50}}]' \
        --query 'Instances[0].InstanceId' \
        --output text)

    echo "Waiting for test instance to be running..."
    aws ec2 wait instance-running --instance-ids "$TEST_INSTANCE_ID" --region $AWS_REGION
    
    # Get instance public IP
    TEST_INSTANCE_IP=$(aws ec2 describe-instances \
        --instance-ids "$TEST_INSTANCE_ID" \
        --region $AWS_REGION \
        --query 'Reservations[0].Instances[0].PublicIpAddress' \
        --output text)
    
    echo "Test instance IP: $TEST_INSTANCE_IP"
    
    # Wait for SSH to be available
    echo "Waiting for SSH to be available on test instance..."
    for i in {1..20}; do
        if ssh -o StrictHostKeyChecking=no -o ConnectTimeout=15 \
        -i "/tmp/${KEY_NAME}.pem" "${AWS_USER}"@"${TEST_INSTANCE_IP}" exit 2>/dev/null; then
            echo "SSH connection established to test instance"
            break
        fi
        echo "Attempt $i of 20: SSH connection failed, retrying after 15 seconds..."
        if [ $i -eq 20 ]; then
            echo "Timeout waiting for SSH access to test instance"
            terminate_test_instance
            exit 1
        fi
        sleep 15
    done
    
    # Test API health endpoint
    echo "Testing API health endpoint..."
    if ! ssh -o StrictHostKeyChecking=no -i "/tmp/${KEY_NAME}.pem" "${AWS_USER}"@"${TEST_INSTANCE_IP}" \
        "curl -s -f -X GET http://localhost:8001/health"; then
        # Terminate the test instance
        terminate_test_instance
        echo "API verification failed. AMI is not valid."
        exit 1
    fi
    
    echo "API health check passed. AMI is valid."
    # Terminate the test instance
    terminate_test_instance
    return 0
}

# Helper function to terminate test instance
terminate_test_instance() {
    if [ -n "$TEST_INSTANCE_ID" ]; then
        echo "Terminating test instance $TEST_INSTANCE_ID"
        aws ec2 terminate-instances --instance-ids "$TEST_INSTANCE_ID" --region $AWS_REGION || true
    fi
}

# Create AMI from instance
create_ami() {
    local ami_type=$1
    
    echo "Attempt to create AMI image..."
    AMI_NAME=$(update_version_file "$ami_type")

    echo "AMI_NAME value: '$AMI_NAME'"
    if [ -z "$AMI_NAME" ]; then
        echo "Error: AMI_NAME is empty, cannot create image"
        exit 1
    fi

    AMI_ID=$(aws ec2 create-image \
        --instance-id "$INSTANCE_ID" \
        --region $AWS_REGION \
        --name "$AMI_NAME" \
        --description "Custom AMI with ${ami_type} requirements - ${AMI_NAME}" \
        --query 'ImageId' \
        --output text)
    
    echo "AMI_ID value: '$AMI_ID'"
    wait_for_ami_creation "$AMI_ID"
    
    echo "${ami_type}: ${AMI_ID} ${AMI_NAME} $(date '+%Y-%m-%d %H:%M:%S')" >> "ami_versions_${ami_type}.txt"
    echo "Created $ami_type AMI: $AMI_NAME ($AMI_ID)"
    eval "${ami_type}_AMI_NAME='$AMI_NAME'"
    eval "${ami_type}_AMI_ID='$AMI_ID'"
    
    # Test the newly created AMI
    echo "Testing the newly created AMI..."
    launch_test_instance "$AMI_ID"
}

# Wait for AMI to be available
wait_for_ami_creation() {
    local ami_id=$1
    
    echo "Waiting for AMI to be available..."
    timeout=3600  # 1 hour timeout
    start_time=$(date +%s)

    while true; do
        status=$(aws ec2 describe-images --region $AWS_REGION --image-ids "$ami_id" --query 'Images[0].State' --output text)
        current_time=$(date +%s)
        if [ "$status" = "available" ]; then
            break
        elif [ "$status" = "failed" ]; then
            echo "AMI creation failed"
            exit 1
        elif [ $((current_time - start_time)) -gt $timeout ]; then
            echo "Timeout waiting for AMI creation"
            exit 1
        fi
        echo "Waiting for AMI creation... Current status: $status"
        sleep 30
    done
}

# Main function to orchestrate AMI creation
launch_and_create_ami() {
    local instance_type=$1
    local install_script=$2
    local ami_type=$3

    # Validate ami_type before proceeding
    if [ "$ami_type" != "cpu" ] && [ "$ami_type" != "gpu" ]; then
        echo "Error: Invalid AMI type '$ami_type' - must be 'cpu' or 'gpu'"
        exit 1
    fi

    echo "Creating $ami_type AMI..."
    
    # Launch instance
    launch_instance "$instance_type" "$ami_type"
    
    # Install requirements
    install_requirements "$install_script"
    
    # Setup API service
    setup_api_service
    
    # Create AMI
    create_ami "$ami_type"
}

####################
### Section 8: Main Execution #####
####################

# Build AMIs based on type
if [ "$INSTALL_TYPE" = "cpu" ]; then
    launch_and_create_ami "$BOOTSTRAP_CPU_INSTANCE_TYPE" "cpu/install.sh" "cpu"
elif [ "$INSTALL_TYPE" = "gpu" ]; then
    launch_and_create_ami "$BOOTSTRAP_GPU_INSTANCE_TYPE" "gpu/install.sh" "gpu"
fi

echo "AMI build completed at $(date)"
