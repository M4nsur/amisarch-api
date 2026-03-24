CREATE SCHEMA studioapp;



CREATE TABLE studioapp.project_categories (
    id          SERIAL       PRIMARY KEY,
    version     BIGINT       NOT NULL DEFAULT 1,
    name        VARCHAR(100) NOT NULL CHECK (char_length(name) BETWEEN 1 AND 100),
    slug        VARCHAR(100) NOT NULL UNIQUE CHECK (slug ~ '^[a-z0-9-]+$'),
    sort_order  INT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE studioapp.projects (
    id          SERIAL       PRIMARY KEY,
    version     BIGINT       NOT NULL DEFAULT 1,
    category_id INT          REFERENCES studioapp.project_categories(id) ON DELETE SET NULL,
    title       VARCHAR(255) NOT NULL CHECK (char_length(title) BETWEEN 1 AND 255),
    slug        VARCHAR(255) NOT NULL UNIQUE CHECK (slug ~ '^[a-z0-9-]+$'),
    location    VARCHAR(255) CHECK (char_length(location) BETWEEN 2 AND 255),
    area_sqm    DECIMAL      CHECK (area_sqm > 0),
    year        INT          CHECK (year BETWEEN 1900 AND 2100),
    description TEXT         CHECK (char_length(description) BETWEEN 1 AND 10000),
    tasks       TEXT         CHECK (char_length(tasks) BETWEEN 1 AND 10000),
    tags        TEXT[],
    published   BOOLEAN      NOT NULL DEFAULT FALSE,
    sort_order  INT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE studioapp.project_images (
    id          SERIAL       PRIMARY KEY,
    version     BIGINT       NOT NULL DEFAULT 1,
    project_id  INT          NOT NULL REFERENCES studioapp.projects(id) ON DELETE CASCADE,
    url         VARCHAR(2048) NOT NULL CHECK (char_length(url) BETWEEN 10 AND 2048),
    type        VARCHAR(50)  NOT NULL DEFAULT 'photo'
                             CHECK (type IN ('photo', 'render', 'blueprint')),
    is_cover    BOOLEAN      NOT NULL DEFAULT FALSE,
    sort_order  INT
);


CREATE UNIQUE INDEX uq_project_cover
    ON studioapp.project_images (project_id)
    WHERE is_cover = TRUE;



CREATE TABLE studioapp.clients (
    id         SERIAL       PRIMARY KEY,
    version    BIGINT       NOT NULL DEFAULT 1,
    name       VARCHAR(255) NOT NULL CHECK (char_length(name) BETWEEN 2 AND 255),
    email      VARCHAR(255) NOT NULL UNIQUE
                            CHECK (
                                email ~ '^[^@\s]+@[^@\s]+\.[^@\s]+$'
                                AND char_length(email) BETWEEN 5 AND 255
                            ),
    phone      VARCHAR(50)  NOT NULL
                            CHECK (
                                phone ~ '^\+?[0-9\s\-\(\)]+$'
                                AND char_length(phone) BETWEEN 7 AND 50
                            ),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE studioapp.orders (
    id              SERIAL       PRIMARY KEY,
    version         BIGINT       NOT NULL DEFAULT 1,
    client_id       INT          NOT NULL REFERENCES studioapp.clients(id) ON DELETE RESTRICT,
    object_type     VARCHAR(50)  NOT NULL
                                 CHECK (object_type IN ('apartment', 'house', 'office')),
    area_sqm        DECIMAL      NOT NULL CHECK (area_sqm > 0),
    floors          INT          CHECK (floors BETWEEN 1 AND 200),
    rooms           INT          CHECK (rooms BETWEEN 1 AND 100),
    services        TEXT[]       NOT NULL,
    budget          DECIMAL      CHECK (budget >= 0),
    estimated_price DECIMAL      NOT NULL CHECK (estimated_price >= 0),
    status          VARCHAR(50)  NOT NULL DEFAULT 'new'
                                 CHECK (status IN ('new', 'negotiation', 'in_progress', 'review', 'done', 'cancelled')),
    notes           TEXT         CHECK (char_length(notes) BETWEEN 1 AND 10000),
    pdf_url         VARCHAR(2048) CHECK (char_length(pdf_url) BETWEEN 10 AND 2048),
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE studioapp.order_files (
    id         SERIAL        PRIMARY KEY,
    version    BIGINT        NOT NULL DEFAULT 1,
    order_id   INT           NOT NULL REFERENCES studioapp.orders(id) ON DELETE CASCADE,
    url        VARCHAR(2048) NOT NULL CHECK (char_length(url) BETWEEN 10 AND 2048),
    name       VARCHAR(255)  NOT NULL CHECK (char_length(name) BETWEEN 1 AND 255),
    type       VARCHAR(50)   NOT NULL
                             CHECK (type IN ('reference', 'contract', 'document')),
    created_at TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE TABLE studioapp.order_comments (
    id         SERIAL      PRIMARY KEY,
    version    BIGINT      NOT NULL DEFAULT 1,
    order_id   INT         NOT NULL REFERENCES studioapp.orders(id) ON DELETE CASCADE,
    text       TEXT        NOT NULL CHECK (char_length(text) BETWEEN 1 AND 10000),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE studioapp.order_status_history (
    id          SERIAL      PRIMARY KEY,
    version     BIGINT      NOT NULL DEFAULT 1,
    order_id    INT         NOT NULL REFERENCES studioapp.orders(id) ON DELETE CASCADE,
    status_from VARCHAR(50) CHECK (status_from IN ('new', 'negotiation', 'in_progress', 'review', 'done', 'cancelled')),
    status_to   VARCHAR(50) NOT NULL
                            CHECK (status_to IN ('new', 'negotiation', 'in_progress', 'review', 'done', 'cancelled')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- INDEXES
-- ============================================================

CREATE INDEX idx_projects_category      ON studioapp.projects (category_id);
CREATE INDEX idx_projects_published     ON studioapp.projects (published);
CREATE INDEX idx_project_images_project ON studioapp.project_images (project_id);
CREATE INDEX idx_orders_client          ON studioapp.orders (client_id);
CREATE INDEX idx_orders_status          ON studioapp.orders (status);
CREATE INDEX idx_order_files_order      ON studioapp.order_files (order_id);
CREATE INDEX idx_order_comments_order   ON studioapp.order_comments (order_id);
CREATE INDEX idx_order_history_order    ON studioapp.order_status_history (order_id);

-- ============================================================
-- AUTO-UPDATE updated_at
-- ============================================================

CREATE OR REPLACE FUNCTION studioapp.set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_projects_updated_at
    BEFORE UPDATE ON studioapp.projects
    FOR EACH ROW EXECUTE FUNCTION studioapp.set_updated_at();

CREATE TRIGGER trg_orders_updated_at
    BEFORE UPDATE ON studioapp.orders
    FOR EACH ROW EXECUTE FUNCTION studioapp.set_updated_at();
