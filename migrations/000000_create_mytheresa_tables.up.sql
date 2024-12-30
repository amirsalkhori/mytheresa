CREATE TABLE products
(
    id          INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    sku         VARCHAR(255) NOT NULL UNIQUE,
    name        VARCHAR(255) NOT NULL,
    category    VARCHAR(255) NOT NULL,
    price       INT UNSIGNED NOT NULL,
    INDEX idx_sku (sku),            
    INDEX idx_name (name),             
    INDEX idx_category (category),        
    INDEX idx_price (price)        
);

CREATE TABLE discounts
(
    id          INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    type ENUM('category', 'sku') NOT NULL,
    identifier VARCHAR(50) NOT NULL,
    percentage FLOAT NOT NULL,
    INDEX idx_type (type), 
    INDEX idx_identifier (identifier), 
    INDEX idx_percentage (percentage)     
);

-- -- Inserting products
-- INSERT INTO products (sku, name, category, price) VALUES
-- ('000001', 'BV Lean leather ankle boots', 'boots', 89000),
-- ('000002', 'BV Lean leather ankle boots', 'boots', 99000),
-- ('000003', 'Ashlington leather ankle boots', 'boots', 71000),
-- ('000004', 'Naima embellished suede sandals', 'sandals', 79500),
-- ('000005', 'Nathane leather sneakers', 'sneakers', 59000);

-- -- Inserting discounts
-- INSERT INTO discounts (type, identifier, percentage) VALUES
-- ('sku', '000003', 15),
-- ('category', 'boots', 30);
