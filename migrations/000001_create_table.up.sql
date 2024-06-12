CREATE TABLE IF NOT EXISTS restaurants (
                                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(64),
    address VARCHAR(128),
    phone_number VARCHAR(32),
    description VARCHAR(128),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at BIGINT DEFAULT 0
    );

CREATE TYPE status AS ENUM('foydalanilmoqda', 'bo''sh', 'bron qilingan');

CREATE TABLE IF NOT EXISTS reservations (
                                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    restaurant_id UUID REFERENCES restaurants(id),
    reservation_time TIMESTAMP,
    status status,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at BIGINT DEFAULT 0
    );

CREATE TABLE IF NOT EXISTS menu (
                                    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id UUID REFERENCES restaurants(id),
    name VARCHAR(64),
    description VARCHAR(128),
    price DECIMAL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at BIGINT DEFAULT 0
    );

CREATE TABLE IF NOT EXISTS reservation_orders (
                                                  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reservation_id UUID REFERENCES reservations(id),
    menu_item_id UUID REFERENCES menu(id),
    quantity INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at BIGINT DEFAULT 0
    );