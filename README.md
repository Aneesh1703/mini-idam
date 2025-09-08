#Mini IDAM + PAM Platform

A minimal Identity & Access Management (IDAM) + Privileged Access Management (PAM) system written in **Go**.  
This project demonstrates user authentication with MFA (TOTP), role-based access control (RBAC), credential vault, and session logging — all without a frontend.

---

## Features

- **User Registration & Login**
  - Secure password hashing
  - TOTP-based MFA (Time-based OTP)
- **Role-Based Access Control (RBAC)**
  - Assign roles (admin, user, etc.)
- **Vault**
  - Store and retrieve encrypted secrets
- **Session Logging**
  - Track user actions and sessions
- **PostgreSQL Database**
  - Persistent storage for users, sessions, and vault entries
