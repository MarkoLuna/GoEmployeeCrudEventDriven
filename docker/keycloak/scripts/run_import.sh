#!/bin/sh

echo "Importing master.."
/opt/bitnami/keycloak/bin/kc.sh import --dir=/import/master --override true 2>/dev/null

echo "Importing dev-realm.."
/opt/bitnami/keycloak/bin/kc.sh import --dir=/import/dev-realm --override true 2>/dev/null

# echo "export realms"
# /opt/bitnami/keycloak/bin/kc.sh export --dir=/export/master
# /opt/bitnami/keycloak/bin/kc.sh export --dir=/export/master --realm master
# /opt/bitnami/keycloak/bin/kc.sh export --dir=/export/dev-realm --realm dev

### Resume normal execution
/opt/bitnami/scripts/keycloak/run.sh
