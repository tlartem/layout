-- Мы имеем таблицы
CREATE TYPE status AS ENUM (
    'pending',
    'active',
    'inactive',
    'banned'
    );

CREATE TABLE IF NOT EXISTS profile
(
    id         UUID PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    name       TEXT,
    age        INT,
    status     status,
    verified   BOOLEAN,
    contacts   JSONB -- email, phone
);

CREATE TABLE IF NOT EXISTS property
(
    profile_id UUID REFERENCES profile (id) ON DELETE CASCADE,
    tags       TEXT[]
);

-- Добавление фейковых данных
BEGIN;

INSERT INTO profile (id, name, age, status, verified, contacts)
VALUES (gen_random_uuid(), 'John Doe', 30, 'active', TRUE,
        jsonb_build_object('email', 'john.doe@example.com', 'phone', '123-456-7890')),
       (gen_random_uuid(), 'Jane Smith', 25, 'pending', FALSE,
        jsonb_build_object('email', 'jane.smith@example.com', 'phone', '234-567-8901')),
       (gen_random_uuid(), 'Alice Johnson', 28, 'inactive', TRUE,
        jsonb_build_object('email', 'alice.johnson@example.com', 'phone', '345-678-9012')),
       (gen_random_uuid(), 'Bob Brown', 35, 'banned', FALSE,
        jsonb_build_object('email', 'bob.brown@example.com', 'phone', '456-789-0123'));

INSERT INTO property (profile_id, tags)
VALUES ((SELECT id FROM profile WHERE name = 'John Doe' LIMIT 1), ARRAY ['home', 'primary']),
       ((SELECT id FROM profile WHERE name = 'Jane Smith' LIMIT 1), ARRAY ['work', 'secondary']),
       ((SELECT id FROM profile WHERE name = 'Alice Johnson' LIMIT 1), ARRAY ['home', 'primary']),
       ((SELECT id FROM profile WHERE name = 'Bob Brown' LIMIT 1), ARRAY ['work', 'secondary']);

COMMIT;

-- Получить профиль из таблицы по ID профиля
EXPLAIN
SELECT *
FROM profile
WHERE id = 'd7cad69e-9759-4efe-a00f-c52fbc2ea689';

-- Получить теги из массива
SELECT tags
FROM property
WHERE profile_id = (SELECT id FROM profile WHERE name = 'John Doe');


-- Сделать append в массив
UPDATE property
SET tags = array_append(tags, 'new')
WHERE profile_id = (SELECT id FROM profile WHERE name = 'John Doe');

-- Удалить тег из массива
UPDATE property
SET tags = array_remove(tags, 'new')
WHERE profile_id = (SELECT id FROM profile WHERE name = 'John Doe');

-- Получить JSONB
SELECT contacts
FROM profile
WHERE id = (SELECT id FROM profile WHERE name = 'John Doe');

-- Получить поле из JSONB
SELECT contacts ->> 'email' AS email
FROM profile
WHERE id = (SELECT id FROM profile WHERE name = 'John Doe');

-- Удалить поле из JSONB
UPDATE profile
SET contacts = contacts - 'phone'
WHERE id = (SELECT id FROM profile WHERE name = 'John Doe');

-- Добавить поле в JSONB
UPDATE profile
SET contacts = contacts || jsonb_build_object('address', '123 Main St')
WHERE id = (SELECT id FROM profile WHERE name = 'John Doe');

-- Добаить индекс на property (profile_id)
CREATE INDEX IF NOT EXISTS idx_property_profile_id ON property (profile_id);

-- Запрос с JOIN
EXPLAIN ANALYSE
SELECT p.id, p.name, p.age, p.status, p.verified, p.contacts, pr.tags
FROM profile p
         JOIN property pr ON p.id = pr.profile_id
WHERE p.id = '34b4b762-4083-4cd3-8265-6a85857de745';