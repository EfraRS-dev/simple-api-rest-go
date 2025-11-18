-- Crear tabla de órdenes
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    customer_id VARCHAR(50) NOT NULL,
    total_amount NUMERIC(10,2) NOT NULL,
    items_count INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Crear tabla de items de órdenes
CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL,
    product_id VARCHAR(50) NOT NULL,
    quantity INTEGER NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

-- Insertar datos de ejemplo solo si no existen
INSERT INTO orders (customer_id, total_amount, items_count)
SELECT 'C123', 200.00, 2
WHERE NOT EXISTS (SELECT 1 FROM orders WHERE customer_id = 'C123');

INSERT INTO items (order_id, product_id, quantity, price)
SELECT 1, 'P001', 2, 50.00
WHERE NOT EXISTS (SELECT 1 FROM items WHERE order_id = 1 AND product_id = 'P001')
UNION ALL
SELECT 1, 'P002', 1, 100.00
WHERE NOT EXISTS (SELECT 1 FROM items WHERE order_id = 1 AND product_id = 'P002');