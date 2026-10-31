-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS public.tbl_admin_users (
	id uuid DEFAULT uuid_generate_v4() NOT NULL,
	email varchar(255) NOT NULL,
	password_hash varchar(255) NULL,
	"name" varchar(100) NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	"token" text NULL,
	CONSTRAINT tbl_admin_users_email_key UNIQUE (email),
	CONSTRAINT tbl_admin_users_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.tbl_categories (
	id uuid DEFAULT uuid_generate_v4() NOT NULL,
	"name" varchar(100) NOT NULL,
	slug varchar(120) NOT NULL,
	sort_order int2 DEFAULT 0 NOT NULL,
	is_active bool NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	CONSTRAINT tbl_categories_name_key UNIQUE (name),
	CONSTRAINT tbl_categories_pkey PRIMARY KEY (id),
	CONSTRAINT tbl_categories_slug_key UNIQUE (slug)
);

CREATE TABLE IF NOT EXISTS public.tbl_settings (
	id uuid DEFAULT uuid_generate_v4() NOT NULL,
	site_name varchar(100) NOT NULL,
	wa_number varchar(20) NOT NULL,
	wa_button_text varchar(50) NOT NULL,
	wa_message_template text NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	CONSTRAINT tbl_settings_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.tbl_assets (
	id uuid DEFAULT uuid_generate_v4() NOT NULL,
	code varchar(50) NOT NULL,
	slug varchar(150) NOT NULL,
	category_id uuid NULL,
	title varchar(255) NOT NULL,
	"location" text NOT NULL,
	description text NULL,
	badge varchar(50) NULL,
	price_label varchar(50) NULL,
	price int8 NOT NULL,
	price_negotiable bool NULL,
	view_count int4 NULL,
	status varchar(20) NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	CONSTRAINT tbl_assets_code_key UNIQUE (code),
	CONSTRAINT tbl_assets_pkey PRIMARY KEY (id),
	CONSTRAINT tbl_assets_slug_key UNIQUE (slug),
	CONSTRAINT tbl_assets_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.tbl_categories(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_assets_category_id ON public.tbl_assets USING btree (category_id);
CREATE INDEX IF NOT EXISTS idx_assets_status ON public.tbl_assets USING btree (status);

CREATE TABLE IF NOT EXISTS public.tbl_catageries_attributes (
	id uuid DEFAULT uuid_generate_v4() NOT NULL,
	category_id uuid NULL,
	"name" varchar(100) NOT NULL,
	data_type varchar(30) NOT NULL,
	unit varchar(20) NULL,
	"options" _text DEFAULT '{}'::text[] NULL,
	is_required bool DEFAULT false NOT NULL,
	show_on_card bool DEFAULT false NOT NULL,
	is_active bool DEFAULT true NOT NULL,
	sort_order int2 DEFAULT 0 NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	CONSTRAINT tbl_catageries_attributes_pkey PRIMARY KEY (id),
	CONSTRAINT tbl_catageries_attributes_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.tbl_categories(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_categories_attributes_category_id ON public.tbl_catageries_attributes USING btree (category_id);

CREATE TABLE IF NOT EXISTS public.tbl_asset_attribute_values (
	id uuid DEFAULT uuid_generate_v4() NOT NULL,
	asset_id uuid NULL,
	attribute_id uuid NULL,
	custome_label varchar(100) NULL,
	value_text varchar(255) NULL,
	value_number int4 NULL,
	sort_order int2 NULL,
	CONSTRAINT tbl_asset_attribute_values_pkey PRIMARY KEY (id),
	CONSTRAINT tbl_asset_attribute_values_asset_id_fkey FOREIGN KEY (asset_id) REFERENCES public.tbl_assets(id) ON DELETE CASCADE,
	CONSTRAINT tbl_asset_attribute_values_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.tbl_catageries_attributes(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_asset_attribute_values_asset_id ON public.tbl_asset_attribute_values USING btree (asset_id);
CREATE INDEX IF NOT EXISTS idx_asset_attribute_values_attribute_id ON public.tbl_asset_attribute_values USING btree (attribute_id);

CREATE TABLE IF NOT EXISTS public.tbl_asset_media (
	id uuid DEFAULT uuid_generate_v4() NOT NULL,
	asset_id uuid NULL,
	type_media varchar(20) NULL,
	file_path text NULL,
	sort_order int2 NULL,
	created_at timestamp NULL,
	CONSTRAINT tbl_asset_media_pkey PRIMARY KEY (id),
	CONSTRAINT tbl_asset_media_asset_id_fkey FOREIGN KEY (asset_id) REFERENCES public.tbl_assets(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_asset_media_asset_id ON public.tbl_asset_media USING btree (asset_id);

CREATE TABLE IF NOT EXISTS public.tbl_asset_owner_info (
	asset_id uuid NOT NULL,
	owner_name varchar(100) NOT NULL,
	owner_contact varchar(50) NOT NULL,
	internal_notes text NULL,
	updated_at timestamp NULL,
	CONSTRAINT tbl_asset_owner_info_pkey PRIMARY KEY (asset_id),
	CONSTRAINT tbl_asset_owner_info_asset_id_fkey FOREIGN KEY (asset_id) REFERENCES public.tbl_assets(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS public.tbl_asset_owner_info;
DROP TABLE IF EXISTS public.tbl_asset_media;
DROP TABLE IF EXISTS public.tbl_asset_attribute_values;
DROP TABLE IF EXISTS public.tbl_catageries_attributes;
DROP TABLE IF EXISTS public.tbl_assets;
DROP TABLE IF EXISTS public.tbl_settings;
DROP TABLE IF EXISTS public.tbl_categories;
DROP TABLE IF EXISTS public.tbl_admin_users;