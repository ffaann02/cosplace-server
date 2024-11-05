BEGIN;

CREATE TABLE IF NOT EXISTS custom_posts (
    post_id VARCHAR(10) PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    description VARCHAR(200) NOT NULL,
    price_range_start FLOAT NOT NULL,
    price_range_end FLOAT NOT NULL,
    status VARCHAR(20) NOT NULL,
    anime_name VARCHAR(100) NOT NULL,
    tag VARCHAR(100) NOT NULL,
    created_by VARCHAR(10) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

CREATE TABLE IF NOT EXISTS custom_post_categories (
    custom_post_tag_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    post_id VARCHAR(10) NOT NULL,
    custom_category_id INT NOT NULL,
    FOREIGN KEY (post_id) REFERENCES custom_posts(post_id),
    FOREIGN KEY (custom_category_id) REFERENCES md_product_categories(category_id)
);

CREATE TABLE IF NOT EXISTS customs_ref_images (
    custom_image_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    post_id VARCHAR(10) NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    image_order INT NOT NULL,
    FOREIGN KEY (post_id) REFERENCES custom_posts(post_id)
);

CREATE TABLE IF NOT EXISTS services (
    service_id VARCHAR(10) PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description VARCHAR(200) NOT NULL,
    created_by VARCHAR(10) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

CREATE TABLE IF NOT EXISTS portfolios (
    portfolio_id VARCHAR(10) PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description VARCHAR(300) NOT NULL,
    created_by VARCHAR(10) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

CREATE TABLE IF NOT EXISTS portfolio_images (
    portfolio_image_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    portfolio_id VARCHAR(10) NOT NULL,
    image_description VARCHAR(100) NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    image_order INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (portfolio_id) REFERENCES portfolios(portfolio_id)
);

COMMIT;
