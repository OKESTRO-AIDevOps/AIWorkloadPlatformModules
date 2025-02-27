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