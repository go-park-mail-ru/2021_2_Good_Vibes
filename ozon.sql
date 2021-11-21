DROP TABLE IF EXISTS customers;

CREATE TABLE customers (
    id serial NOT NULL PRIMARY KEY,
    name varchar(30),
    email varchar(30),
    password varchar(100)
);


DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
CREATE TABLE IF NOT EXISTS categories (
    id serial NOT NULL PRIMARY KEY,
    name varchar(100) not null ,
    lft int not null ,
    rgt int not null
);

INSERT INTO categories(name, lft, rgt)
VALUES ('ALL_THINGS', 1, 18),
       ('CLOTHES', 2, 13),
       ('ELECTRONICS', 14, 17),
       ('CLOTHES_MEN', 3, 12),
       ('CLOTHES_UP_MEN', 4, 5),
       ('SHOES_MEN', 6, 11),
       ('SNICKERS_MEN', 7, 10),
       ('SNICKERS_ADIDAS_MEN', 8, 9),
       ('PHONES', 16, 21);

CREATE TABLE IF NOT EXISTS products (
    id serial NOT NULL PRIMARY KEY,
     image varchar(50),
     name varchar(50) not null ,
      price numeric(10,2),
      rating numeric(2,1),
      category_id int not null references categories(id)
);

