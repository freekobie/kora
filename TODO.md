# File Storage Platform – TODO Checklist

A feature roadmap for building a cloud-based file storage system using a Golang backend and React frontend.

---

## PHASE 1: Project Setup

### Base Environment
- [X] Initialize Go module for backend (go mod init)
- [X] Scaffold React app using Vite
- [X] Setup shared .env config for dev
- [X] Configure Docker / Docker Compose (Go + React + DB)

---

## PHASE 2: Backend (Golang)

### Core Features

#### Users & Authentication
- [X] User registration and login
- [ ] Password reset email flow

#### File Upload & Storage
- [ ] Chunked file upload (multi-part)
- [ ] Save metadata (size, type, checksum, owner_id, timestamps)
- [ ] Upload and manage file storage in Google Cloud Storage
- [ ] Validate and limit file types and size

#### Folders & Navigation
- [ ] Support nested folder hierarchy (parent_id relationships)
- [ ] Create, rename, move, delete folders
- [ ] Display folder tree

#### File Browsing
- [ ] List all files and folders under a directory
- [ ] Sort by name, date, size
- [ ] Filter/search by file type or name

#### Sharing & Permissions
- [ ] Share file or folder via public link
- [ ] Share with specific user (email) with view/edit access
- [ ] Create permission table (resource_id, user_id, access_level)

#### File Versioning
- [ ] Keep previous versions of uploaded files
- [ ] Allow rollback to earlier version

#### Trash Bin
- [ ] Soft delete files (moved to trash)
- [ ] Empty trash manually or after 30 days
- [ ] Restore from trash

#### Storage Quotas
- [ ] Track used storage per user
- [ ] Prevent upload when over quota
- [ ] Show storage usage in account settings

#### File Previews & Thumbnails
- [ ] Support preview for PDF, images, videos
- [ ] Use background job or 3rd-party service to generate previews

#### Downloads
- [ ] Single file download
- [ ] Multiple files as .zip archive

---

## PHASE 3: Frontend (React.js)

### UI Pages
- [ ] Login & Register
- [ ] Dashboard / File explorer
- [ ] Folder view and breadcrumb
- [ ] Upload modal / drag-and-drop area
- [ ] File preview screen
- [ ] Trash view
- [ ] Account settings (storage usage, password reset)

### Components & UX
- [ ] File/folder tree view
- [ ] Responsive layout (mobile/tablet)
- [ ] Upload progress bar and retry handling
- [ ] Rename, move, delete dialogs
- [ ] Context menu for files/folders

---

## PHASE 4: Admin Tools

- [ ] Admin dashboard
- [ ] View user accounts and storage usage
- [ ] Remove abusive/oversized files
- [ ] Monitor system errors and file storage metrics

---

## PHASE 5: Final Touches

### CLI Client
- [ ] Build a Go-based CLI tool for interacting with the platform
- [ ] Support login/auth via token or credentials
- [ ] Support file upload, download, delete
- [ ] Support folder listing and navigation
- [ ] Allow file sharing from the CLI
- [ ] Provide output in plain text or JSON

### Security
- [ ] Validate file content type and extensions
- [ ] Rate limiting on uploads/downloads
- [ ] CSRF protection (for cookies-based sessions)
- [ ] Audit logs (uploads, shares, deletes)

### Testing
- [ ] Unit tests for file APIs and auth
- [ ] Integration tests with GCS and DB
- [ ] Frontend tests (React Testing Library)
- [ ] E2E tests (Cypress/Playwright)

### DevOps & Monitoring
- [ ] Logging (e.g. zap + file rotation)
- [ ] Error tracking (e.g., Sentry)
- [ ] Metrics via Prometheus/Grafana
- [ ] Health check endpoint
- [ ] CI/CD pipeline with GitHub Actions
- [ ] Deploy backend to Google Cloud Compute Engine
- [ ] Deploy frontend to GCE or Firebase Hosting

---

## Documentation

- [ ] README for dev setup (Docker, envs)
- [ ] API docs with Swagger/OpenAPI
- [ ] Guide: how file upload and GCS storage works
- [ ] Guide: sharing and permissions model
- [ ] Deployment instructions (GCP, GitHub Actions)

