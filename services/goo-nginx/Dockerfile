FROM fedora

MAINTAINER pioh "thepioh@zoho.com"

RUN dnf install --setopt=tsflags=nodocs -y nginx \
 && dnf update --setopt=tsflags=nodocs -y \
 && dnf clean all \
 && ln -sf /dev/stdout /var/log/nginx/access.log \
 && ln -sf /dev/stderr /var/log/nginx/error.log

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]
