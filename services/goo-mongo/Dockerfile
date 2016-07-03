FROM       mongo

ADD . /setup

ENTRYPOINT ["/setup/entrypoint.sh"]

EXPOSE 27017
CMD ["mongod"]
