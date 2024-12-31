CREATE TABLE categories
(
    id          INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,    
    INDEX idx_name (name)
);

CREATE TABLE products
(
    id          INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    sku         VARCHAR(255) NOT NULL UNIQUE,
    name        VARCHAR(255) NOT NULL,
    price       INT UNSIGNED NOT NULL,
    category_id INT UNSIGNED NOT NULL,
    INDEX idx_sku (sku),            
    INDEX idx_name (name),             
    INDEX idx_price (price),      
    constraint fk_categories_products
        foreign key (category_id) references categories (id)  
);

CREATE TABLE discounts
(
    id   INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    sku VARCHAR(50) NOT NULL,
    category_id INT UNSIGNED NOT NULL,
    percentage INT UNSIGNED NOT NULL,
    INDEX idx_percentage (percentage),
    INDEX idx_sku (sku),
    constraint fk_categories_discounts
        foreign key (category_id) references categories (id)       
);

