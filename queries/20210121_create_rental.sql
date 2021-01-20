CREATE TABLE IF NOT EXISTS returns (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    t_id INT NOT NULL,
    s_id INT NOT NULL,
    amount INT NOT NULL,
    date DATE NOT NULL,
    -- GOOD (300), LOW_BROKEN (402), MID_BROKEN (405), LOST (410),
    status INT NOT NULL,
    claim INT,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by INT,
    FOREIGN KEY ( modified_by ) REFERENCES users( id ) ON DELETE SET NULL,
    FOREIGN KEY ( t_id ) REFERENCES transactions( id ),
    FOREIGN KEY ( s_id ) REFERENCES shipment( id )
);

CREATE TABLE IF NOT EXISTS extends (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    next_snapshot INT NOT NULL,
    prev_snapshot INT NOT NULL,
    amount INT NOT NULL,
    FOREIGN KEY ( next_snapshot ) REFERENCES snapshot_item( id ),
    FOREIGN KEY ( prev_snapshot ) REFERENCES snapshot_item( id )
)