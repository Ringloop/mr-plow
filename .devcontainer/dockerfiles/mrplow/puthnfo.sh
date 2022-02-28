#!/bin/bash

USER_UID=$(awk 'NR==1' /hnfo)
USER_GID=$(awk 'NR==2' /hnfo)
DKR_GID=$(awk 'NR==3' /hnfo)
VSUID=$(expr $USER_UID + 67)
VSGID=$(expr $USER_GID + 67)

echo putinfo.sh called with user: [$USERNAME, $USER_UID, $USER_GID, $DKR_GID] 
usermod --uid $VSUID vscode
groupmod --gid $VSGID vscode
groupadd --gid $USER_GID $USERNAME
useradd --uid $USER_UID --gid $USER_GID -m -s /bin/bash $USERNAME
cp /home/vscode/.bash_logout /home/$USERNAME
cp /home/vscode/.bashrc /home/$USERNAME
cp /home/vscode/.profile /home/$USERNAME 
cp /home/vscode/.zshrc /home/$USERNAME 
cp -R /home/vscode/.oh-my-zsh /home/$USERNAME 
chown -R $USERNAME:$USERNAME /home/$USERNAME
echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME
chmod 0440 /etc/sudoers.d/$USERNAME
groupmod --gid $DKR_GID docker
usermod -a -G $DKR_GID $USERNAME