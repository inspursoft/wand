# Wand

![wand-magic](./docs/images/magic.jpeg)

## Introduction

Wand provides quick and convenient experience for DevOps. By integrating two core components of DevOps for which Gogs and Jenkins, Wand can help you setup DevOps workflow automatically.

## Features
* Supports setup DevOps workflow with automation
* Supports Git-based repository management
* Supports Travis-CI with Jenkins
* Supports multi slave nodes with KVM virtual machine managed with registry service to Jenkins master
* Supports deployment as Docker containers with docker-compose

## Get Started

1. Get source code from Github.
2. Compile and build source code into Docker images.
3. Setup KVM virtual machine as multi slave nodes.
4. Config to start.
5. Start the service

## Typical Usage

1. Register user as Gogs account.

  * Open brower locate to the Gogs URL.
  
  * Sign up user to the Gogs.

2. Create repository with registered user.

  * Sign in Gogs with registered user.

  * Discover current repositories to the current user.
  
  * Create repository to the current user.

3. Checkout source codes from created repository.
  
  * Copy clone URL from repository
  
  * Clone repository into local directory  

4. Create travis.yml to arrange CI procedures.

  * Create travis.yml

5. Commit changes to Gogs repository.

  * Submit your changes to the repository

6. Check Jenkins CI working status after commit.

  * Jenkins will be running CI as travis.yml described
