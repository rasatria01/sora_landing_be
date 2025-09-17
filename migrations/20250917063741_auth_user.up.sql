create table public.users
(
    id           varchar(27)                                     not null
        primary key,
    name         text                                            not null,
    created_at   timestamp with time zone default now()          not null,
    updated_at   timestamp with time zone default now()          not null,
    roles        text[]                                          not null,
    email        text,
    status       text                     default 'Active'::text not null
);
create unique index users_email_uindex
    on public.users (email);

create table public.authentications
(
    id               varchar(27)                                        not null
        primary key,
    user_id          varchar(27)
        references public.users,
    password         text                                               not null,
    refresh_token_id varchar(27),
    created_at       timestamp with time zone default CURRENT_TIMESTAMP not null,
    updated_at       timestamp with time zone default CURRENT_TIMESTAMP not null
);



