#!/bin/sh
# postinst script for uchiwa
#

set -e

# summary of how this script can be called and ordering:
#  http://www.debian.org/doc/debian-policy/ch-maintainerscripts.html
#  http://www.proulx.com/~bob/debian/hints/DpkgScriptOrder

create_uchiwa_user_group() {
    # create uchiwa group if missing
    if ! getent group uchiwa >/dev/null; then
        groupadd -r uchiwa
    fi

    # create sensu group if missing
    if ! getent group sensu >/dev/null; then
        groupadd -r sensu
    fi

    # create uchiwa user
    if ! getent passwd uchiwa >/dev/null; then
        useradd -r -g uchiwa -G sensu -d /opt/uchiwa \
            -s /bin/false -c "Uchiwa, a Sensu dashboard" uchiwa
    fi
}

fix_logrotate_permissions() {
    chmod 644 /etc/logrotate.d/uchiwa
}

case "$1" in
    configure)
        create_uchiwa_user_group
        fix_logrotate_permissions
        ;;

    abort-upgrade|abort-remove|abort-deconfigure)
        ;;

    *)
        echo "postinst called with unknown argument \`$1'" >&2
        exit 1
        ;;
esac

exit 0
