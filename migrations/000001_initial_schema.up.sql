-- Create sword_masters table
CREATE TABLE IF NOT EXISTS sword_masters (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    bio TEXT NOT NULL DEFAULT '',
    birth_year INTEGER,
    death_year INTEGER,
    image_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sword_masters_name ON sword_masters(name);

-- Create fighting_books table
CREATE TABLE IF NOT EXISTS fighting_books (
    id SERIAL PRIMARY KEY,
    sword_master_id INTEGER NOT NULL REFERENCES sword_masters(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    publication_year INTEGER,
    cover_image_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fighting_books_sword_master ON fighting_books(sword_master_id);
CREATE INDEX idx_fighting_books_title ON fighting_books(title);

-- Create chapters table
CREATE TABLE IF NOT EXISTS chapters (
    id SERIAL PRIMARY KEY,
    fighting_book_id INTEGER NOT NULL REFERENCES fighting_books(id) ON DELETE CASCADE,
    chapter_number INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_chapter_number_per_book UNIQUE(fighting_book_id, chapter_number)
);

CREATE INDEX idx_chapters_fighting_book ON chapters(fighting_book_id);
CREATE INDEX idx_chapters_number ON chapters(fighting_book_id, chapter_number);

-- Create techniques table
CREATE TABLE IF NOT EXISTS techniques (
    id SERIAL PRIMARY KEY,
    chapter_id INTEGER NOT NULL REFERENCES chapters(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    instructions TEXT NOT NULL DEFAULT '',
    video_url TEXT,
    thumbnail_url TEXT,
    order_in_chapter INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_techniques_chapter ON techniques(chapter_id);
CREATE INDEX idx_techniques_order ON techniques(chapter_id, order_in_chapter);
