#! /bin/ash
CMD="/usr/bin/kamino"

if [[ ! -z "$AS_USER" ]]
then
    USERNAME=$(echo $AS_USER|sed -e 's/:.*//')
    UID=$(echo $AS_USER|sed -e 's/[^:]*://' -e 's/:.*//' )
    GID=$(echo $AS_USER|sed -e 's/.*://')
fi

if [[ -z "$USERNAME" ]] || [[ -z "$UID" ]] ||  [[  -z "$GID" ]]
then
    exec /usr/bin/kamino "$@"
fi

GROUPNAME=$(grep :$GID: /etc/group | sed -e 's/:.*//') 
if [[ -z "$GROUPNAME" ]]
then
    GROUPNAME=$USERNAME
    addgroup -g $GID $GROUPNAME
fi

if ! grep -q x:$UID: /etc/passwd
then
    adduser -D -G $GROUPNAME -u $UID -h /home/$USERNAME $USERNAME
fi

ln -s /.ssh /home/$USERNAME/.ssh
chown -R $USERNAME:$GROUPNAME /.ssh

exec gosu $USERNAME:$GROUPNAME /usr/bin/kamino "$@"