--
-- PostgreSQL database dump
--

\restrict cg7KBhas2xpT0ODwjipiOoeAe17tjI1WpVEAooFNgljENA49haHLRgPBuPGfDqv

-- Dumped from database version 18.0
-- Dumped by pg_dump version 18.0

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: area_type_enum; Type: TYPE; Schema: public; Owner: cash-cow-user
--

CREATE TYPE public.area_type_enum AS ENUM (
    'city',
    'town',
    'village'
);


ALTER TYPE public.area_type_enum OWNER TO "cash-cow-user";

--
-- Name: cattle_class_enum; Type: TYPE; Schema: public; Owner: cash-cow-user
--

CREATE TYPE public.cattle_class_enum AS ENUM (
    'male_calf',
    'female_calf',
    'steer',
    'heifer',
    'cow',
    'bull'
);


ALTER TYPE public.cattle_class_enum OWNER TO "cash-cow-user";

--
-- Name: cattle_sex; Type: TYPE; Schema: public; Owner: cash-cow-user
--

CREATE TYPE public.cattle_sex AS ENUM (
    'male',
    'female',
    'unknown'
);


ALTER TYPE public.cattle_sex OWNER TO "cash-cow-user";

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: areas; Type: TABLE; Schema: public; Owner: cash-cow-user
--

