CREATE TABLE IF NOT EXISTS transactions (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    date DATE NOT NULL,
    customer VARCHAR(20) NOT NULL,
    -- SALES (1), RENTAL (2),
    type INT NOT NULL,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP,
    paid_date DATE,
    done_date DATE,
    modified_by INT,
    UNIQUE ( date, customer, type ),
    FOREIGN KEY ( modified_by ) REFERENCES users( id ) ON DELETE SET NULL,
    FOREIGN KEY ( customer ) REFERENCES customers( code )
);

CREATE TABLE IF NOT EXISTS snapshot (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    t_id INT NOT NULL,
    project INT,
    address VARCHAR(255) NOT NULL,
    total_price INT NOT NULL,
    deposit INT NOT NULL,
    discount INT,
    shipping_fee INT,
    FOREIGN KEY ( t_id ) REFERENCES transactions( id ),
    FOREIGN KEY ( project ) REFERENCES projects( id )
);

CREATE TABLE IF NOT EXISTS snapshot_item (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    t_id INT NOT NULL,
    item VARCHAR(20) NOT NULL,
    amount INT NOT NULL,
    price INT NOT NULL,
    claim INT,
    time_unit INT,
    duration INT,
    FOREIGN KEY ( t_id ) REFERENCES transactions( id ),
    FOREIGN KEY ( item ) REFERENCES items( code )
);

CREATE TABLE IF NOT EXISTS payments (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    t_id INT,
    name VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    amount INT NOT NULL,
    -- CASH (100), TRANSFER (200),
    method INT NOT NULL,
    -- DEBIT (100), CREDIT (200),
    account INT NOT NULL,
    accept_by INT,
    FOREIGN KEY ( t_id ) REFERENCES transactions( id ),
    FOREIGN KEY ( accept_by ) REFERENCES users( id ) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS shipment (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    t_id INT NOT NULL,
    i_id INT NOT NULL,
    amount INT NOT NULL,
    date DATE NOT NULL,
    deadline DATE NOT NULL,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by INT,
    FOREIGN KEY ( modified_by ) REFERENCES users( id ) ON DELETE SET NULL,
    FOREIGN KEY ( t_id ) REFERENCES transactions( id ),
    FOREIGN KEY ( i_id ) REFERENCES snapshot_item( id )
);
