CREATE ROLE pr LOGIN
    ENCRYPTED PASSWORD 'md56786a615d656be84b5b95f264076eaa0'
    NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;

CREATE DATABASE pr
  WITH OWNER = pr
       ENCODING = 'UTF8'
       TABLESPACE = pg_default
       CONNECTION LIMIT = -1;

\connect pr

CREATE TABLE links
(
    hostname text NOT NULL,
    update_at timestamp with time zone DEFAULT now(),
    CONSTRAINT links_pkey PRIMARY KEY (hostname)
)
WITH (
    OIDS=FALSE
);
ALTER TABLE links
    OWNER TO pr;

CREATE TABLE proxies
(
    hostname text NOT NULL,
    host text,
    port text,
    work boolean,
    anon boolean,
    checks integer,
    create_at timestamp with time zone DEFAULT now(),
    update_at timestamp with time zone,
    response integer,
    CONSTRAINT proxies_pkey PRIMARY KEY (hostname)
)
WITH (
    OIDS=FALSE
);
ALTER TABLE proxies
    OWNER TO pr;

INSERT INTO links (hostname) VALUES 
    ('https://hidester.com/proxydata/php/data.php?mykey=data&offset=0&limit=1000&orderBy=latest_check&sortOrder=DESC&country=&port=&type=undefined&anonymity=undefined&ping=undefined&gproxy=2'),
    ('http://gatherproxy.com/embed/'),
    ('http://txt.proxyspy.net/proxy.txt'),
    ('http://webanetlabs.net/publ/24'),
    ('http://awmproxy.com/freeproxy.php'),
    ('http://www.samair.ru/proxy/type-01.htm'),
    ('https://www.us-proxy.org/'),
    ('http://free-proxy-list.net/'),
    ('http://www.proxynova.com/proxy-server-list/'),
    ('http://proxyserverlist-24.blogspot.ru/'),
    ('http://gatherproxy.com/'),
    ('https://hidemy.name/ru/proxy-list/'),
    ('https://hidemy.name/en/proxy-list/?type=hs&anon=34#list'),
    ('https://free-proxy-list.com'),
    ('https://free-proxy-list.com/?search=1&page=&port=&type%5B%5D=http&type%5B%5D=https&level%5B%5D=anonymous&level%5B%5D=high-anonymous&speed%5B%5D=2&speed%5B%5D=3&connect_time%5B%5D=2&connect_time%5B%5D=3&up_time=40&search=Search'),
    ('http://www.idcloak.com/proxylist/free-proxy-servers-list.html'),
    ('https://premproxy.com/list/'),
    ('https://proxy-list.org/english/index.php'),
    ('https://www.sslproxies.org/')
ON CONFLICT (hostname) DO NOTHING;