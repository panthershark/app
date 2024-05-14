-- migrate:up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- accounts can be users or services
CREATE TABLE public.accounts (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    email TEXT NOT NULL UNIQUE
);

-- accounts with a person attached are a user.
CREATE TABLE public.person (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    account_id uuid NOT NULL REFERENCES public.accounts(id),
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT ''
);

-- storing application configurations in the database can be really helpful.
-- things like feature flags, external configurations, etc. 
-- what's nice - the env only need to be configured with the database url.
-- env settings are not easy to mutate from web interfaces.
CREATE TABLE public.configs (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    slug text NOT NULL UNIQUE,
    config_section jsonb NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

--
-- migrate:down
DROP FUNCTION public.check_account;
DROP TABLE public.person;
DROP TABLE public.accounts;
DROP TABLE public.configs;
DROP EXTENSION IF EXISTS "uuid-ossp";
