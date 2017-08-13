ALTER TABLE proxies ALTER COLUMN host SET DEFAULT '';
UPDATE proxies SET host = '' WHERE host IS NULL;
ALTER TABLE proxies ALTER COLUMN host SET NOT NULL;

ALTER TABLE proxies ALTER COLUMN port SET DEFAULT '';
UPDATE proxies SET port = '' WHERE port IS NULL;
ALTER TABLE proxies ALTER COLUMN port SET NOT NULL;

ALTER TABLE proxies ALTER COLUMN work SET DEFAULT FALSE;
UPDATE proxies SET work = FALSE WHERE work IS NULL;
ALTER TABLE proxies ALTER COLUMN work SET NOT NULL;

ALTER TABLE proxies ALTER COLUMN anon SET DEFAULT FALSE;
UPDATE proxies SET anon = FALSE WHERE anon IS NULL;
ALTER TABLE proxies ALTER COLUMN anon SET NOT NULL;

ALTER TABLE proxies ALTER COLUMN checks SET DEFAULT 0;
UPDATE proxies SET checks = 0 WHERE checks IS NULL;
ALTER TABLE proxies ALTER COLUMN checks SET NOT NULL;

ALTER TABLE proxies ALTER COLUMN create_at SET DEFAULT now();
UPDATE proxies SET create_at = now() WHERE create_at IS NULL;
ALTER TABLE proxies ALTER COLUMN create_at SET NOT NULL;

ALTER TABLE proxies ALTER COLUMN update_at SET DEFAULT now();
UPDATE proxies SET update_at = now() WHERE update_at IS NULL;
ALTER TABLE proxies ALTER COLUMN update_at SET NOT NULL;

ALTER TABLE proxies ALTER COLUMN response SET DEFAULT 0;
UPDATE proxies SET response = 0 WHERE response IS NULL;
ALTER TABLE proxies ALTER COLUMN response SET NOT NULL;
