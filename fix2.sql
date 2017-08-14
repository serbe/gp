ALTER TABLE links ALTER COLUMN update_at SET DEFAULT now();
UPDATE links SET update_at = now() WHERE update_at IS NULL;
ALTER TABLE links ALTER COLUMN update_at SET NOT NULL;

ALTER TABLE links ADD COLUMN iterate boolean;
ALTER TABLE links ALTER COLUMN iterate SET DEFAULT false;
UPDATE links SET iterate = FALSE WHERE iterate IS NULL;
ALTER TABLE links ALTER COLUMN iterate SET NOT NULL;