CREATE TABLE IF NOT EXISTS users
(
    id      bigserial   not null primary key,
    name    text unique not null,
    balance bigint      not null default 0 check (balance >= 0)
);

CREATE TYPE quest_type as ENUM ('usual', 'random');

CREATE TABLE IF NOT EXISTS quests
(
    id          bigserial not null primary key,
    name        text      not null unique,
    description text      not null,
    cost        bigint    not null check (cost >= 0 and cost <= 1000),
    type        quest_type not null
);

CREATE TABLE IF NOT EXISTS balance_history
(
    id      bigserial not null primary key,
    user_id bigint    not null references users (id) on delete cascade,
    quest_id bigint    null references quests (id) on delete SET NULL,
    created timestamp not null default now(),
    balance bigint    not null,
    CONSTRAINT quest_unique UNIQUE NULLS NOT DISTINCT (user_id, quest_id)
);
