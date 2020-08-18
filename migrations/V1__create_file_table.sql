create table if not exists processed_files (
    file_name text unique,
    process_date date,
    primary key (file_name)
);