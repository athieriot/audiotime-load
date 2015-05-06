# Audiotime Load

Go lang based script to populate a Postgresql database with the content of an Audible datafeed.
Prepare the tables with a Trigram index in order to optimize fuzzy search (And allow Autocompletion)

# Run in Docker

docker build -t load . && docker run -it --rm --link=audiotime_db_1:db --name audiotime-load load