CREATE TABLE public.areas (
    id bigint NOT NULL,
    region_id bigint NOT NULL,
    name text NOT NULL,
    area_type public.area_type_enum NOT NULL,
    latitude double precision,
    longitude double precision,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.areas OWNER TO "cash-cow-user";

--
-- Name: areas_id_seq; Type: SEQUENCE; Schema: public; Owner: cash-cow-user
--

CREATE SEQUENCE public.areas_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.areas_id_seq OWNER TO "cash-cow-user";

--
-- Name: areas_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: cash-cow-user
--

ALTER SEQUENCE public.areas_id_seq OWNED BY public.areas.id;


--
-- Name: breeds; Type: TABLE; Schema: public; Owner: cash-cow-user
--

CREATE TABLE public.breeds (
    id bigint NOT NULL,
    name text NOT NULL,
    description text,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.breeds OWNER TO "cash-cow-user";

--
-- Name: breeds_id_seq; Type: SEQUENCE; Schema: public; Owner: cash-cow-user
--

CREATE SEQUENCE public.breeds_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.breeds_id_seq OWNER TO "cash-cow-user";

--
-- Name: breeds_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: cash-cow-user
--

ALTER SEQUENCE public.breeds_id_seq OWNED BY public.breeds.id;


--
-- Name: cattle; Type: TABLE; Schema: public; Owner: cash-cow-user
--

CREATE TABLE public.cattle (
    id bigint NOT NULL,
    owner_id bigint NOT NULL,
    breed_id bigint NOT NULL,
    tag_number text NOT NULL,
    sex public.cattle_sex NOT NULL,
    age_months integer NOT NULL,
    weight_kg double precision,
    vaccinations text,
    medical_history text,
    is_pregnant boolean,
    is_castrated boolean,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.cattle OWNER TO "cash-cow-user";

--
-- Name: cattle_id_seq; Type: SEQUENCE; Schema: public; Owner: cash-cow-user
--

CREATE SEQUENCE public.cattle_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.cattle_id_seq OWNER TO "cash-cow-user";

--
-- Name: cattle_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: cash-cow-user
--

ALTER SEQUENCE public.cattle_id_seq OWNED BY public.cattle.id;


--
-- Name: listing_prices; Type: TABLE; Schema: public; Owner: cash-cow-user
--

CREATE TABLE public.listing_prices (
    listing_id bigint NOT NULL,
    cattle_class public.cattle_class_enum NOT NULL,
    price_per_kg numeric(10,2) NOT NULL,
    quantity integer NOT NULL
);


ALTER TABLE public.listing_prices OWNER TO "cash-cow-user";

--
-- Name: listings; Type: TABLE; Schema: public; Owner: cash-cow-user
--

CREATE TABLE public.listings (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    area_id bigint NOT NULL,
    region_id bigint NOT NULL,
    title text NOT NULL,
    description text,
    latitude double precision,
    longitude double precision,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.listings OWNER TO "cash-cow-user";

--
-- Name: listings_cattle; Type: TABLE; Schema: public; Owner: cash-cow-user
--

CREATE TABLE public.listings_cattle (
    listing_id bigint NOT NULL,
    cattle_id bigint NOT NULL
);


ALTER TABLE public.listings_cattle OWNER TO "cash-cow-user";

--
-- Name: listings_id_seq; Type: SEQUENCE; Schema: public; Owner: cash-cow-user
--

CREATE SEQUENCE public.listings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.listings_id_seq OWNER TO "cash-cow-user";

--
-- Name: listings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: cash-cow-user
--

ALTER SEQUENCE public.listings_id_seq OWNED BY public.listings.id;


--
-- Name: permissions; Type: TABLE; Schema: public; Owner: cash-cow-user
--

CREATE TABLE public.permissions (
    id bigint NOT NULL,
    code text NOT NULL
);


ALTER TABLE public.permissions OWNER TO "cash-cow-user";

--
-- Name: permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: cash-cow-user
--

CREATE SEQUENCE public.permissions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.permissions_id_seq OWNER TO "cash-cow-user";

--
-- Name: permissions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: cash-cow-user
--

ALTER SEQUENCE public.permissions_id_seq OWNED BY public.permissions.id;


--
-- Name: regions; Type: TABLE; Schema: public; Owner: cash-cow-user
--

CREATE TABLE public.regions (
    id bigint NOT NULL,
    name text NOT NULL,
    code text NOT NULL
);


ALTER TABLE public.regions OWNER TO "cash-cow-user";

--
-- Name: regions_id_seq; Type: SEQUENCE; Schema: public; Owner: cash-cow-user
--

CREATE SEQUENCE public.regions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.regions_id_seq OWNER TO "cash-cow-user";

--
-- Name: regions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: cash-cow-user
--

ALTER SEQUENCE public.regions_id_seq OWNED BY public.regions.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: cash-cow-user
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO "cash-cow-user";

--
-- Name: tokens; Type: TABLE; Schema: public; Owner: cash-cow-user
--

CREATE TABLE public.tokens (
    hash bytea NOT NULL,
    user_id bigint NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    scope text NOT NULL
);


ALTER TABLE public.tokens OWNER TO "cash-cow-user";

--
-- Name: users; Type: TABLE; Schema: public; Owner: cash-cow-user
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    email text NOT NULL,
    password_hash text NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    farmer_id integer,
    phone_number text,
    is_activated boolean DEFAULT false NOT NULL,
    is_verified boolean DEFAULT false NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.users OWNER TO "cash-cow-user";

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: cash-cow-user
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO "cash-cow-user";

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: cash-cow-user
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: users_permissions; Type: TABLE; Schema: public; Owner: cash-cow-user
--

CREATE TABLE public.users_permissions (
    user_id bigint NOT NULL,
    permission_id bigint NOT NULL
);


ALTER TABLE public.users_permissions OWNER TO "cash-cow-user";

--
-- Name: areas id; Type: DEFAULT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.areas ALTER COLUMN id SET DEFAULT nextval('public.areas_id_seq'::regclass);


--
-- Name: breeds id; Type: DEFAULT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.breeds ALTER COLUMN id SET DEFAULT nextval('public.breeds_id_seq'::regclass);


--
-- Name: cattle id; Type: DEFAULT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.cattle ALTER COLUMN id SET DEFAULT nextval('public.cattle_id_seq'::regclass);


--
-- Name: listings id; Type: DEFAULT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.listings ALTER COLUMN id SET DEFAULT nextval('public.listings_id_seq'::regclass);


--
-- Name: permissions id; Type: DEFAULT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.permissions ALTER COLUMN id SET DEFAULT nextval('public.permissions_id_seq'::regclass);


--
-- Name: regions id; Type: DEFAULT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.regions ALTER COLUMN id SET DEFAULT nextval('public.regions_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: areas areas_pkey; Type: CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.areas
    ADD CONSTRAINT areas_pkey PRIMARY KEY (id);


--
-- Name: breeds breeds_name_key; Type: CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.breeds
    ADD CONSTRAINT breeds_name_key UNIQUE (name);


--
-- Name: breeds breeds_pkey; Type: CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.breeds
    ADD CONSTRAINT breeds_pkey PRIMARY KEY (id);


--
-- Name: cattle cattle_pkey; Type: CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.cattle
    ADD CONSTRAINT cattle_pkey PRIMARY KEY (id);


--
-- Name: cattle cattle_tag_number_key; Type: CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.cattle
    ADD CONSTRAINT cattle_tag_number_key UNIQUE (tag_number);


--
-- Name: listing_prices listing_prices_pkey; Type: CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.listing_prices
    ADD CONSTRAINT listing_prices_pkey PRIMARY KEY (listing_id, cattle_class);


--
-- Name: listings_cattle listings_cattle_pkey; Type: CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.listings_cattle
    ADD CONSTRAINT listings_cattle_pkey PRIMARY KEY (listing_id, cattle_id);


--
-- Name: listings listings_pkey; Type: CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.listings
    ADD CONSTRAINT listings_pkey PRIMARY KEY (id);


--
-- Name: permissions permissions_code_key; Type: CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_code_key UNIQUE (code);


--
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);


--
-- Name: regions regions_code_key; Type: CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.regions
    ADD CONSTRAINT regions_code_key UNIQUE (code);


--
-- Name: regions regions_name_key; Type: CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.regions
    ADD CONSTRAINT regions_name_key UNIQUE (name);


--
-- Name: regions regions_pkey; Type: CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.regions
    ADD CONSTRAINT regions_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: tokens tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_pkey PRIMARY KEY (hash);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users_permissions users_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.users_permissions
    ADD CONSTRAINT users_permissions_pkey PRIMARY KEY (user_id, permission_id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: areas fk_areas_region_id; Type: FK CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.areas
    ADD CONSTRAINT fk_areas_region_id FOREIGN KEY (region_id) REFERENCES public.regions(id) ON DELETE RESTRICT;


--
-- Name: cattle fk_cattle_breed_id; Type: FK CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.cattle
    ADD CONSTRAINT fk_cattle_breed_id FOREIGN KEY (breed_id) REFERENCES public.breeds(id) ON DELETE RESTRICT;


--
-- Name: listings_cattle fk_cattle_id; Type: FK CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.listings_cattle
    ADD CONSTRAINT fk_cattle_id FOREIGN KEY (cattle_id) REFERENCES public.cattle(id) ON DELETE CASCADE;


--
-- Name: cattle fk_cattle_owner_id; Type: FK CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.cattle
    ADD CONSTRAINT fk_cattle_owner_id FOREIGN KEY (owner_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: listings_cattle fk_listing_id; Type: FK CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.listings_cattle
    ADD CONSTRAINT fk_listing_id FOREIGN KEY (listing_id) REFERENCES public.listings(id) ON DELETE CASCADE;


--
-- Name: listing_prices fk_listing_id; Type: FK CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.listing_prices
    ADD CONSTRAINT fk_listing_id FOREIGN KEY (listing_id) REFERENCES public.listings(id) ON DELETE CASCADE;


--
-- Name: listings fk_listings_area_id; Type: FK CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.listings
    ADD CONSTRAINT fk_listings_area_id FOREIGN KEY (area_id) REFERENCES public.areas(id) ON DELETE CASCADE;


--
-- Name: listings fk_listings_region_id; Type: FK CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.listings
    ADD CONSTRAINT fk_listings_region_id FOREIGN KEY (region_id) REFERENCES public.regions(id) ON DELETE CASCADE;


--
-- Name: listings fk_listings_user_id; Type: FK CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.listings
    ADD CONSTRAINT fk_listings_user_id FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: tokens fk_tokens_user_id; Type: FK CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT fk_tokens_user_id FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: users_permissions fk_users_permissions_permission_id; Type: FK CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.users_permissions
    ADD CONSTRAINT fk_users_permissions_permission_id FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE;


--
-- Name: users_permissions fk_users_permissions_user_id; Type: FK CONSTRAINT; Schema: public; Owner: cash-cow-user
--

ALTER TABLE ONLY public.users_permissions
    ADD CONSTRAINT fk_users_permissions_user_id FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict cg7KBhas2xpT0ODwjipiOoeAe17tjI1WpVEAooFNgljENA49haHLRgPBuPGfDqv

