#!/bin/sh
# prerm script for uchiwa
#

set -e

# summary of how this script can be called and ordering:
#  http://www.debian.org/doc/debian-policy/ch-maintainerscripts.html
#  http://www.proulx.com/~bob/debian/hints/DpkgScriptOrder

# try to stop any running uchiwa services (not all will be running)
stop_uchiwa_service() {
    if [ -x "/etc/init.d/uchiwa" ]; then
        if [ -x "`which invoke-rc.d 2>/dev/null`" ]; then
            invoke-rc.d uchiwa stop || true
        else
            /etc/init.d/uchiwa stop || true
        fi
    fi
}

case "$1" in
    remove|purge)
        stop_uchiwa_service
        ;;

    upgrade|deconfigure)
        ;;

    *)
        echo "prerm called with unknown argument \`$1'" >&2
        exit 1
        ;;
esac

exit 0
