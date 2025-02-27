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