/**
 * Auth client for token issuance and validation.
 * Connects to backend auth service.
 * By:- Faisal Hanif | imfanee@gmail.com
 */

const getAuthBaseUrl = (): string => {
  if (typeof window !== 'undefined') {
    const host = window.location.hostname;
    return host === 'localhost' ? 'http://localhost:8081' : `${window.location.origin}/auth`;
  }
  return process.env.NEXT_PUBLIC_AUTH_URL ?? 'http://localhost:8081';
};

const getSignalingBaseUrl = (): string => {
  if (typeof window !== 'undefined') {
    const host = window.location.hostname;
    return host === 'localhost' ? 'http://localhost:8080' : `${window.location.origin}/signal`;
  }
  return process.env.NEXT_PUBLIC_SIGNALING_URL ?? 'http://localhost:8080';
};

export async function fetchToken(userId: string): Promise<string> {
  const res = await fetch(`${getAuthBaseUrl()}/auth/token`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ userId }),
  });
  if (!res.ok) {
    throw new Error('Failed to fetch token');
  }
  const data = await res.json();
  return data.token;
}

export function getSignalingUrl(): string {
  return getSignalingBaseUrl();
}
