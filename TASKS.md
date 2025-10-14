## ✅ Backend Tasks (Go)

### 🏗️ Setup

- [x] Initialize Go project (`backend/`)
- [x] PostgreSQL setup & schema
- [x] Redis setup
- [x] `.env` config (DB, Redis, JWT, etc.)
- [] DB & Redis connection modules

---

### 🔐 Authentication

- [x] User registration & login
- [x] Password hashing (bcrypt)=-
- [x] JWT generation & middleware
- [] Email verification
- [x] Password reset flow

---

### 📡 Monitor Management

- [ ] CRUD endpoints for monitors
- [ ] Pause/resume monitor
- [ ] Link monitors to user/team

### 🔍 Monitor Checks

- [ ] HTTP checker (timeout, SSL, status/body validation)
- [ ] Store results in DB
- [ ] Measure response time
- [ ] Multi-region support via env config

### ⏱️ Scheduler & Queue

- [ ] Job queue setup (BullMQ or Go alternative)
- [ ] Worker to run checks and store results
- [ ] Add/update/remove jobs on monitor change

### ⚠️ Incident Handling

- [ ] Detect failures/recoveries
- [ ] Create/resolve incidents
- [ ] Track incident history & affected regions
- [ ] Maintenance mode support

### 📣 Alerts

- [ ] Email alerts via Resend
- [ ] Alert channels API (email, webhook, Slack)
- [ ] Link channels to monitors
- [ ] Alert rate-limiting / cooldown
- [ ] Webhook delivery with signatures

### 📄 Status Pages

- [ ] Status pages CRUD
- [ ] Public status page endpoint
- [ ] Add/remove monitors

### 👥 Teams & Access Control

- [ ] Teams: create, invite, accept
- [ ] Role-based permissions (owner/admin/member/viewer)
- [ ] Team switch support

### 💳 Billing

- [ ] Stripe integration (checkout + webhooks)
- [ ] Enforce plan limits (monitors, intervals, regions)
- [ ] Plans config (free, pro, enterprise)

### 🔑 API Access

- [ ] API key generation & storage
- [ ] API key authentication middleware
- [ ] Scoped permissions
- [ ] Rate limiting per key

### 📝 Incident Tools

- [ ] Acknowledge incident
- [ ] Add notes
- [ ] Postmortem fields (cause, resolution)
- [ ] Notify team on changes

### 👤 User Settings

- [ ] Update profile (name, email, password)
- [ ] Notification preferences (critical only, quiet hours)
- [ ] Timezone settings
- [ ] Delete account

### 🚀 Onboarding & UX

- [ ] Send welcome email
- [ ] Onboarding steps (create monitor, add alerts, view dashboard)
- [ ] Seed demo data

### 🧪 Testing

- [ ] Unit tests
- [ ] Integration tests
- [ ] Load tests (simulate 1K+ monitors)

### 📚 Documentation

- [ ] OpenAPI spec
- [ ] API reference docs
- [ ] Dev docs (setup, deployment)
- [ ] FAQ / troubleshooting

---

### ⚙️ DevOps & Deployment

- [ ] Deploy API & workers (Railway/Render)
- [ ] Set up health checks
- [ ] Add monitoring/logging (Sentry, etc.)
- [ ] Automated DB backups

### ✅ Launch Readiness

- [ ] All critical features functional
- [ ] Security review done
- [ ] Payments tested
- [ ] Email delivery verified
- [ ] Performance acceptable
