-- 000002_create_products_and_categories.up.sql
BEGIN;

CREATE TABLE IF NOT EXISTS md_product_categories (
    category_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    category_name VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS sellers (
    seller_id VARCHAR(10) NOT NULL PRIMARY KEY,
    shop_name VARCHAR(50) NOT NULL,
    verify BOOLEAN NOT NULL DEFAULT FALSE,
    accept_credit_card BOOLEAN NOT NULL DEFAULT FALSE,
    accept_qr_prompt_pay BOOLEAN NOT NULL DEFAULT FALSE,
    joined_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (seller_id) REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS products (
    product_id VARCHAR(10) NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price FLOAT NOT NULL,
    quantity INT NOT NULL,
    rent BOOLEAN DEFAULT FALSE,
    rent_return_date DATETIME,
    description VARCHAR(300),
    `condition` VARCHAR(50),
    size VARCHAR(5),
    region VARCHAR(50),
    created_by VARCHAR(10) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (created_by) REFERENCES sellers(seller_id)
);

CREATE TABLE IF NOT EXISTS product_categories (
    product_id VARCHAR(10) NOT NULL,
    category_id INT NOT NULL,
    PRIMARY KEY (product_id, category_id),  -- Composite primary key
    FOREIGN KEY (product_id) REFERENCES products(product_id),
    FOREIGN KEY (category_id) REFERENCES md_product_categories(category_id)
);


CREATE TABLE IF NOT EXISTS product_images (
    product_image_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    product_id VARCHAR(10) NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products(product_id)
);

COMMIT;