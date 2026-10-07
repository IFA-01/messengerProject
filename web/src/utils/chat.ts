export function getPeerName(chatName: string, myUsername: string): string {
  if (myUsername) {
    if (chatName.startsWith(`${myUsername}-`)) return chatName.slice(myUsername.length + 1);
    if (chatName.endsWith(`-${myUsername}`)) return chatName.slice(0, -(myUsername.length + 1));
  }
  return chatName;
}

export function getAvatarLetter(chatName: string, myUsername: string): string {
  return getPeerName(chatName, myUsername).charAt(0).toUpperCase() || '?';
}
