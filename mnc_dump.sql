--
-- PostgreSQL database dump
--

-- Dumped from database version 14.13 (Homebrew)
-- Dumped by pg_dump version 14.13 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: transactions; Type: TABLE; Schema: public; Owner: user1
--

CREATE TABLE public.transactions (
    transaction_id text NOT NULL,
    user_id text,
    type text,
    amount bigint,
    remarks text,
    balance_before bigint,
    balance_after bigint,
    created_date timestamp with time zone,
    top_up_id text
);


ALTER TABLE public.transactions OWNER TO user1;

--
-- Name: users; Type: TABLE; Schema: public; Owner: user1
--

CREATE TABLE public.users (
    user_id text NOT NULL,
    first_name text,
    last_name text,
    phone_number text,
    address text,
    pin text,
    balance bigint,
    created_date timestamp with time zone,
    updated_date timestamp with time zone
);


ALTER TABLE public.users OWNER TO user1;

--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: user1
--

COPY public.transactions (transaction_id, user_id, type, amount, remarks, balance_before, balance_after, created_date, top_up_id) FROM stdin;
5605dd4b-5940-49a5-8292-02d5ad03b45b	89220e37-7ca1-457b-956f-49971967ffd3	CREDIT	500000		0	500000	2024-10-24 00:31:13.171229+07	\N
f8d2c335-9799-4ea7-87cd-d2830d246493	89220e37-7ca1-457b-956f-49971967ffd3		500000		500000	1000000	2024-10-24 00:41:45.198287+07	\N
537e206c-f3a2-45a6-ae1e-dcbc24328ff6	89220e37-7ca1-457b-956f-49971967ffd3		500000		1000000	1500000	2024-10-24 00:42:18.146927+07	\N
88feb768-cd0c-40d9-b90a-539592ee6221	89220e37-7ca1-457b-956f-49971967ffd3		500000		1500000	2000000	2024-10-24 00:46:04.91103+07	6d105c03-e73a-497e-90bc-127206a81689
	89220e37-7ca1-457b-956f-49971967ffd3		500000		2000000	2500000	2024-10-24 00:47:30.467016+07	b684d69f-ba10-47af-bb98-9bb4aa0f5e76
a5f72c19-2b37-4a26-b238-c526be2003c1	89220e37-7ca1-457b-956f-49971967ffd3		500000		3000000	3500000	2024-10-24 00:50:53.093469+07	\N
49244ec2-2f58-4628-98fb-d2bad0d79770	89220e37-7ca1-457b-956f-49971967ffd3	DEBIT	100000	Pulsa Telkomsel 100k	3500000	3400000	2024-10-24 00:53:23.586294+07	\N
383f80ed-050a-48e6-a56c-0cf54d781618	89220e37-7ca1-457b-956f-49971967ffd3	DEBIT	100000	Pulsa Telkomsel 100k	3400000	3300000	2024-10-24 00:53:46.92196+07	\N
b0dc43c7-c6b9-4d54-b4bd-2e208667a44b	89220e37-7ca1-457b-956f-49971967ffd3	DEBIT	30000	Hadiah Ultah	3300000	3270000	2024-10-24 01:19:23.828309+07	\N
92d59871-13d5-47bc-8837-f26a6bd0c487	89220e37-7ca1-457b-956f-49971967ffd3	CREDIT	30000	Hadiah Ultah	3300000	3330000	2024-10-24 01:19:23.829197+07	\N
c5eb3b17-f052-4039-833b-69db6c1c0b43	89220e37-7ca1-457b-956f-49971967ffd3	DEBIT	60000	Hadiah Ultah	3330000	3270000	2024-10-24 01:19:23.875735+07	\N
cfc660ef-31d7-445a-83fa-f0c8d9a81c14	89220e37-7ca1-457b-956f-49971967ffd3	CREDIT	60000	Hadiah Ultah	3330000	3390000	2024-10-24 01:19:23.876403+07	\N
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: user1
--

COPY public.users (user_id, first_name, last_name, phone_number, address, pin, balance, created_date, updated_date) FROM stdin;
e26c138d-c617-4949-a2b1-d1c283715a34	Guntur	Saputro	0811255502	Jl. Kebon Sirih No. 1	123456	0	2024-10-24 00:48:22.84836+07	0001-01-01 07:07:12+07:07:12
89220e37-7ca1-457b-956f-49971967ffd3	Tom	Araya	0811255501	Jl. Diponegoro No. 215	123456	3390000	2024-10-24 00:25:38.303667+07	2024-10-24 02:02:41.035525+07
\.


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: user1
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (transaction_id);


--
-- Name: users uni_users_phone_number; Type: CONSTRAINT; Schema: public; Owner: user1
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_phone_number UNIQUE (phone_number);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: user1
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);


--
-- PostgreSQL database dump complete
--

