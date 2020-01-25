-- drop table users;

create table if not exists users
(
    "userId" bigint not null,
    "menuId" int,
    "lat" double precision,
	"lon" double precision,
	"geom" geometry(POINT, 4326),

    primary key ("userId")
);
