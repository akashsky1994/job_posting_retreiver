
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

sudo apt-get update
sudo apt-get install docker-compose-plugin
docker compose version

cd ~
git clone https://github.com/akashsky1994/job_posting_retreiver.git
cd job_posting_retreiver

docker compose build
docker compose up -d