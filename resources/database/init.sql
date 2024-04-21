begin transaction;

CREATE TABLE IF NOT EXISTS commands  (
    id serial primary key ,
    code text,
    description text,
    created_at timestamp
);

CREATE TABLE IF NOT EXISTS completed_commands  (
    id serial primary key ,
    command_id int,
    result text,
    completed_at timestamp,
    status varchar(255),
    ppid int,
    foreign key(command_id) references commands(id) on delete set null
);

end transaction;