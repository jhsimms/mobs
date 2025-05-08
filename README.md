# MOBS - Managed Object Store

MOBS is a multi-tenant object storage system designed to provide isolated, secure storage for each tenant while abstracting away the underlying cloud provider implementation.

## Status

🚧 **Early Development** 🚧

## Overview

MOBS enables tenants to store, manage, and access their unstructured data (such as images, videos, backups, logs, and other file types) without managing the underlying infrastructure themselves.

Key features (planned):
- Secure multi-tenant isolation
- Cloud provider abstraction
- Scalable architecture
- Comprehensive access controls

## Implementation Choices

### Why I chose this

I implemented a simple CLI interface with local database storage to demonstrate the core tenant provisioning functionality. While I originally planned to incorporate LocalStack and automation for cloud infrastructure simulation, I pared back the MVP due to timing constraints. I prioritized getting a functional "API" (via CLI) in place that showcases good architectural patterns and creates a foundation we can start experimenting with immediately.

### Considerations, decisions, and assumptions made

- Local storage: Used CloverDB instead of AWS services for simplicity and zero-dependency setup to get moving quickly. Swap out later.
- Abstraction layers: Created clear interface boundaries to make future cloud migration straightforward
- Validation: Incorporated a variety of validations for simple demonstration
- CLI vs. API: Chose CLI for immediate usability without requiring a server component. I constructed the code to support an easy transition.
- Assumed single-user: Designed for local development use without authentication concerns
- Set up basic tooling and scripting for ease of development and developer happiness as it is important.
- Established a baseline for architectural decision records for lightweight documentation.

### Next few improvements I'd make
- Implement AWS storage adapter using DynamoDB and S3 behind the same interface
- Add comprehensive input validation and error handling
- Create a REST API with the same service layer
- Wire up Terraform generation for actual infrastructure provisioning
- Add user authentication and multi-tenancy
