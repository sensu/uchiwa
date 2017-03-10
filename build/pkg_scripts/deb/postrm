#!/bin/sh
# postrm script for uchiwa
#

set -e

# summary of how this script can be called and ordering:
#  http://www.debian.org/doc/debian-policy/ch-maintainerscripts.html
#  http://www.proulx.com/~bob/debian/hints/DpkgScriptOrder

purge_uchiwa_service() {
    update-rc.d uchiwa remove >/dev/null || true
    if [ -f "/etc/init.d/uchiwa" ]; then
        rm /etc/init.d/uchiwa
    fi
}

purge_uchiwa_files() {
    if [ -d "/opt/uchiwa" ]; then
        rm -r /opt/uchiwa
    fi
}

purge_uchiwa_user_group() {
    if getent passwd uchiwa >/dev/null; then
        userdel -f uchiwa
    fi
    if getent group uchiwa >/dev/null; then
        groupdel -f uchiwa
    fi
}

case "$1" in
    purge)
        purge_uchiwa_service
        purge_uchiwa_files
        purge_uchiwa_user_group
        ;;

    remove|upgrade|abort-upgrade|abort-remove|abort-deconfigure)
        ;;

    *)
        echo "postinst called with unknown argument \`$1'" >&2
        exit 1
        ;;
esac

exit 0
