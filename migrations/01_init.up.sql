begin;

create table if not exists counter_data
(
    id serial primary key,
    get_data timestamp,
    client_info text not null,
    operation varchar(15) not null
);

commit;