
-- CREATE administrative_regions TABLE
CREATE TABLE administrative_regions (
                                        id integer NOT NULL,
                                        "name" varchar(255) NOT NULL,
                                        name_en varchar(255) NOT NULL,
                                        code_name varchar(255) NULL,
                                        code_name_en varchar(255) NULL,
                                        CONSTRAINT administrative_regions_pkey PRIMARY KEY (id)
);


-- CREATE administrative_units TABLE
CREATE TABLE administrative_units (
                                      id integer NOT NULL,
                                      full_name varchar(255) NULL,
                                      full_name_en varchar(255) NULL,
                                      short_name varchar(255) NULL,
                                      short_name_en varchar(255) NULL,
                                      code_name varchar(255) NULL,
                                      code_name_en varchar(255) NULL,
                                      CONSTRAINT administrative_units_pkey PRIMARY KEY (id)
);


-- CREATE provinces TABLE
CREATE TABLE provinces (
                           code varchar(20) NOT NULL,
                           "name" varchar(255) NOT NULL,
                           name_en varchar(255) NULL,
                           full_name varchar(255) NOT NULL,
                           full_name_en varchar(255) NULL,
                           code_name varchar(255) NULL,
                           administrative_unit_id integer NULL,
                           administrative_region_id integer NULL,
                           CONSTRAINT provinces_pkey PRIMARY KEY (code)
);

-- CREATE districts TABLE
CREATE TABLE districts (
                           code varchar(20) NOT NULL,
                           "name" varchar(255) NOT NULL,
                           name_en varchar(255) NULL,
                           full_name varchar(255) NULL,
                           full_name_en varchar(255) NULL,
                           code_name varchar(255) NULL,
                           province_code varchar(20) NULL,
                           administrative_unit_id integer NULL,
                           CONSTRAINT districts_pkey PRIMARY KEY (code)
);

-- CREATE wards TABLE
CREATE TABLE wards (
                       code varchar(20) NOT NULL,
                       "name" varchar(255) NOT NULL,
                       name_en varchar(255) NULL,
                       full_name varchar(255) NULL,
                       full_name_en varchar(255) NULL,
                       code_name varchar(255) NULL,
                       district_code varchar(20) NULL,
                       administrative_unit_id integer NULL,
                       CONSTRAINT wards_pkey PRIMARY KEY (code)
);
