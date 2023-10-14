CREATE TABLE delivery(
                         id SERIAL PRIMARY KEY,
                         name varchar(50),
                         phone varchar(12),
                         zip varchar(10),
                         city varchar(50),
                         address varchar(60),
                         region varchar(50),
                         email varchar(50)
);

CREATE TABLE item(
                     chrt_id int PRIMARY KEY,
                     track_number varchar(50),
                     price int,
                     rid varchar(50),
                     name varchar(100),
                     sale int,
                     size varchar(3),
                     total_price int,
                     nm_id int,
                     brand varchar(100),
                     status int
);

CREATE TABLE orders(
                       order_uid text PRIMARY KEY,
                       track_number varchar(50),
                       entry varchar(10),
                       locale varchar(2),
                       internal_signature varchar(10),
                       customer_id varchar(20),
                       delivery_service varchar(20),
                       shardkey varchar(20),
                       sm_id int,
                       date_created date,
                       oof_shard varchar(10),
                       delivery_id int,
                       payment_id text,
                       items_order_id int
);

CREATE TABLE items_in_order(
                               order_id text,
                               item_id int
);

CREATE TABLE payment(
                        transaction text PRIMARY KEY,
                        request_id varchar(50),
                        currency varchar(3),
                        provider varchar(10),
                        amount int,
                        payment_dt int,
                        bank varchar(20),
                        delivery_cost int,
                        goods_total int,
                        custom_fee int

);

ALTER TABLE orders ADD constraint delivery_fk foreign key (delivery_id) references delivery(id);
ALTER TABLE orders ADD constraint payment_fk foreign key (payment_id) references payment(transaction);

ALTER TABLE items_in_order ADD CONSTRAINT item_fk foreign key (item_id) references item(chrt_id);
ALTER TABLE items_in_order ADD CONSTRAINT order_fk foreign key (order_id) references orders(order_uid);