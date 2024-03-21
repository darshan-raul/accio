create_postgres:
	docker run --name postgres -e POSTGRES_PASSWORD=asdf123 --network=host -d -v /home/codespace/postgres:/var/lib/postgresql/data postgres