-- WARNING: This schema is for context only and is not meant to be run.
-- Table order and constraints may not be valid for execution.

CREATE TABLE public.tenants (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  name character varying NOT NULL,
  slug character varying NOT NULL UNIQUE,
  status character varying NOT NULL DEFAULT 'active'::character varying,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone NOT NULL DEFAULT now(),
  deleted_at timestamp with time zone,
  description text,
  address text,
  phone_number character varying,
  open_hours character varying,
  logo_url text,
  banner_url text,
  user_id uuid,
  is_public boolean DEFAULT true,
  CONSTRAINT tenants_pkey PRIMARY KEY (id),
  CONSTRAINT fk_tenants_user FOREIGN KEY (user_id) REFERENCES public.users(id)
);
CREATE TABLE public.users (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  email character varying NOT NULL UNIQUE,
  password_hash text NOT NULL,
  full_name character varying NOT NULL,
  role USER-DEFINED NOT NULL,
  is_active boolean NOT NULL DEFAULT true,
  last_login_at timestamp with time zone,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone NOT NULL DEFAULT now(),
  deleted_at timestamp with time zone,
  CONSTRAINT users_pkey PRIMARY KEY (id)
);
CREATE TABLE public.categories (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  tenant_id uuid NOT NULL,
  name character varying NOT NULL,
  description text,
  sort_order integer NOT NULL DEFAULT 0,
  is_active boolean NOT NULL DEFAULT true,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone NOT NULL DEFAULT now(),
  deleted_at timestamp with time zone,
  CONSTRAINT categories_pkey PRIMARY KEY (id),
  CONSTRAINT fk_categories_tenant FOREIGN KEY (tenant_id) REFERENCES public.tenants(id)
);
CREATE TABLE public.menus (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  tenant_id uuid NOT NULL,
  category_id uuid,
  name character varying NOT NULL,
  description text,
  price numeric NOT NULL CHECK (price >= 0::numeric),
  image_url text,
  is_available boolean NOT NULL DEFAULT true,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone NOT NULL DEFAULT now(),
  deleted_at timestamp with time zone,
  CONSTRAINT menus_pkey PRIMARY KEY (id),
  CONSTRAINT fk_menus_tenant FOREIGN KEY (tenant_id) REFERENCES public.tenants(id),
  CONSTRAINT fk_menus_category FOREIGN KEY (category_id) REFERENCES public.categories(id)
);
CREATE TABLE public.dining_tables (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  tenant_id uuid NOT NULL,
  table_name character varying NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone NOT NULL DEFAULT now(),
  deleted_at timestamp with time zone,
  CONSTRAINT dining_tables_pkey PRIMARY KEY (id),
  CONSTRAINT fk_dining_tables_tenant FOREIGN KEY (tenant_id) REFERENCES public.tenants(id)
);
CREATE TABLE public.orders (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  tenant_id uuid NOT NULL,
  user_id uuid,
  dining_tables_id uuid,
  status character varying NOT NULL,
  total_price numeric NOT NULL CHECK (total_price >= 0::numeric),
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone NOT NULL DEFAULT now(),
  deleted_at timestamp with time zone,
  updated_by uuid,
  CONSTRAINT orders_pkey PRIMARY KEY (id),
  CONSTRAINT fk_orders_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
  CONSTRAINT fk_orders_tenant FOREIGN KEY (tenant_id) REFERENCES public.tenants(id),
  CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES public.users(id),
  CONSTRAINT fk_orders_dining_table FOREIGN KEY (dining_tables_id) REFERENCES public.dining_tables(id)
);
CREATE TABLE public.order_items (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  order_id uuid NOT NULL,
  menu_id uuid NOT NULL,
  menu_name character varying NOT NULL,
  quantity integer NOT NULL CHECK (quantity > 0),
  unit_price numeric NOT NULL CHECK (unit_price >= 0::numeric),
  subtotal numeric NOT NULL CHECK (subtotal >= 0::numeric),
  notes text,
  deleted_at timestamp with time zone,
  CONSTRAINT order_items_pkey PRIMARY KEY (id),
  CONSTRAINT fk_order_items_order FOREIGN KEY (order_id) REFERENCES public.orders(id),
  CONSTRAINT fk_order_items_menu FOREIGN KEY (menu_id) REFERENCES public.menus(id)
);