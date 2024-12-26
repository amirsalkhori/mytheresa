CREATE TABLE products
(
    id          INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    sku         VARCHAR(255) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    category    VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP NULL,
    updated_at  TIMESTAMP NULL,
    INDEX idx_sku (sku),            
    INDEX idx_name (name),             
    INDEX idx_category (category)        
);

CREATE TABLE prices
(
    id                    INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    product_id            INT UNSIGNED NOT NULL,
    original              INT UNSIGNED NOT NULL,
    final                 INT UNSIGNED NOT NULL,
    discount_percentage   VARCHAR(10) NULL,
    currency              VARCHAR(10) NOT NULL,
    CONSTRAINT fk_products_price
        FOREIGN KEY (product_id) REFERENCES products (id),
    INDEX idx_original (original),       
    INDEX idx_final (final)             
);
