CREATE TABLE "cloudproviders" (
  "id" SERIAL PRIMARY KEY ,
  "name" varchar,
  "slug" varchar
);

INSERT INTO cloudproviders (name, slug) VALUES ('Amazon Web Services', 'aws');
INSERT INTO cloudproviders (name, slug) VALUES ('Google Cloud Platform', 'gcp');
INSERT INTO cloudproviders (name, slug) VALUES ('Microsoft Azure', 'azure');


CREATE TABLE "projects" (
  "id" SERIAL PRIMARY KEY ,
  "name" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "stacks" (
  "id" SERIAL PRIMARY KEY ,
  "name" varchar,
  "project_id" integer  REFERENCES projects (id) ON DELETE CASCADE
);

CREATE TABLE "resources" (
  "id" SERIAL PRIMARY KEY , 
  "name" varchar,
  "cloud_provider_id" integer REFERENCES cloudproviders (id) ON DELETE CASCADE,
  "stack_id" integer REFERENCES stacks (id) ON DELETE CASCADE
);