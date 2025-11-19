
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT timezone('utc', now())
);


CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    developer VARCHAR(255),
    release DATE NOT NULL, 
    created_at TIMESTAMPTZ DEFAULT timezone('utc', now())
);


CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    rating SMALLINT NOT NULL CHECK (rating >= 1 AND rating <= 10),
    text_review TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT timezone('utc', now()),

    UNIQUE (user_id, game_id)
);