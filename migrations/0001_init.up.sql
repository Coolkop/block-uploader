create table if not exists chunk
(
    id         serial primary key,
    file_id    uuid        not null,
    storage_id int         not null,
    part_order int         not null,
    created_at timestamptz not null default now()
);

create index if not exists idx_chunk_file_id on chunk (file_id);