# **Go + MySQL Application with Docker Compose**

This project is for abundance which is glucose tracking application for final year project.

## **Prerequisites**

- Docker and Docker Compose installed on your machine.
- A `.env` file.

## **Getting Started**

### **1. Clone the Repository**

```bash
git clone https://github.com/your-username/your-repository.git
cd your-repository
```

### **2. Create a .env File**

Create a .env file in the root directory and fill in necessary credentials.

```bash
touch .env
```

Add the following contents:

```plaintext
PUBLIC_HOST=""
BACKEND_HOST=""

PORT=""

DB_USER=""
DB_PASSWD=""
DB_HOST=""
DB_PORT=""
DB_NAME=""

JWT_SECRET=""
JWT_EXP=9999

IS_PRODUCTION=false

SMTP_HOST=""
SMTP_PORT=""
SMTP_USERNAME=""
SMTP_PASSWORD=""
EMAIL_FROM=""
```

### **3. Build and Start Services**

Use Docker Compose to build and start the services.

```bash
docker-compose up --build
```

## **Common Commands**

### **Start the Project**

Start all services defined in `docker-compose.yml`.

```bash
docker-compose up
```

### **Stop the Project**

```bash
docker-compose down
```

### **Rebuild After Changes**

```bash
docker-compose up --build
```

### **Access MySQL Inside the Container**

```bash
docker exec -it mysql_db mysql -u root -p
```
