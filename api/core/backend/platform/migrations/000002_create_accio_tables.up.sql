CREATE TABLE "cloudproviders" (
  "id" SERIAL PRIMARY KEY ,
  "name" varchar UNIQUE,
  "slug" varchar UNIQUE
);

INSERT INTO cloudproviders (name, slug) VALUES ('Amazon Web Services', 'aws');
INSERT INTO cloudproviders (name, slug) VALUES ('Google Cloud Platform', 'gcp');
INSERT INTO cloudproviders (name, slug) VALUES ('Microsoft Azure', 'azure');


CREATE TABLE "projects" (
  "id" SERIAL PRIMARY KEY ,
  "name" varchar UNIQUE,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "stacks" (
  "id" SERIAL PRIMARY KEY ,
  "name" varchar UNIQUE,
  "project_id" integer  REFERENCES projects (id) ON DELETE CASCADE,
  "created_at" timestamp,
  "updated_at" timestamp
);


CREATE TABLE "resourcetypes" (
  "id" SERIAL PRIMARY KEY , 
  "name" varchar UNIQUE
);

-- insert initial set of resource types

INSERT INTO resourcetypes (name) VALUES ('virtualmachines');
INSERT INTO resourcetypes (name) VALUES ('containers');
INSERT INTO resourcetypes (name) VALUES ('objectstore');
INSERT INTO resourcetypes (name) VALUES ('filestore');
INSERT INTO resourcetypes (name) VALUES ('blockstore');
INSERT INTO resourcetypes (name) VALUES ('relationaldb');
INSERT INTO resourcetypes (name) VALUES ('nosqldb');
INSERT INTO resourcetypes (name) VALUES ('virutalnetwork');
INSERT INTO resourcetypes (name) VALUES ('cdn');
INSERT INTO resourcetypes (name) VALUES ('loadbalancer');
INSERT INTO resourcetypes (name) VALUES ('secretmanager');
INSERT INTO resourcetypes (name) VALUES ('encryption');
INSERT INTO resourcetypes (name) VALUES ('monitoring');
INSERT INTO resourcetypes (name) VALUES ('queue');
INSERT INTO resourcetypes (name) VALUES ('apigw');


CREATE TABLE "resources" (
  "id" SERIAL PRIMARY KEY , 
  "name" varchar,
  "cloud_prov_id" integer REFERENCES cloudproviders (id) ON DELETE CASCADE,
  "stack_id" integer REFERENCES stacks (id) ON DELETE CASCADE,
  "res_type_id" integer REFERENCES resourcetypes (id) ON DELETE CASCADE,
  "created_at" timestamp,
  "updated_at" timestamp
);
