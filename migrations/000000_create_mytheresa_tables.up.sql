CREATE TABLE products
(
    id          INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    sku         VARCHAR(255) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    category    VARCHAR(255) NOT NULL,
    price       INT UNSIGNED NOT NULL,
    currency    VARCHAR(10) NOT NULL,
    INDEX idx_sku (sku),            
    INDEX idx_name (name),             
    INDEX idx_category (category),        
    INDEX idx_price (price)        
);

CREATE TABLE discounts
(
    id          INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    sku         VARCHAR(255) NOT NULL,
    category    VARCHAR(255) NOT NULL,
    percentage  VARCHAR(10) NOT NULL,
    INDEX idx_sku (sku),            
    INDEX idx_category (category), 
    INDEX idx_percentage (percentage)     
);