# **ML Workload Resource Stress Test API**

## **Overview**

This project provides a **stress test API** to measure the load and performance of key system resources during **ML Workload Pipeline execution**. It supports controlled testing for **CPU, Memory, Disk I/O, Network I/O, and GPU** resources.

---

## **Table of Contents**

1. [Base Code](#base-code)
2. [Features](#features)
3. [API Configuration](#api-configuration)
4. [Execution Guide](#execution-guide)
5. [Curl Test Scripts](#curl-test-scripts)
6. [Installation Steps](#installation-steps)

## **Base Code**

- **GitHub Repository**: [KetiOps Hybrid-Cloud](https://github.com/ketiops/Hybrid-Cloud/tree/main)

- **Code Summary**:
  - Provides stress testing for resource usage in **ML Workload Pipelines**.
  - Supported resources: CPU, Memory, Disk I/O, Network I/O, GPU.

---

## **Features**

- **API for Stress Testing**:
  - Modular APIs for each resource (CPU, GPU, Memory, Disk, Network).
  - **Arguments** allow control over stress duration and intensity.

- **Sequential and Parallel Execution**:
  - Stress tests can run **sequentially** or **in parallel** using **multi-processing**.

- **All-in-One Stress Test**:
  - Combine multiple resource tests in a single API call.

  ## **API Configuration**

### **CPU Stress API**
- **Endpoint**: `/cpu_stress`
- **Arguments**:
  - `duration` (int): Stress duration in seconds.
  - `cpu_num` (int): Number of CPUs to stress.
  - `percentage` (int): CPU utilization percentage.

### **GPU Stress API**
- **Endpoint**: `/gpu_stress`
- **Arguments**:
  - `duration` (int): Stress duration in seconds.

### **Memory Stress API**
- **Endpoint**: `/memory_stress`
- **Arguments**:
  - `duration` (int): Stress duration in seconds.
  - `mem_amount` (int): Memory usage in MB.

### **Disk I/O Stress API**
- **Endpoint**: `/disk_stress`
- **Arguments**:
  - `duration` (int): Stress duration in seconds.
  - `size_mb` (int): Size of data for disk I/O in MB.

### **Network I/O Stress API**
- **Endpoint**: `/network_stress`
- **Arguments**:
  - `duration` (int): Stress duration in seconds.
  - `net_url` (string): Target URL.
  - `net_port` (int): Target port.
  - `network_mod` (string): Network mode name.

### **All-in-One Stress API**
- **Endpoint**: `/all_in_one`
- **Arguments**:
  - `cpu_stress`, `gpu_stress`, `memory_stress`, `disk_stress`, `network_stress` (bool): Enable stress testing for each resource.
  - Resource-specific arguments as required.

  ## **Execution Guide**

### **Environment**

- Two directories:
  1. `new_resource_stress`: Resource stress APIs.
     - Location: `10.0.1.112:~/new_resource_stress`
  2. `network_stress_flask`: Network stress external API.
     - Location: `10.0.2.193:~/network_stress_flask`

### **Steps**

#### 1. **Build Docker Image**
```bash
# Navigate to 'new_resource_stress' directory and build Docker image
docker build -t resource-stress:0.0.3 -f Dockerfile .
```

#### 2. **Run Resource Stress Container**
```bash
docker run -it --rm -p 5000:5000 --gpus all resource-stress:0.0.3
```

#### 3. **Run Network Stress Flask**
```bash
# Navigate to 'network_stress_flask' directory and start Flask app
python3 app.py --port 5001
```

## **Curl Test Scripts**

The following scripts are available to test individual APIs:

1. **CPU Stress Test**: `cpu_stress_curl.sh`
    ```bash
    ./cpu_stress_curl.sh
    ```

2. **GPU Stress Test**: `gpu_stress_curl.sh`
    ```bash
    ./gpu_stress_curl.sh
    ```

3. **Memory Stress Test**: `memory_stress_curl.sh`
    ```bash
    ./memory_stress_curl.sh
    ```

4. **Disk Stress Test**: `disk_stress_curl.sh`
    ```bash
    ./disk_stress_curl.sh
    ```

5. **Network Stress Test**: `network_stress_curl.sh`
    ```bash
    ./network_stress_curl.sh
    ```

6. **All-in-One Stress Test**: `all_in_one_stress_curl.sh`
    ```bash
    ./all_in_one_stress_curl.sh
    ```