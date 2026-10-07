import { describe, expect, it } from 'vitest';

import { getAvatarLetter, getPeerName } from './chat';

describe('getPeerName', () => {
  it('returns the other participant when I created the chat', () => {
    expect(getPeerName('IFA-alice', 'IFA')).toBe('alice');
  });

  it('returns the other participant when they created the chat', () => {
    expect(getPeerName('alice-IFA', 'IFA')).toBe('alice');
  });

  it('falls back to the full name when it has no separator', () => {
    expect(getPeerName('General', 'IFA')).toBe('General');
  });

  it('handles usernames that contain hyphens', () => {
    expect(getPeerName('ifa-01-john-doe', 'ifa-01')).toBe('john-doe');
    expect(getPeerName('john-doe-ifa-01', 'ifa-01')).toBe('john-doe');
  });

  it('returns the full name before my username is known', () => {
    expect(getPeerName('IFA-alice', '')).toBe('IFA-alice');
  });
});

describe('getAvatarLetter', () => {
  it('uses the uppercase first letter of the peer', () => {
    expect(getAvatarLetter('IFA-bob', 'IFA')).toBe('B');
  });

  it('returns a placeholder for an empty name', () => {
    expect(getAvatarLetter('', 'IFA')).toBe('?');
  });
});
