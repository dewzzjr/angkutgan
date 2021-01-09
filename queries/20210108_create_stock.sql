CREATE TABLE IF NOT EXISTS stock (
    code VARCHAR(20) NOT NULL PRIMARY KEY,
    asset INT NOT NULL,
    inventory INT NOT NULL,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP,
    modified_by INT,
    FOREIGN KEY ( modified_by ) REFERENCES users( id ) ON DELETE SET NULL,
    FOREIGN KEY ( code ) REFERENCES items( code )
);

CREATE TABLE IF NOT EXISTS stock_history (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    t_id INT NOT NULL,
    r_id INT,
    item VARCHAR(20) NOT NULL,
    date DATE NOT NULL,
    amount INT NOT NULL,
    -- SOLD (100), 
    -- RESTOCK (200),
    -- LOW_BROKEN (302), MID_BROKEN (305), LOST (310),
    -- LOW_REPAIR (402), MID_REPAIR (405),
    type INT NOT NULL,
    FOREIGN KEY ( t_id ) REFERENCES transactions( id ),
    FOREIGN KEY ( r_id ) REFERENCES returns( id ),
    FOREIGN KEY ( item ) REFERENCES items( code )
);