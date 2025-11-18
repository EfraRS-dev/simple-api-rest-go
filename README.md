# Simple API REST - Go

API REST para gestión de órdenes construida con Go, Gin y PostgreSQL.

## 🚀 Inicio Rápido

### 1. Ejecutar con Docker Compose

```bash
docker-compose up --build
```

La API estará disponible en `http://localhost:5001`

### 2. Variables de Entorno

Las variables de entorno ya están configuradas en `docker-compose.yaml`:

```env
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=orders_db
```

Para ejecución local (sin Docker), crea un archivo `.env`:

```env
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=orders_db
```

## 📍 Rutas Disponibles

### Health Check

```http
GET /health
```

### Obtener todas las órdenes

```http
GET /api/orders
```

### Crear nueva orden

```http
POST /api/orders
Content-Type: application/json

{
  "customer_id": "C123",
  "items": [
    {
      "product_id": "P001",
      "quantity": 2,
      "price": 50.00
    },
    {
      "product_id": "P002",
      "quantity": 1,
      "price": 100.00
    }
  ]
}
```

**Respuesta:**

```json
{
  "order_id": 1,
  "customer_id": "C123",
  "total_amount": 200.00,
  "items_count": 2,
  "processing_date": "2025-11-17T10:30:45Z"
}
```

### Eliminar orden

```http
DELETE /api/orders/:id
```

## 🛠️ Comandos Útiles

```bash
# Ver logs
docker-compose logs -f

# Detener servicios
docker-compose down

# Detener y eliminar volúmenes
docker-compose down -v

# Reconstruir solo la API
docker-compose up --build api
```

## 📦 Base de Datos

PostgreSQL expuesto en puerto `5433` (host) y `5432` (contenedor).

Conexión desde host:

```bash
psql -h localhost -p 5433 -U postgres -d orders_db
```
