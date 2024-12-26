create table products
(
    id                    int unsigned auto_increment
        primary key,
    sku         varchar(255) not null,
    name         varchar(255) not null,
    category         varchar(255) not null,
    created_at            timestamp                                           null,
    updated_at            timestamp                                           null
);

create table prices
(
    id                   int unsigned auto_increment
        primary key,
    product_id             int unsigned   not null,
    original             int unsigned   not null,
    final int unsigned not null,
    discount_percentage TINYINT UNSIGNED not null,
    currency VARCHAR(10) not null,
    constraint fk_products_price
        foreign key (product_id) references products (id)
);