CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    product_name VARCHAR(255) NOT NULL,
    product_description TEXT,
    product_images TEXT[],
    product_price DECIMAL(10, 2),
    compressed_product_images TEXT[]
);
