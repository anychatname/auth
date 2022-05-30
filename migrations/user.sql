create table if not exists users (
    id serial unique not null,
    nickname varchar unique,
    email varchar unique,
    password varchar,
    picture varchar,
    created_at timestamp not null,

    primary key (id)
);

create table if not exists external_user_auth (
    id varchar unique not null,
    user_id integer not null,
    email varchar not null,
    platform varchar not null,
    picture varchar not null,
    created_at timestamp not null,

    primary key (id),
    foreign key (user_id) references users(id)
);

create table if not exists user_session (
    id varchar unique not null,
    tmp_id varchar unique,
    user_id integer not null,
    logged_at timestamp not null,
    last_seen_at timestamp not null,
    logged_with varchar,
    actived boolean,

    primary key (id),
    foreign key (user_id) references users(id)
);

create unique index idx_user_id_actived on user_session(user_id, actived);

create table if not exists events (
	id serial unique not null,
	name varchar unique,
	created_at timestamp default now()
);

create table if not exists sudo (
	id serial unique not null,
	session_id varchar not null,
	duration_in_secs integer not null,
	created_at timestamp not null,

	primary key(id),
    foreign key (session_id) references user_session(id)
);

create table if not exists sudo_events (
	id serial unique not null,
	sudo_id integer not null,
	event_id integer not null, 
	created_at timestamp not null,

	primary key(id),
    foreign key (sudo_id) references sudo(id),
    foreign key (event_id) references events(id)
);
