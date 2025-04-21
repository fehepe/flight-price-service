## AWS‑Native HTTPS/TLS + ECS + Redis Integration


### 1. Request a Public Certificate in ACM

1. **AWS Console → Certificate Manager → Request certificate**  
2. **Request a public certificate**, enter your domain (`api.flight.example.com`), choose **DNS validation**, and finish.  
3. Add the provided CNAME in your DNS zone and wait for **Issued** status.

---

### 2. Create an Internet‑Facing Application Load Balancer

1. **EC2 → Load Balancers → Create → Application Load Balancer**  
2. Name it (e.g. `flight‑alb`), Scheme **internet‑facing**, Listeners **HTTP (80)** and **HTTPS (443)**.  
3. Attach Security Group allowing **80/443** inbound.  
4. Create or select a Target Group (HTTP/8080, health check `/health`).  
5. Finish creation.

---

### 3. Redirect HTTP → HTTPS

1. In the ALB’s **Listeners** tab, edit port 80 → **Redirect** to HTTPS 443 (status code 301).  

---

### 4. Deploy to ECS

#### a. Cluster & Task

- **ECS → Clusters → Create Cluster** (Fargate)  
- **Task Definitions → Create**  
  - Launch type: **Fargate**  
  - Container image: `<repo>/flight-price-service:<tag>`  
  - Port mapping: container 8080 → host 8080  
  - Add env var `REDIS_URL` (filled later)

#### b. Service

- **Clusters → [your cluster] → Services → Create**  
- Launch type: **Fargate**, select your Task Definition  
- Service name: `flight-service`, tasks: as needed  
- Load balancer: select `flight‑alb` HTTPS listener → forward to your target group  
- Create service

---

### 5. Point DNS at the ALB

- **Route 53 → Hosted zones → Create record**  
- Name: `api.flight.example.com` → Alias to your ALB  
- Save

---

### 6. Provision ElastiCache Redis

1. **AWS Console → ElastiCache → Create**  
2. Choose **Redis**, engine version latest (e.g. 6.x), cluster mode **Disabled** (single‑node or small replica set)  
3. VPC & Subnet group: use same VPC as your ECS tasks  
4. Security Group: allow inbound on **6379** from your ECS tasks’ SG  
5. Finish and note the **Primary endpoint** (e.g. `my-redis.xxxxxx.use1.cache.amazonaws.com:6379`)

---

### 7. Wire Redis into Your Service

1. **ECS → Task Definitions → [your task] → Create new revision**  
2. In your container’s **Environment** section, add:
Name: REDIS_URL Value: redis://my-redis.xxxxxx.use1.cache.amazonaws.com:6379
3. Update (or redeploy) your ECS Service to use the new task revision.

---

### 8. Verify & Monitor

- **HTTPS health check**: `https://api.flight.example.com/health`  
- **Redis connectivity**: confirm cache hits/misses in your app logs or CloudWatch  
- **Certificate renewal**: managed automatically by ACM  
- **ElastiCache metrics**: view CPU, connections, evictions in CloudWatch  

---

✨ **Result:**  
Flight‑Price‑Service is now running on ECS, fronted by an AWS ALB with ACM TLS certificates, automatically redirecting HTTP→HTTPS, and backed by a secure ElastiCache Redis cluster for caching.  
