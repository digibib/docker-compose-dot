.PHONY: all test
IMAGE ?= digibib/docker-compose-dot

all: reload provision run

reload: halt up

halt:
	vagrant halt || true

up:
	vagrant up --no-provision

provision:
	sleep 3 && vagrant provision

run: delete
	@echo "======= RUNNING DOCKER-COMPOSE-DOT CONTAINER ======\n"
	@vagrant ssh -c 'sudo docker run -d --name docker-compose-dot -p 9999:9999 $(IMAGE)'

stop:
	@echo "======= STOPPING DOCKER-COMPOSE-DOT CONTAINER ======\n"
	vagrant ssh -c 'sudo docker stop docker-compose-dot' || true

delete: stop
	@echo "======= DELETING DOCKER-COMPOSE-DOT CONTAINER ======\n"
	vagrant ssh -c 'sudo docker rm docker-compose-dot' || true

test:
	vagrant ssh -c 'docker stats --no-stream docker-compose-dot'

login: # needs EMAIL, PASSWORD, USERNAME
	@ vagrant ssh -c 'docker login --email=$(EMAIL) --username=$(USERNAME) --password=$(PASSWORD)'

TAG = "$(shell git rev-parse HEAD)"

tag:
	vagrant ssh -c 'docker tag -f $(IMAGE) $(IMAGE):$(TAG)'

push: tag
	@echo "======= PUSHING KOHA CONTAINER ======\n"
	vagrant ssh -c 'docker push $(IMAGE):$(TAG)'

docker_cleanup:
	@echo "cleaning up unused containers and images"
	@vagrant ssh -c '/vagrant/docker-cleanup.sh'