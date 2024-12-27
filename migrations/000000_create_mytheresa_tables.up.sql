CREATE TABLE products
(
    id          INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    sku         VARCHAR(255) NOT NULL,
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