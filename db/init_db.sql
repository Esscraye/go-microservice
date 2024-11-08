-- Créer des tables
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    password VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    category VARCHAR(100),
    price DECIMAL
);

-- Insérer des utilisateurs
INSERT INTO users (name, email, password) VALUES
('User1', 'user1@example.com', 'password'),
('User2', 'user2@example.com', 'password');

-- Insérer des produits
INSERT INTO products (name, category, price) VALUES
('Product1', 'Category1', 100.0),
('Product2', 'Category2', 200.0);