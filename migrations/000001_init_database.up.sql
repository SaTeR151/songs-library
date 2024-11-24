CREATE TABLE IF NOT EXISTS songs(
    id SERIAL NOT NULL PRIMARY KEY,
    song text NOT NULL,
    name_group text NOT NULL,
    release_date DATE,
    text TEXT,
    link text
);
SET datestyle = "German, DMY";
CREATE INDEX IF NOT EXISTS songs_name ON songs (song);
CREATE INDEX IF NOT EXISTS songs_group ON songs (name_group);
CREATE INDEX IF NOT EXISTS songs_release_date ON songs (release_date);